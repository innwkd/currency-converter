//go:generate sh -c "mockery -inpkg -name RateProvider -print > file.tmp && mv file.tmp rate_provider_mock.go"
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
	Name() string
	GetRate(pair CurrencyPair) (CurrencyRate, error)
}
