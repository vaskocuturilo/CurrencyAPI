package providers

import (
	"CurrencyAPI/internal/models"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

const bisDefaultBaseURL = "https://stats.bis.org/api/v2/data/dataflow/BIS/WS_XRU/1.0/D..USD.A?startPeriod=%s&format=csv"

func FetchBankInternationalSettlements() ([]models.CurrencyResponse, error) {
	startPeriod := time.Now().AddDate(0, 0, -14).Format("2006-01-02")
	url := fmt.Sprintf(bisDefaultBaseURL, startPeriod)

	headers := map[string]string{
		"Accept": "application/vnd.sdmx.data+csv;version=1.0.0",
	}

	body, err := fetchURL(url, headers)
	if err != nil {
		return nil, fmt.Errorf("fetch BIS rates: %w", err)
	}

	reader := csv.NewReader(bytes.NewReader(body))
	headersRow, err := reader.Read()
	if err != nil {
		slog.Error("Failed to parse BIS CSV header", "error", err)
		return nil, fmt.Errorf("parse BIS CSV header: %w", err)
	}

	colIndex := make(map[string]int, len(headersRow))
	for i, col := range headersRow {
		colIndex[strings.TrimSpace(col)] = i
	}

	dateIdx, okDate := colIndex["TIME_PERIOD"]
	rateIdx, okRate := colIndex["OBS_VALUE"]
	if !okDate || !okRate {
		return nil, fmt.Errorf("BIS CSV missing TIME_PERIOD or OBS_VALUE columns")
	}

	quoteIdx, okQuote := colIndex["CURRENCY"]
	if !okQuote {
		quoteIdx, okQuote = colIndex["REF_AREA"]
	}
	if !okQuote {
		return nil, fmt.Errorf("BIS CSV missing CURRENCY or REF_AREA column")
	}

	var currencies []models.CurrencyResponse
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Warn("Skipping malformed BIS CSV row", "error", err)
			continue
		}
		if len(row) <= dateIdx || len(row) <= rateIdx || len(row) <= quoteIdx {
			continue
		}

		quote := strings.ToUpper(strings.TrimSpace(row[quoteIdx]))
		if models.Validate(quote) != nil {
			continue
		}

		rate, err := strconv.ParseFloat(strings.TrimSpace(row[rateIdx]), 64)
		if err != nil {
			continue
		}

		currencies = append(currencies, models.CurrencyResponse{
			Date:  strings.TrimSpace(row[dateIdx]),
			Base:  "USD",
			Quote: quote,
			Rate:  rate,
		})
	}

	if len(currencies) == 0 {
		return nil, fmt.Errorf("no valid currency rates found in BIS response")
	}

	return currencies, nil
}
