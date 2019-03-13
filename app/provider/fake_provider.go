package provider

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/yddmat/currency-converter/app/types"
)

type FakeProvider struct {
	Rate       decimal.Decimal
	UpdateTime time.Time
}

func (fp *FakeProvider) Name() string {
	return "fake"
}

func (fp *FakeProvider) GetRate(pair types.CurrencyPair) (types.CurrencyRate, error) {
	return types.CurrencyRate{
		Pair:      pair,
		Value:     fp.Rate,
		UpdatedAt: fp.UpdateTime,
	}, nil
}
