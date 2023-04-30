package cmd

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const fileConfig = "cmd/config.ini"

type config struct {
	serverHost     string
	refreshTimeSec int64

	filenameSMSData     string
	filenameVoiceData   string
	filenameEmailData   string
	filenameBillingData string
	host                string
}

var c config

func ReadConfig() {

	conf, err := readConfig()
	if err != nil {
		log.Fatalf("config: %s", err)
	}
	c = *conf
}

func readConfig() (*config, error) {

	data, err := openConfig()
	if err != nil {
		return nil, err
	}
	records := readDataConfig(data)

	if err = checkFilePath(records); err != nil {
		return nil, err
	}
	if err = checkHost(records); err != nil {
		return nil, err
	}

	refreshTime, err := strconv.Atoi(records["dataIsUpToDateForSeconds:"])
	if err != nil {
		if newErr := appendConfig("dataIsUpToDateForSeconds:"); newErr != nil {
			return nil, newErr
		}
		return nil, fmt.Errorf("incorrect [dataIsUpToDateForSeconds] "+
			"check \"%s\"\n%s", fileConfig, err)
	}

	return &config{
		serverHost:          records["serverHost:"],
		refreshTimeSec:      int64(refreshTime),
		filenameSMSData:     records["filenameSMSData:"],
		filenameVoiceData:   records["filenameVoiceData:"],
		filenameEmailData:   records["filenameEmailData:"],
		filenameBillingData: records["filenameBillingData:"],
		host:                records["incomingDataHost:"],
	}, nil
}

func openConfig() ([]byte, error) {

	data, err := os.ReadFile(fileConfig)
	if err != nil {
		if newErr := makeConfig(); newErr != nil {
			return nil, newErr
		}
		return nil, fmt.Errorf("%s\n"+
			"[%s] is corrupt and has been replaced\n"+
			"check it and rerun application", err, fileConfig)
	}
	return data, nil
}

func makeConfig() error {

	file, err := os.OpenFile(fileConfig, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Println(err)
		}
	}()

	_, err = file.Write([]byte(defaultConfig()))
	if err != nil {
		return err
	}
	return nil
}

func defaultConfig() string {

	return `## serverHost: адрес поднимаемого сервера для сбора данных
## dataIsUpToDateForSeconds: время хранения кэша в секундах

## filenameSMSData: путь к файлу данных о сервисе СМС
## filenameVoiceData: путь к файлу данных о сервисе ГолосовыхУслуг
## filenameEmailData: путь к файлу данных о сервисе ЭлектроннойПочты
## filenameBillingData: путь к файлу данных о сервисе Оплаты
## incomingDataHost: адрес сервера данных о сервисах ММС, Поддержки и Инцидентов



serverHost:                 localhost:8080
dataIsUpToDateForSeconds:   30

filenameSMSData:            ./skillbox-diploma/sms.data
filenameVoiceData:          ./skillbox-diploma/voice.data
filenameEmailData:          ./skillbox-diploma/email.data
filenameBillingData:        ./skillbox-diploma/billing.data
incomingDataHost:           127.0.0.1:8383`
}

func readDataConfig(data []byte) map[string]string {

	lines := strings.Split(string(data), "\n")
	records := make(map[string]string)
	for _, line := range lines {
		if line == "\r" || line == "" || line[:1] == "#" {
			continue
		}
		rows := strings.Fields(line)
		if len(rows) != 2 {
			continue
		}
		records[rows[0]] = rows[1]
	}
	return records
}

func checkFilePath(records map[string]string) error {

	var (
		filesPath = []string{
			"filenameSMSData:",
			"filenameVoiceData:",
			"filenameEmailData:",
			"filenameBillingData:"}
	)

	for _, val := range filesPath {
		if records[val] == "" {
			if err := appendConfig(val); err != nil {
				return err
			}
			return fmt.Errorf("not enough data [%s] check \"%s\"\n",
				val[:len(val)-1], fileConfig)
		}
	}
	return nil
}

func appendConfig(s string) error {

	file, err := os.OpenFile(fileConfig, os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if file.Close() != nil {
			log.Println(err)
		}
	}()

	str := strings.Split(defaultConfig(), "\n")
	for _, val := range str {
		if strings.Index(val, s) == 0 {
			_, err = file.Write([]byte("\n" + val))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func checkHost(records map[string]string) error {

	var (
		hosts  = []string{"serverHost:", "incomingDataHost:"}
		msg, i = "", 0
	)

	for j, val := range hosts {
		i = j
		parts := strings.Split(records[val], ":")
		switch len(parts) {
		case 0: // empty
			msg = "not enough data"
			i = j
		case 1, 2: // IPv4
			msg = checkIPv4(parts)
		default: // IPv6
			msg = checkIPv6(records[val])
		}
		if msg != "" {
			break
		}
	}
	if msg == "" {
		return nil
	}
	if err := appendConfig(hosts[i]); err != nil {
		return err
	}
	return fmt.Errorf("%s [%s] check \"%s\"\n",
		msg, hosts[i][:len(hosts[i])-1], fileConfig)
}

func checkIPv4(parts []string) string {

	parseIP := net.ParseIP(parts[0])
	if parseIP == nil && strings.ToLower(parts[0]) != "localhost" {
		return "incorrect IPv4"
	}

	if len(parts) == 2 {
		port, err := strconv.Atoi(parts[1])
		if err != nil || 1<<10 >= port || port >= 1<<16 {
			return "incorrect port"
		}
	}
	return ""
}

func checkIPv6(s string) string {

	if !strings.ContainsAny(s, "[]") {
		parseIP := net.ParseIP(s)
		if parseIP == nil {
			return "incorrect IPv6"
		}
		return ""
	}

	port, err := strconv.Atoi(s[strings.LastIndex(s, ":")+1:])
	if err != nil || 1<<10 >= port || port >= 1<<16 {
		return "incorrect port"
	}

	parseIP := net.ParseIP(s[1:strings.
		Index(s, "]")])
	if parseIP == nil {
		return "incorrect IPv6"
	}
	return ""
}
