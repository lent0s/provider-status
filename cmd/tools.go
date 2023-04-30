package cmd

import (
	"fmt"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

func fieldIndex(jsonField string, s reflect.Type) (index int, err error) {

	for i := 0; i < s.NumField(); i++ {
		if s.Field(i).Tag.Get("json") == jsonField {
			return i, nil
		}
	}
	return -1, fmt.Errorf("sort: unknown field %s\n", jsonField)
}

// view returns: 0 - "###", 1 - "XXX", 2 - full name (ENG),
// 3 - full name (RUS), another number - "XX"
func getCountryISO3166_1(country string, view int) (string, error) {

	list, err := getISO3166_1()
	if err != nil {
		return "", err
	}

	toView := ""

	switch len(country) {
	case 2:
		country = strings.ToTitle(country)
		if list[country] != nil {
			toView = country
		}
	case 3:
		_, err := strconv.Atoi(country)
		if err == nil {
			for key, val := range list {
				if country == val[0] {
					toView = key
				}
			}
		} else {
			country = strings.ToTitle(country)
			for key, val := range list {
				if country == val[1] {
					toView = key
				}
			}
		}
	default:
		if len(country) < 2 {
			break
		}
		country = strings.ToLower(country)
		for key, val := range list {
			if country == strings.ToLower(val[2]) ||
				country == strings.ToLower(val[3]) {
				toView = key
			}
		}
	}

	if toView == "" {
		return "", fmt.Errorf("country [%s] is not on the ISO3166_1\n",
			country)
	}

	switch view {
	case 0, 1, 2, 3:
		return list[toView][view], nil
	default:
		return toView, nil
	}
}

func getISO3166_1() (map[string][]string, error) {

	m := make(map[string][]string)

	data, err := getRowsCSV("cmd/lib/countryISO3166_1.txt")
	if err != nil {
		return nil, err
	}

	for _, str := range data {
		if str == "" || str[0] != '"' {
			continue
		}
		endEN := 22 + strings.Index(str[22:], `"`)

		m[str[1:3]] = []string{
			str[8:11],
			str[15:18],
			str[22:endEN],
			str[endEN+4 : len(str)-2]}
	}

	return m, nil
}

func getResultStatus(res ResultSetT) bool {

	if res.SMS == nil || res.MMS == nil || res.VoiceCall == nil ||
		res.Email == nil || res.Support == nil || res.Incidents == nil {
		return false
	}
	if getBillingStatus() == nil {
		return false
	}
	return true
}

// push double slice of list with full country name and sorted on f1 & f2 fields
// to res
func makeDoubleList(res *ResultSetT, list interface{}, f1, f2 string) error {

	s := reflect.ValueOf(list)
	if s.Kind() != reflect.Slice {
		return fmt.Errorf("DoubleList: not a slice\n")
	}
	if s.Len() == 0 {
		return nil
	}

	f1 = strings.ToLower(f1)
	f2 = strings.ToLower(f2)
	elementType := s.Type().Elem()
	fieldIndex, err := fieldIndex(f1, elementType)
	if err != nil {
		return fmt.Errorf("DoubleList: %s", err)
	}

	for data := 0; data < s.Len(); data++ {
		temp, _ := getCountryISO3166_1(s.Index(data).
			Field(fieldIndex).String(), 2)
		s.Index(data).Field(fieldIndex).SetString(temp)
	}

	if err := sortListByField(list, f1, 1); err != nil {
		return fmt.Errorf("DoubleList: %s", err)
	}
	temp := reflect.MakeSlice(s.Type(), s.Len(), s.Len())
	reflect.Copy(temp, s)

	if err := sortListByField(list, f2, 1); err != nil {
		return fmt.Errorf("DoubleList: %s", err)
	}

	r := reflect.ValueOf(res)
	for i := 0; i < r.Elem().NumField(); i++ {
		if reflect.ValueOf(*res).Field(i).Type() == reflect.SliceOf(s.Type()) {
			fieldIndex = i
		}
	}
	r.Elem().Field(fieldIndex).
		Set(reflect.Append(r.Elem().Field(fieldIndex), temp, s))

	return nil
}

// option make direction of sorted list: 0 - [Z-A, Я-А, 9-0, true-false],
// another - [A-Z, А-Я, 0-9, false-true]
func sortListByField(list interface{}, jsonField string, option int) error {

	s := reflect.ValueOf(list)
	if s.Kind() != reflect.Slice {
		return fmt.Errorf("sort: not a slice\n")
	}
	if s.Len() == 0 || s.Len() == 1 {
		return nil
	}

	elementType := s.Type().Elem()
	jsonField = strings.ToLower(jsonField)
	fieldIndex, err := fieldIndex(jsonField, elementType)
	if err != nil {
		return err
	}

	less := func(i, j int) bool {
		var vi, vj reflect.Value
		switch option {
		case 0:
			vj = s.Index(i).Field(fieldIndex)
			vi = s.Index(j).Field(fieldIndex)
		default:
			vi = s.Index(i).Field(fieldIndex)
			vj = s.Index(j).Field(fieldIndex)
		}

		switch vi.Interface().(type) {
		case string:
			return reflect.ValueOf(vi.Interface()).String() <
				reflect.ValueOf(vj.Interface()).String()
		case float32:
			return reflect.ValueOf(vi.Interface()).Float() <
				reflect.ValueOf(vj.Interface()).Float()
		case float64:
			return reflect.ValueOf(vi.Interface()).Float() <
				reflect.ValueOf(vj.Interface()).Float()
		case uint64:
			return reflect.ValueOf(vi.Interface()).Uint() <
				reflect.ValueOf(vj.Interface()).Uint()
		default:
			return reflect.ValueOf(vi.Interface()).Int() <
				reflect.ValueOf(vj.Interface()).Int()
		}
	}

	ifBool := func() {
		dir := false
		if option != 0 {
			dir = true
		}
		max := s.Len() - 1
		temp := reflect.New(s.Index(0).Type()).Elem()
		for i := 0; i < max; i++ {
			if s.Index(i).Field(fieldIndex).Bool() == dir {
				temp.Set(s.Index(i))
				s.Index(i).Set(s.Index(max))
				s.Index(max).Set(temp)
				max--
				i--
			}
		}
	}

	if elementType.Field(fieldIndex).Type.Kind() == reflect.Bool {
		ifBool()
	} else {
		sort.SliceStable(list, less)
	}
	return nil
}

func sysType() int {

	if runtime.GOARCH == "amd64" {
		return 64
	} else {
		return 32
	}
}
