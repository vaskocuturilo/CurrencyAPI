package providers

import (
	"CurrencyAPI/models"
	"log/slog"
	"sync"
	"time"
)

type FetchFunc func() ([]models.CurrencyResponse, error)
type Provider struct {
	Name  string
	Fetch FetchFunc
}

func GetProviders() []Provider {
	return []Provider{
		{Name: "European Central Bank", Fetch: FetchEuropeanCentraBankRates},
		{Name: "Banca d'Italia", Fetch: FetchBancaDItaliaRates},
		{Name: "Banca Națională a României", Fetch: FetchBancaNationalRomaniaRates},
	}
}

func FetchAllCurrencies() []models.CurrencyResponse {
	start := time.Now()
	registeredProviders := GetProviders()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var allCurrencies []models.CurrencyResponse

	slog.Info("Starting concurrent fetch from all registered currency providers", "provider_count", len(registeredProviders))

	for _, p := range registeredProviders {
		wg.Add(1)

		go func(provider Provider) {
			defer wg.Done()
			providerStart := time.Now()
			rates, err := provider.Fetch()
			if err != nil {
				slog.Error("Failed fetching currency rates from provider",
					"provider", provider.Name,
					"error", err,
					"latency_ms", time.Since(providerStart).Milliseconds(),
				)
				return
			}
			mu.Lock()
			allCurrencies = append(allCurrencies, rates...)
			mu.Unlock()
			slog.Info("Successfully loaded currencies from provider",
				"provider", provider.Name,
				"count", len(rates),
				"latency_ms", time.Since(providerStart).Milliseconds(),
			)
		}(p)
	}

	wg.Wait()

	slog.Info("Completed currency aggregation across all providers",
		"total_currencies_loaded", len(allCurrencies),
		"total_latency_ms", time.Since(start).Milliseconds(),
	)

	return allCurrencies
}
