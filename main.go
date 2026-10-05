package main

import (
	"CurrencyAPI/cron"
	"CurrencyAPI/handlers"
	"CurrencyAPI/jobs"
	"CurrencyAPI/providers"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Fetching initial daily exchange rates ...")

	initialRates := providers.FetchAllCurrencies()

	if len(initialRates) == 0 {
		log.Fatalf("Initialization error: Failed to load rates from banks")
	}

	log.Printf("Successfully loaded %d total currency entries into cache", len(initialRates))

	currencyHandler := handlers.NewCurrencyHandler(initialRates)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	jobManager := cron.NewJobManager(ctx)

	jobManager.RegisterJob(jobs.NewCurrencyUpdateJob(currencyHandler))

	go jobManager.StartScheduler()

	http.HandleFunc("/currencies", handlers.CORSMiddleware(currencyHandler.GetCurrencies))

	port := ":8080"

	server := &http.Server{Addr: port}

	go func() {
		log.Printf("Microservice running on http://localhost%s/currencies", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down gracefully ...")

	jobManager.Stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped successfully.")
}
