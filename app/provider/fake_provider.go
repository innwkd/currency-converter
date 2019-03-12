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

func (fp *FakeProvider) GetRate(pair types.CurrencyPair) (types.CurrencyRate, error) {
	return types.CurrencyRate{
		Pair:      pair,
		Rate:      fp.Rate,
		Provider:  "fake",
		UpdatedAt: fp.UpdateTime,
	}, nil
}
