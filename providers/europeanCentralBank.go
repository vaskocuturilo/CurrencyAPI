package providers

import (
	"CurrencyAPI/models"
	"encoding/xml"
	"fmt"
	"log/slog"
	"strconv"
)

const europeanCentralBankURL = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

type XMLRateCube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}

type XMLTimeCube struct {
	Time  string        `xml:"time,attr"`
	Rates []XMLRateCube `xml:"Cube"`
}

type XMLOuterCube struct {
	TimeCube XMLTimeCube `xml:"Cube"`
}

type XMLEnvelope struct {
	XMLName xml.Name     `xml:"Envelope"`
	Cube    XMLOuterCube `xml:"Cube"`
}

func FetchEuropeanCentraBankRates() ([]models.CurrencyResponse, error) {
	headers := map[string]string{
		"Accept": "application/xml, text/xml, */*",
	}

	body, err := fetchURL(europeanCentralBankURL, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch European Central Bank data: %w", err)
	}

	var env XMLEnvelope
	if err := xml.Unmarshal(body, &env); err != nil {
		slog.Error("Failed to parse European Central Bank XML payload", "error", err)
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	date := env.Cube.TimeCube.Time
	var currencies []models.CurrencyResponse
	parseErrors := 0

	for _, item := range env.Cube.TimeCube.Rates {
		rateVal, err := strconv.ParseFloat(item.Rate, 64)
		if err != nil {
			parseErrors++
			slog.Warn("Failed to parse rate for currency from ECB", "currency", item.Currency, "raw_rate", item.Rate, "error", err)
			continue
		}

		currencies = append(currencies, models.CurrencyResponse{
			Date:  date,
			Base:  "EUR",
			Quote: item.Currency,
			Rate:  rateVal,
		})
	}

	if parseErrors > 0 {
		slog.Warn("Encountered unparseable rates in ECB dataset", "skipped_count", parseErrors)
	}

	return currencies, nil
}
