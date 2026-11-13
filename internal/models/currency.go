package models

import (
	"errors"
	"strings"
)

type CurrencyResponse struct {
	Date  string  `json:"date"`
	Base  string  `json:"base"`
	Quote string  `json:"quote"`
	Rate  float64 `json:"rate"`
}

var (
	ErrNotFound              = errors.New("currency not found")
	ErrCurrencyRequired      = errors.New("currency required")
	ErrCurrencyInvalidLength = errors.New("currency code must be exactly 3 characters")
	InternalError            = "Internal error"
)

func Validate(value string) error {
	if strings.TrimSpace(value) == "" {
		return ErrCurrencyRequired
	}

	if len(value) != 3 {
		return ErrCurrencyInvalidLength
	}

	return nil
}
