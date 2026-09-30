package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

const ecbURL = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

type CurrencyResponse struct {
	Date  string  `json:"date"`
	Base  string  `json:"base"`
	Quote string  `json:"quote"`
	Rate  float64 `json:"rate"`
}

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

var cachedCurrencies []CurrencyResponse

func fetchAndParseECBRates() ([]CurrencyResponse, error) {
	resp, err := http.Get(ecbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ECB data: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("failed to fetch body data: %w", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ECB API returned status: %s", resp.Status)
	}

	var env XMLEnvelope
	if err := xml.NewDecoder(resp.Body).Decode(&env); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	date := env.Cube.TimeCube.Time
	var currencies []CurrencyResponse

	for _, item := range env.Cube.TimeCube.Rates {
		rateVal, err := strconv.ParseFloat(item.Rate, 64)
		if err != nil {
			log.Printf("Warning: failed to parse rate for %s: %v", item.Currency, err)
			continue
		}

		currencies = append(currencies, CurrencyResponse{
			Date:  date,
			Base:  "EUR",
			Quote: item.Currency,
			Rate:  rateVal,
		})
	}

	return currencies, nil
}

func currenciesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cachedCurrencies); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func main() {
	var err error
	log.Println("Fetching daily exchange rates from ECB...")

	cachedCurrencies, err = fetchAndParseECBRates()

	if err != nil {
		log.Fatalf("Initialization error: %v", err)
	}

	log.Printf("Successfully loaded %d currencies for date %s", len(cachedCurrencies), cachedCurrencies[0].Date)

	http.HandleFunc("/currencies", currenciesHandler)

	port := ":8080"

	log.Printf("Microservice running on http://localhost%s/currencies", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
