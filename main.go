package main

import (
	"CurrencyAPI/cron"
	"CurrencyAPI/handlers"
	"CurrencyAPI/internal/config"
	"CurrencyAPI/jobs"
	"CurrencyAPI/providers"
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	slog.Info("Fetching initial daily exchange rates ...")

	initialRates := providers.FetchAllCurrencies()

	if len(initialRates) == 0 {
		slog.Warn("Initialization error: Failed to load rates from banks")
	}

	slog.Info("Successfully loaded total currency entries into cache", "total currency", len(initialRates))

	currencyHandler := handlers.NewCurrencyHandler(initialRates)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	jobManager := cron.NewJobManager(ctx)

	jobManager.RegisterJob(jobs.NewCurrencyUpdateJob(currencyHandler))

	go jobManager.StartScheduler()

	r := gin.Default()

	api := r.Group("/api/v1")

	api.GET("/currencies", handlers.CORSMiddleware(currencyHandler.GetCurrencies))

	cfg := config.Load()

	srv := &http.Server{Addr: net.JoinHostPort(cfg.Server.Host, cfg.Server.Port), Handler: r}

	go func() {
		slog.Info("Server is starting on port", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Warn("Listen error", "error", err)
		}
	}()

	<-ctx.Done()
	slog.Info("Shutting down gracefully ...")

	jobManager.Stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Info("Server shutdown error", "error", err)
	}

	slog.Info("Server stopped successfully.")
}
