package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func getRowsCSV(filename string) ([]string, error) {

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("access to the file [%s] is denied: %s",
			filename, err)
	}
	return strings.Split(string(data), "\n"), nil
}

func makeEmailMap(e []EmailData) map[string][][]EmailData {

	m := make(map[string][][]EmailData)
	for _, data := range e {
		if m[data.Country] == nil {
			m[data.Country] = make([][]EmailData, 2, 2)
		}
		m[data.Country][0] = append(m[data.Country][0], data)
	}

	for country := range m {
		sort.Slice(m[country][0], func(i, j int) bool {
			return m[country][0][i].DeliveryTime < m[country][0][j].DeliveryTime
		})

		lenM := len(m[country][0])
		m[country][1] = make([]EmailData, lenM, lenM)
		copy(m[country][1], m[country][0])
		count := 3
		if lenM < 3 {
			count = lenM
		}
		m[country][0] = m[country][0][:count]
		m[country][1] = m[country][1][lenM-count:]
	}

	return m
}

func providerNotCorrect(provider, providerType string) (string, error) {

	var (
		correctProvidersVoice = []string{"TransparentCalls", "E-Voice", "JustPhone"}
		correctProvidersSMS   = []string{"Topolo", "Rond", "Kildy"}
		correctProvidersEmail = []string{"Comcast",
			"MSN", "Gmail", "Orange", "RediffMail",
			"AOL", "Live", "Hotmail", "Protonmail",
			"GMX", "Yahoo", "Yandex", "Mail.ru"}
		correctProviders []string
		allowed          = []string{"sms", "mms", "voice", "email"}
	)

	providerType = strings.ToLower(providerType)
	switch providerType {
	case allowed[0]:
		correctProviders = correctProvidersSMS
	case allowed[1]:
		correctProviders = correctProvidersSMS
	case allowed[2]:
		correctProviders = correctProvidersVoice
	case allowed[3]:
		correctProviders = correctProvidersEmail
	default:
		return "", fmt.Errorf("unknown provider type [%s]\n"+
			"allowed: %v\n", strings.ToTitle(providerType), allowed)
	}

	if provider == "" {
		return "", fmt.Errorf("provider is empty\n")
	}
	provider = strings.ToLower(provider)
	for _, name := range correctProviders {
		if provider == strings.ToLower(name) {
			return name, nil
		}
	}
	return "", fmt.Errorf("[%s] is unknown provider\n",
		strings.ToTitle(provider))
}
