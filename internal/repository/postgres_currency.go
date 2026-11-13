package repository

import (
	"CurrencyAPI/internal/models"
	"context"
	"database/sql"
	"errors"
	"log/slog"
)

type CurrencyRepository interface {
	Get(ctx context.Context, currency string) (*models.CurrencyResponse, error)
	GetAll(ctx context.Context) ([]*models.CurrencyResponse, error)
	SaveAll(ctx context.Context, currencies []models.CurrencyResponse) error
}
type PostgresCurrencyRepository struct {
	db *sql.DB
}

func NewPostgresCurrencyRepository(db *sql.DB) *PostgresCurrencyRepository {
	return &PostgresCurrencyRepository{db: db}
}

func (r *PostgresCurrencyRepository) Get(ctx context.Context, currency string) (*models.CurrencyResponse, error) {
	query := `SELECT date, base, quote,rate FROM currencies WHERE quote = $1`

	c := &models.CurrencyResponse{}

	err := r.db.QueryRowContext(ctx, query, currency).Scan(
		&c.Date, &c.Base, &c.Quote, &c.Rate)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNotFound
	}

	return c, err
}

func (r *PostgresCurrencyRepository) GetAll(ctx context.Context) ([]*models.CurrencyResponse, error) {
	currencies := make([]*models.CurrencyResponse, 0)

	query := `SELECT date, base, quote, rate FROM currencies`

	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		c := &models.CurrencyResponse{}
		if err := rows.Scan(&c.Date, &c.Base, &c.Quote, &c.Rate); err != nil {
			return nil, err
		}
		currencies = append(currencies, c)
	}

	return currencies, err
}

func (r *PostgresCurrencyRepository) SaveAll(ctx context.Context, currencies []models.CurrencyResponse) error {
	if len(currencies) == 0 {
		return nil
	}

	query := `
		INSERT INTO currencies (date, base, quote, rate)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (date, base, quote) 
		DO UPDATE SET rate = EXCLUDED.rate
	`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, c := range currencies {
		_, err := stmt.ExecContext(ctx, c.Date, c.Base, c.Quote, c.Rate)
		if err != nil {
			slog.Error("Failed to save currency record", "quote", c.Quote, "date", c.Date, "error", err)
			return err
		}
	}

	return nil
}
