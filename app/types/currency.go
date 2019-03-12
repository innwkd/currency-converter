package types

import (
	"github.com/shopspring/decimal"
	"fmt"
	"time"
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
	CachedRates() []CurrencyRate
	AllowedBases() []CurrencyPair
	CacheDuration() time.Duration
}
