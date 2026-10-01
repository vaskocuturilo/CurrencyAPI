package main

import (
	"CurrencyAPI/handlers"
	"CurrencyAPI/providers"
	"log"
	"net/http"
)

func main() {
	log.Println("Fetching daily exchange rates ...")

	cachedCurrencies := providers.FetchAllCurrencies()

	if len(cachedCurrencies) == 0 {
		log.Fatalf("Initialization error: Failed to load rates from banks")
	}

	log.Printf("Successfully loaded %d total currency entries into cache", len(cachedCurrencies))

	currencyHandler := handlers.NewCurrencyHandler(cachedCurrencies)

	http.HandleFunc("/currencies", currencyHandler.GetCurrencies)

	port := ":8080"
	log.Printf("Microservice running on http://localhost%s/currencies", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
