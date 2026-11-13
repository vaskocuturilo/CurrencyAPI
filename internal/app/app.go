package app

import (
	"CurrencyAPI/internal/config"
	"CurrencyAPI/internal/controller"
	"CurrencyAPI/internal/cron"
	"CurrencyAPI/internal/handlers"
	"CurrencyAPI/internal/jobs"
	"CurrencyAPI/internal/models"
	"CurrencyAPI/internal/providers"
	"CurrencyAPI/internal/repository"
	"CurrencyAPI/internal/service"
	"CurrencyAPI/migrations"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type App struct {
	cfg        *config.Config
	db         *sql.DB
	redis      *redis.Client
	srv        *http.Server
	jobManager *cron.JobManager
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	db, err := openPostgres(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := migrations.RunMigrations(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrations: %w", err)
	}
	slog.Info("Migrations applied")

	rdb := newRedis(ctx, cfg)

	repo := repository.NewPostgresCurrencyRepository(db)
	initialRates := loadInitialRates(ctx, repo)

	currencyHandler := handlers.NewCurrencyHandler(initialRates)

	jm := cron.NewJobManager(ctx)
	jm.RegisterJob(jobs.NewCurrencyUpdateJob(currencyHandler, repo))

	ctr := controller.NewUserController(service.NewCurrencyService(repo))

	return &App{
		cfg:        cfg,
		db:         db,
		redis:      rdb,
		jobManager: jm,
		srv: &http.Server{
			Addr:    net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
			Handler: newRouter(ctr),
		},
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	go a.jobManager.StartScheduler()

	errCh := make(chan error, 1)
	go func() {
		slog.Info("Server is starting", "addr", a.srv.Addr)
		if err := a.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		a.jobManager.Stop()
		return fmt.Errorf("http server: %w", err)
	case <-ctx.Done():
		slog.Info("Shutting down gracefully ...")
	}

	a.jobManager.Stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // move to cfg.Server.ShutdownTimeout
	defer cancel()

	if err := a.srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}
	slog.Info("Server stopped successfully")
	return nil
}

func (a *App) Close() {
	if err := a.db.Close(); err != nil {
		slog.Warn("closing postgres", "error", err)
	}
	if err := a.redis.Close(); err != nil {
		slog.Warn("closing redis", "error", err)
	}
}

func openPostgres(ctx context.Context, cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.Postgres.ConnString())
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, cfg.Server.DBTimeout)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	slog.Info("Postgres is ready")
	return db, nil
}

func newRedis(ctx context.Context, cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           0,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	})

	pingCtx, cancel := context.WithTimeout(ctx, cfg.Server.RedisTimeout)
	defer cancel()

	if err := rdb.Ping(pingCtx).Err(); err != nil {
		slog.Warn("Redis is not reachable, working without cache", "error", err)
	} else {
		slog.Info("Connected to Redis")
	}
	return rdb
}

func loadInitialRates(ctx context.Context, repo repository.CurrencyRepository) []models.CurrencyResponse {
	fetchCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	slog.Info("Fetching initial daily exchange rates ...")
	rates := providers.FetchAllCurrencies(fetchCtx, repo)
	if len(rates) > 0 {
		slog.Info("Loaded and saved currency entries", "total", len(rates))
		return rates
	}

	slog.Warn("External APIs failed, falling back to database")
	dbRates, err := repo.GetAll(fetchCtx)
	if err != nil || len(dbRates) == 0 {
		slog.Warn("No rates available from APIs or database yet", "error", err)
		return nil
	}

	rates = make([]models.CurrencyResponse, len(dbRates))
	for i, r := range dbRates {
		rates[i] = *r
	}
	slog.Info("Loaded rates from database fallback", "count", len(rates))
	return rates
}

func newRouter(ctr *controller.CurrencyController) *gin.Engine {
	rl := tollbooth.NewLimiter(20, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Second})
	rl.SetMessage("Too many requests, try again later.")

	router := gin.Default()
	router.Use(handlers.CORSMiddleware())
	router.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	api := router.Group("/api/v1")
	api.GET("/currencies", tollbooth_gin.LimitHandler(rl), ctr.Route)

	return router
}
