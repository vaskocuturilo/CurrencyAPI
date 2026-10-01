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

func FetchAllCurrencies() []models.CurrencyResponse {
	providers := []Provider{{Name: "europeanCentralBankURL", Fetch: FetchEuropeanCentraBankRates}}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var allCurrencies []models.CurrencyResponse

	for _, p := range providers {
		wg.Add(1)

		go func(provider Provider) {
			defer wg.Done()
			rates, err := provider.Fetch()
			if err != nil {
				log.Printf("Error fetching : %v", provider.Name)
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
