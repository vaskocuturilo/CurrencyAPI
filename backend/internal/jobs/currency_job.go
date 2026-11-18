package jobs

import (
	"CurrencyAPI/internal/config"
	"CurrencyAPI/internal/handlers"
	"CurrencyAPI/internal/providers"
	"CurrencyAPI/internal/repository"
	"context"
	"log"
)

type CurrencyUpdateJob struct {
	handler *handlers.CurrencyHandler
	repo    repository.CurrencyRepository
}

func NewCurrencyUpdateJob(handler *handlers.CurrencyHandler, repo repository.CurrencyRepository) *CurrencyUpdateJob {
	return &CurrencyUpdateJob{
		handler: handler,
		repo:    repo,
	}
}

func (j *CurrencyUpdateJob) Name() string {
	return "Currency Update Job"
}

func (j *CurrencyUpdateJob) Schedule() string {
	cfg := config.Load()

	return cfg.Schedule.Time
}

func (j *CurrencyUpdateJob) Run(ctx context.Context) error {
	log.Println("Running scheduled currency update job...")

	newRates := providers.FetchAllCurrencies(ctx, j.repo)
	if len(newRates) == 0 {
		log.Println("[Warning] Cron fetch returned 0 rates. Keeping existing cache.")
		return nil
	}

	j.handler.UpdateCache(newRates)
	log.Printf("Successfully updated cache with %d rates via cron.", len(newRates))
	return nil
}
