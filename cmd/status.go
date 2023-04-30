package cmd

import (
	"encoding/json"
	"log"
	"strings"
	"sync"
)

type ResultSetT struct {
	SMS       [][]SMSData              `json:"sms"`
	MMS       [][]MMSData              `json:"mms"`
	VoiceCall []VoiceCallData          `json:"voice_call"`
	Email     map[string][][]EmailData `json:"email"`
	Billing   BillingData              `json:"billing"`
	Support   []int                    `json:"support"`
	Incidents []IncidentData           `json:"incident"`
}

func getResultData() ResultSetT {

	resultData := &ResultSetT{}
	wg := sync.WaitGroup{}
	wg.Add(7)

	go func() {
		if err := makeDoubleList(resultData, getSMSstatus(),
			"country", "provider"); err != nil {
			log.Printf("SMS: %s", err)
		}
		wg.Done()
	}()

	go func() {
		if err := makeDoubleList(resultData, getMMSstatus(),
			"country", "provider"); err != nil {
			log.Printf("MMS: %s", err)
		}
		wg.Done()
	}()

	go func() {
		resultData.VoiceCall = getVoiceCallStatus()
		wg.Done()
	}()

	go func() {
		resultData.Email = makeEmailMap(getEmailStatus())
		wg.Done()
	}()

	go func() {
		if billing := getBillingStatus(); billing != nil {
			resultData.Billing = billing[0]
		}
		wg.Done()
	}()

	go func() {
		resultData.Support = makeSupportReport(getSupportStatus())
		wg.Done()
	}()

	go func() {
		resultData.Incidents = makeIncidentReport(getIncidentStatus())
		wg.Done()
	}()

	wg.Wait()
	return *resultData
}

type ResultT struct {
	Status bool       `json:"status"`
	Data   ResultSetT `json:"data"`
	Error  string     `json:"error"`
}

func getResultT() *ResultT {

	res := &ResultT{}
	data := getResultData()
	res.Status = getResultStatus(data)
	if res.Status {
		res.Data = data
	} else {
		res.Error = "Error on collect data"
	}
	return res
}

func getResultTEmail(country string) ([]byte, error) {

	res := getResultT()

	if country != "" {
		for c := range res.Data.Email {
			if c == country {
				continue
			}
			delete(res.Data.Email, c)
		}
	}

	tempEmail := make(map[string][][]EmailData)
	for c, val := range res.Data.Email {
		temp, _ := getCountryISO3166_1(c, 2)
		tempEmail[temp] = val
	}
	res.Data.Email = tempEmail

	data, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return nil, err
	}

	if country != "" {
		c, _ := getCountryISO3166_1(country, 2)
		str := string(data)

		str = str[:strings.Index(str, `"email": `)+9] +
			str[strings.Index(str, `"email": `)+9+9+len(c)+3:]

		str = str[:strings.Index(str, ",\n    \"billing\": ")-6] +
			str[strings.Index(str, ",\n    \"billing\": "):]

		data = []byte(str)
	}

	return data, nil
}
