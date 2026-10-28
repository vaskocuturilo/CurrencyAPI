package providers

import (
	"CurrencyAPI/models"
	"encoding/xml"
	"fmt"
	"log/slog"
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
		slog.Error("Failed to parse Banca Națională a României XML payload", "error", err)

		return nil, fmt.Errorf("failed to parse banca national romania XML: %w", err)
	}
	date := dataSet.Body.Cube.Date
	if date == "" {
		date = dataSet.Header.PublishingDate
	}

	var currencies []models.CurrencyResponse
	parseErrors := 0

	for _, r := range dataSet.Body.Cube.Rates {
		valStr := strings.TrimSpace(r.Value)
		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			parseErrors++
			slog.Debug("Skipping malformed rate entry from BNR", "currency", r.Currency, "value", valStr, "error", err)
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

	if parseErrors > 0 {
		slog.Warn("Encountered unparseable rates while processing BNR data", "skipped_count", parseErrors)
	}

	return currencies, err
}
