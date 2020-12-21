//go:generate sh -c "mockery -inpkg -name RateStorage -print > file.tmp && mv file.tmp rate_storage_mock.go"
package types

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// Currency in ISO 4217 format
type Currency string

type CurrencyPair struct {
	Base Currency `json:"from"`
	To   Currency `json:"to"`
}

func (cp CurrencyPair) String() string {
	return fmt.Sprintf("%s_%s", cp.Base, cp.To)
}

// Conversion represents one conversion result
type Conversion struct {
	Result       decimal.Decimal `json:"result"`
	CurrencyRate CurrencyRate    `json:"currency_rate"`
}

type Converter interface {
	Convert(pair CurrencyPair, amount decimal.Decimal) (Conversion, error)
}

type ConverterStat interface {
	CachedRates() ([]CurrencyRate, error)
	AllowedPair() []CurrencyPair
	CacheDuration() time.Duration
}

type RateStorage interface {
	// Set currency rate for specified pair
	// If pair already have value it will be returned without updating
	Set(pair CurrencyPair, rate CurrencyRate, duration time.Duration) (CurrencyRate, error)

	// Get currency rate
	Get(pair CurrencyPair) (CurrencyRate, error)

	// GetAll is returning all stored rates
	GetAll() ([]CurrencyRate, error)
}
