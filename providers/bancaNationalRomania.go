package providers

import (
	"CurrencyAPI/models"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

var bancaNationalRomaniaURL = "https://curs.bnr.ro/nbrfxrates.xml"

type BNRDataSet struct {
	XMLName xml.Name  `xml:"DataSet"`
	Header  BNRHeader `xml:"Header"`
	Body    BNRBody   `xml:"Body"`
}

type BNRHeader struct {
	PublishingDate string `xml:"PublishingDate"`
}

type BNRBody struct {
	Cube BNRCube `xml:"Cube"`
}

type BNRCube struct {
	Date  string    `xml:"date,attr"`
	Rates []BNRRate `xml:"Rate"`
}

type BNRRate struct {
	Currency   string `xml:"currency,attr"`
	Multiplier int    `xml:"multiplier,attr"`
	Value      string `xml:",innerxml"`
}

func FetchBancaNationalRomaniaRates() ([]models.CurrencyResponse, error) {
	headers := map[string]string{
		"Accept": "application/xml, text/xml, */*",
	}

	body, err := fetchURL(bancaNationalRomaniaURL, headers)

	if err != nil {
		return nil, fmt.Errorf("banca national romania network error: %w", err)
	}

	var dataSet BNRDataSet

	if err := xml.Unmarshal(body, &dataSet); err != nil {
		return nil, fmt.Errorf("failed to parse banca national romania XML: %w", err)
	}
	var currencies []models.CurrencyResponse

	date := dataSet.Body.Cube.Date
	if date == "" {
		date = dataSet.Header.PublishingDate
	}

	for _, r := range dataSet.Body.Cube.Rates {
		valStr := strings.TrimSpace(r.Value)
		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			continue
		}

		rate := val
		if r.Multiplier > 1 {
			rate = val / float64(r.Multiplier)
		}

		currencies = append(currencies, models.CurrencyResponse{
			Date:  date,
			Base:  "RON",
			Quote: r.Currency,
			Rate:  rate,
		})
	}

	return currencies, err
}
