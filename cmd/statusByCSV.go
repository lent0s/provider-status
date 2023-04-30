package cmd

import (
	"log"
	"sort"
	"strconv"
	"strings"
)

type SMSData struct {
	Country      string `json:"country"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
	Provider     string `json:"provider"`
}

func getSMSstatus() (SMSStatus []SMSData) {

	rows, err := getRowsCSV(c.filenameSMSData)
	if err != nil {
		log.Printf("SMSstatus: %s\n", err)
		return
	}

	for num, row := range rows {
		if len(row) == 0 {
			log.Printf("row [%d] is empty\n", num+1)
			continue
		}

		field := strings.Split(row, ";")
		if len(field) != 4 {
			log.Printf("row [%d] is corrupted\n", num+1)
			continue
		}
		SMSdata := &SMSData{}
		if SMSdata.Country,
			err = getCountryISO3166_1(field[0], 69); err != nil {
			log.Printf("row [%d]: %s", num+1, err)
			continue
		}
		if _, err := strconv.Atoi(field[1]); err != nil {
			log.Printf("row [%d] is corrupted: %s", num+1, err)
			continue
		}
		SMSdata.Bandwidth = field[1]

		if _, err = strconv.Atoi(field[2]); err != nil {
			log.Printf("row [%d] is corrupted: %s", num+1, err)
			continue
		}
		SMSdata.ResponseTime = field[2]

		if SMSdata.Provider,
			err = providerNotCorrect(field[3], "SMS"); err != nil {
			log.Printf("row [%d]: %s", num+1, err)
			continue
		}

		SMSStatus = append(SMSStatus, *SMSdata)
	}
	return
}

type VoiceCallData struct {
	Country             string  `json:"country"`
	Bandwidth           string  `json:"bandwidth"`
	ResponseTime        string  `json:"response_time"`
	Provider            string  `json:"provider"`
	ConnectionStability float32 `json:"connection_stability"`
	TTFB                int     `json:"ttfb"`
	VoicePurity         int     `json:"voice_purity"`
	MedianOfCallsTime   int     `json:"median_of_calls_time"`
}

func getVoiceCallStatus() (VCStatus []VoiceCallData) {

	rows, err := getRowsCSV(c.filenameVoiceData)
	if err != nil {
		log.Printf("VoiceCallStatus: %s", err)
		return
	}

	for num, row := range rows {
		if len(row) == 0 {
			log.Printf("row [%d] is empty\n", num+1)
			continue
		}

		field := strings.Split(row, ";")
		if len(field) != 8 {
			log.Printf("row [%d] is corrupted\n", num+1)
			continue
		}
		VoiceData := &VoiceCallData{}

		if VoiceData.Country,
			err = getCountryISO3166_1(field[0], 69); err != nil {
			log.Printf("row [%d]: %s", num+1, err)
			continue
		}
		if _, err = strconv.Atoi(field[1]); err != nil {
			log.Printf("row [%d] is corrupted: %s", num+1, err)
			continue
		} else {
			VoiceData.Bandwidth = field[1]
		}
		if _, err = strconv.Atoi(field[2]); err != nil {
			log.Printf("row [%d] is corrupted: %s", num+1, err)
			continue
		} else {
			VoiceData.ResponseTime = field[2]
		}
		if VoiceData.Provider, err = providerNotCorrect(field[3],
			"Voice"); err != nil {
			log.Printf("row [%d]: %s", num+1, err)
			continue
		}

		connectionStability64, err := strconv.ParseFloat(field[4], sysType())
		if err != nil {
			log.Printf("row [%d] is corrupted: %s", num+1, err)
			continue
		} else {
			VoiceData.ConnectionStability = float32(connectionStability64)
		}

		if VoiceData.TTFB, err = strconv.Atoi(field[5]); err != nil {
			log.Printf("row [%d] is corrupted: %s", num+1, err)
			continue
		}
		if VoiceData.VoicePurity, err = strconv.Atoi(field[6]); err != nil {
			log.Printf("row [%d] is corrupted: %s", num+1, err)
			continue
		}
		if VoiceData.MedianOfCallsTime,
			err = strconv.Atoi(field[7]); err != nil {
			log.Printf("row [%d] is corrupted: %s", num+1, err)
			continue
		}

		VCStatus = append(VCStatus, *VoiceData)
	}
	sort.Slice(VCStatus, func(i, j int) bool {
		return VCStatus[i].Country < VCStatus[j].Country
	})
	return
}

type EmailData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	DeliveryTime int    `json:"delivery_time"`
}

func getEmailStatus() (EmailStatus []EmailData) {

	rows, err := getRowsCSV(c.filenameEmailData)
	if err != nil {
		log.Printf("EmailStatus: %s\n", err)
		return
	}

	for num, row := range rows {
		if len(row) == 0 {
			log.Printf("row [%d] is empty\n", num+1)
			continue
		}

		field := strings.Split(row, ";")
		if len(field) != 3 {
			log.Printf("row [%d] is corrupted\n", num+1)
			continue
		}
		emailData := &EmailData{}

		if emailData.Country,
			err = getCountryISO3166_1(field[0], 69); err != nil {
			log.Printf("row [%d]: %s", num+1, err)
			continue
		}
		if emailData.Provider, err = providerNotCorrect(field[1],
			"Email"); err != nil {
			log.Printf("row [%d]: %s", num+1, err)
			continue
		}
		if emailData.DeliveryTime, err = strconv.Atoi(field[2]); err != nil {
			log.Printf("row [%d] is corrupted: %s", num+1, err)
			continue
		}

		EmailStatus = append(EmailStatus, *emailData)
	}
	return
}

type BillingData struct {
	CreateCustomer bool `json:"create_customer"`
	Purchase       bool `json:"purchase"`
	Payout         bool `json:"payout"`
	Recurring      bool `json:"recurring"`
	FraudControl   bool `json:"fraud_control"`
	CheckoutPage   bool `json:"checkout_page"`
}

func getBillingStatus() (billingStatus []BillingData) {

	rows, err := getRowsCSV(c.filenameBillingData)
	if err != nil {
		log.Printf("BillingStatus: %s\n", err)
		return
	}
	for num, row := range rows {
		if len(row) == 0 {
			log.Printf("row [%d] is empty\n", num+1)
			continue
		}
		if []byte(row)[len([]byte(row))-1] == 13 {
			row = string([]byte(row)[:len([]byte(row))-1])
		}
		if len(row) != 6 {
			log.Printf("row [%d] is corrupted\n", num+1)
			continue
		}

		data, err := strconv.ParseInt(row, 2, 64)
		if err != nil {
			log.Printf("row [%d] is corrupted: %s", num+1, err)
			continue
		}
		dataDec := uint8(data)

		billingData := &BillingData{}
		for i := 1; i < 64; i *= 2 {
			switch dataDec & uint8(i) {
			case 1:
				billingData.CreateCustomer = true
			case 2:
				billingData.Purchase = true
			case 4:
				billingData.Payout = true
			case 8:
				billingData.Recurring = true
			case 16:
				billingData.FraudControl = true
			case 32:
				billingData.CheckoutPage = true
			}
		}
		billingStatus = append(billingStatus, *billingData)
	}

	return
}
