package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type CurrencyRate struct {
	Pair  CurrencyPair    `json:"pair"`
	Value decimal.Decimal `json:"value"`

	Provider  string    `json:"provider"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RateProvider interface {
	GetRate(pair CurrencyPair) (CurrencyRate, error)
}
