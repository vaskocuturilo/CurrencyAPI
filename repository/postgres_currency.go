package repository

import (
	"CurrencyAPI/models"
	"context"
	"database/sql"
	"errors"
)

type PostgresCurrencyRepository struct {
	db *sql.DB
}

func NewPostgresCurrencyRepository(db *sql.DB) *PostgresCurrencyRepository {
	return &PostgresCurrencyRepository{db: db}
}

func (r *PostgresCurrencyRepository) Get(ctx context.Context, id string) (*models.CurrencyResponse, error) {
	query := `SELECT id, date, base, quote,rate 
              FROM currencies_db WHERE id = $1`

	currency := &models.CurrencyResponse{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&currency.Date, &currency.Base, &currency.Quote, &currency.Rate)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrNotFound
	}

	return currency, err
}

func (r *PostgresCurrencyRepository) GetAll(ctx context.Context) ([]*models.CurrencyResponse, error) {
	events := make([]*models.CurrencyResponse, 0)

	query := `SELECT * FROM currencies_db`

	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		c := &models.CurrencyResponse{}
		if err := rows.Scan(&c.Base, &c.Rate, &c.Quote, &c.Rate); err != nil {
			return nil, err
		}
		events = append(events, c)
	}

	return events, err
}
