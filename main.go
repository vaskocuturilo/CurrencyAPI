package main

import (
	"CurrencyAPI/cron"
	"CurrencyAPI/handlers"
	"CurrencyAPI/internal/config"
	"CurrencyAPI/jobs"
	"CurrencyAPI/migrations"
	"CurrencyAPI/models"
	"CurrencyAPI/providers"
	"CurrencyAPI/repository"
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.Postgres.ConnString())

	if err != nil {
		slog.Error("Failed to init Postgres: ", "error", err)
		os.Exit(1)
	}

	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.DBTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		slog.Error("Postgres is not ready: ", "error", err)
		os.Exit(1)
	} else {
		slog.Info("Postgres is ready")
	}

	err = migrations.RunMigrations(db)

	if err != nil {
		slog.Error("Failed to init migration process: ", "error", err)
		os.Exit(1)
	}

	slog.Info("Successfully init migration process")

	currencyRepo := repository.NewPostgresCurrencyRepository(db)

	slog.Info("Fetching initial daily exchange rates ...")

	fetchCtx, fetchCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer fetchCancel()

	initialRates := providers.FetchAllCurrencies(fetchCtx, currencyRepo)

	if len(initialRates) == 0 {
		slog.Warn("Failed to fetch rates from external APIs on startup. Attempting to load from database fallback...")
		dbRates, dbErr := currencyRepo.GetAll(fetchCtx)
		if dbErr == nil && len(dbRates) > 0 {
			initialRates = make([]models.CurrencyResponse, len(dbRates))
			for i, r := range dbRates {
				initialRates[i] = *r
			}
			slog.Info("Successfully loaded rates from database fallback", "count", len(initialRates))
		} else {
			slog.Warn("Initialization warning: No rates available from APIs or database yet")
		}
	} else {
		slog.Info("Successfully loaded and saved currency entries", "total_currencies", len(initialRates))
	}

	currencyHandler := handlers.NewCurrencyHandler(initialRates)

	ctxSignal, cancelSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancelSignal()

	jobManager := cron.NewJobManager(ctxSignal)

	jobManager.RegisterJob(jobs.NewCurrencyUpdateJob(currencyHandler, currencyRepo))

	go jobManager.StartScheduler()

	router := gin.Default()

	router.Use(handlers.CORSMiddleware())

	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	api := router.Group("/api/v1")

	api.GET("/currencies", currencyHandler.GetCurrencies)

	srv := &http.Server{Addr: net.JoinHostPort(cfg.Server.Host, cfg.Server.Port), Handler: router}

	go func() {
		slog.Info("Server is starting on port", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Warn("Listen error", "error", err)
		}
	}()

	<-ctxSignal.Done()
	slog.Info("Shutting down gracefully ...")

	jobManager.Stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Info("Server shutdown error", "error", err)
	}

	slog.Info("Server stopped successfully.")
}
