package providers

import (
	"CurrencyAPI/models"
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

const bancaDItaliaLatestURL = "https://tassidicambio.bancaditalia.it/terzevalute-wf-web/rest/v1.0/latestRates?lang=en"

func FetchBancaDItaliaRates() ([]models.CurrencyResponse, error) {
	headers := map[string]string{
		"Accept": "text/csv, text/plain, */*",
	}

	body, err := fetchURL(bancaDItaliaLatestURL, headers)
	if err != nil {
		return nil, fmt.Errorf("Banca d'Italia network error: %w", err)
	}

	r := csv.NewReader(bytes.NewReader(body))
	r.Comment = '#'
	r.FieldsPerRecord = -1

	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse Banca d'Italia CSV payload: %w", err)
	}

	currencyMap := make(map[string]models.CurrencyResponse)

	for i, record := range records {
		if i == 0 || len(record) < 9 {
			continue
		}

		isoCode := strings.TrimSpace(record[2])
		euroRateStr := strings.TrimSpace(record[4])
		refDate := strings.TrimSpace(record[8])

		if isoCode == "" || euroRateStr == "" || euroRateStr == "N.A." {
			continue
		}

		rateVal, err := strconv.ParseFloat(euroRateStr, 64)
		if err != nil || rateVal <= 0 {
			continue
		}

		if _, exists := currencyMap[isoCode]; !exists {
			currencyMap[isoCode] = models.CurrencyResponse{
				Date:  refDate,
				Base:  "EUR",
				Quote: isoCode,
				Rate:  rateVal,
			}
		}
	}

	currencies := make([]models.CurrencyResponse, 0, len(currencyMap))
	for _, currency := range currencyMap {
		currencies = append(currencies, currency)
	}

	return currencies, nil
}
