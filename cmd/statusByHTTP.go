package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

type MMSData struct {
	Country      string `json:"country"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
	Provider     string `json:"provider"`
}

func getMMSstatus() (MMSStatus []MMSData) {

	url := fmt.Sprintf("http://%s/mms", c.host)
	body, err := getResponseHTTP(url)
	if err != nil {
		log.Printf("MMS status: %s", err)
		return
	}
	if err = json.Unmarshal(body, &MMSStatus); err != nil {
		log.Printf("MMS status: %s", err)
		return
	}

	for num, val := range MMSStatus {
		if MMSStatus[num].Country,
			err = getCountryISO3166_1(val.Country, 69); err != nil {
			log.Printf("rec [%d]: %s", num+1, err)
			MMSStatus = append(MMSStatus[:num], MMSStatus[num+1:]...)
			continue
		}
		if MMSStatus[num].Provider,
			err = providerNotCorrect(val.Provider, "MMS"); err != nil {
			log.Printf("rec [%d]: %s", num+1, err)
			MMSStatus = append(MMSStatus[:num], MMSStatus[num+1:]...)
			continue
		}
		if _, err = strconv.Atoi(val.Bandwidth); err != nil {
			log.Printf("rec [%d] is corrupted: %s", num+1, err)
			MMSStatus = append(MMSStatus[:num], MMSStatus[num+1:]...)
			continue
		}
		if _, err = strconv.Atoi(val.ResponseTime); err != nil {
			log.Printf("rec [%d] is corrupted: %s", num+1, err)
			MMSStatus = append(MMSStatus[:num], MMSStatus[num+1:]...)
			continue
		}
	}

	return
}

type SupportData struct {
	Topic         string `json:"topic"`
	ActiveTickets int    `json:"active_tickets"`
}

func getSupportStatus() (supportStatus []SupportData) {

	url := fmt.Sprintf("http://%s/support", c.host)
	body, err := getResponseHTTP(url)
	if err != nil {
		log.Printf("Support status: %s", err)
		return
	}
	if err = json.Unmarshal(body, &supportStatus); err != nil {
		log.Printf("Support status: %s", err)
		return
	}

	return
}

type IncidentData struct {
	Topic  string `json:"topic"`
	Status string `json:"status"`
}

func getIncidentStatus() (incidentStatus []IncidentData) {

	url := fmt.Sprintf("http://%s/accendent", c.host)
	body, err := getResponseHTTP(url)
	if err != nil {
		log.Printf("Incident status: %s", err)
		return
	}
	if err = json.Unmarshal(body, &incidentStatus); err != nil {
		log.Printf("Incident status: %s", err)
		return
	}

	for num, val := range incidentStatus {
		if incidentStatus[num].Status, err = checkCorrectIncidentStatus(val.Status); err != nil {
			log.Printf("rec [%d]: %s", num+1, err)
			incidentStatus = append(incidentStatus[:num], incidentStatus[num+1:]...)
			continue
		}
	}

	return
}
