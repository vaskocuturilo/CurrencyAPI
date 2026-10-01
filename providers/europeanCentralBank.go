package providers

import (
	"CurrencyAPI/models"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
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
	response, err := http.Get(europeanCentralBankURL)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch European Central Bank data: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to fetch body European Central Bank: %w", err)
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("European Central Bank API returned status: %s", response.Status)
	}

	var env XMLEnvelope
	if err := xml.NewDecoder(response.Body).Decode(&env); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	date := env.Cube.TimeCube.Time
	var currencies []models.CurrencyResponse

	for _, item := range env.Cube.TimeCube.Rates {
		rateVal, err := strconv.ParseFloat(item.Rate, 64)
		if err != nil {
			log.Printf("[European Central Bank Warning] failed to parse rate for %s: %v", item.Currency, err)
			continue
		}

		currencies = append(currencies, models.CurrencyResponse{
			Date:  date,
			Base:  "EUR",
			Quote: item.Currency,
			Rate:  rateVal,
		})
	}

	return currencies, nil
}
