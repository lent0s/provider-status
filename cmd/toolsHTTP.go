package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func checkCorrectIncidentStatus(s string) (string, error) {

	allowed := []string{"active", "closed"}

	s = strings.ToLower(s)

	for _, val := range allowed {
		if s == val {
			return val, nil
		}
	}
	return "", fmt.Errorf("wrong incident status [%s]\n"+
		"allowed: %s\n", s, allowed)
}

func getResponseHTTP(url string) ([]byte, error) {

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response failed with status code"+
			" [%d] and body:\n%s", res.StatusCode, body)
	}
	if err != nil {
		return nil, err
	}
	return body, nil
}

func makeIncidentReport(i []IncidentData) []IncidentData {

	if err := sortListByField(i, "status", 1); err != nil {
		log.Println(err)
	}

	return i
}

func makeSupportReport(s []SupportData) []int {

	var (
		tickets        = 0
		usage          = 2
		technicians    = 7
		ticketsPerHour = 18
	)

	for _, ticket := range s {
		tickets += ticket.ActiveTickets
	}
	if tickets < 9 {
		usage = 1
	}
	if tickets > 16 {
		usage = 3
	}

	responseTimeMinutes := int((3600 / (float32(ticketsPerHour) /
		7.0 * float32(technicians))) * float32(tickets) / 60)

	return []int{usage, responseTimeMinutes}
}
