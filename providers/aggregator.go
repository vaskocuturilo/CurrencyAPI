package providers

import (
	"CurrencyAPI/models"
	"log"
	"sync"
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
	}
}

func FetchAllCurrencies() []models.CurrencyResponse {
	registeredProviders := GetProviders()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var allCurrencies []models.CurrencyResponse

	for _, p := range registeredProviders {
		wg.Add(1)

		go func(provider Provider) {
			defer wg.Done()
			rates, err := provider.Fetch()
			if err != nil {
				log.Printf("[Error] Failed fetching from %s: %v", provider.Name, err)
				return
			}
			mu.Lock()
			allCurrencies = append(allCurrencies, rates...)
			mu.Unlock()
			log.Printf("Loaded %d currencies from European Central Bank", len(rates))
		}(p)
	}

	wg.Wait()
	return allCurrencies
}
