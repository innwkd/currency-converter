package types

import "github.com/shopspring/decimal"

// Currency represents currency in ISO 4217 format
type Currency string

type CurrencyPair struct {
	From Currency `json:"from"`
	To   Currency `json:"to"`
}

// Conversion represents one conversion result
type Conversion struct {
	Amount decimal.Decimal `json:"amount"`
	Result decimal.Decimal `json:"result"`

	Pair CurrencyPair `json:"currency"`
	Rate CurrencyRate `json:"rate"`
}

type Converter interface {
	Convert(from, to Currency, amount decimal.Decimal) (Conversion, error)
}
