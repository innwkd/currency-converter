package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type CurrencyRate struct {
	Pair CurrencyPair    `json:"pair"`
	Rate decimal.Decimal `json:"rate"`

	Provider  string    `json:"provider"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RateProvider interface {
	GetRate(pair CurrencyPair) (CurrencyRate, error)
}
