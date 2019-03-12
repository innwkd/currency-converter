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

	CurrencyRate CurrencyRate `json:"currency_rate"`
}

type Converter interface {
	Convert(pair CurrencyPair, amount decimal.Decimal) (Conversion, error)
}
