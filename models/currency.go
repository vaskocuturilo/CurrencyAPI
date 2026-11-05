package models

import "errors"

type CurrencyResponse struct {
	Date  string  `json:"date"`
	Base  string  `json:"base"`
	Quote string  `json:"quote"`
	Rate  float64 `json:"rate"`
}

var (
	ErrNotFound   = errors.New("currency not found")
	InternalError = "Internal error"
)
