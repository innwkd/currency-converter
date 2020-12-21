package provider

import (
	"time"

	"currency-converter/app/types"

	"github.com/shopspring/decimal"
)

type fakeProvider struct {
	rate       decimal.Decimal
	updateTime time.Time
}

func NewFakeProvider(rate decimal.Decimal, updateTime time.Time) *fakeProvider {
	return &fakeProvider{
		rate:       rate,
		updateTime: updateTime,
	}
}

func (fp *fakeProvider) Name() string {
	return "fake"
}

func (fp *fakeProvider) GetRate(pair types.CurrencyPair) (types.CurrencyRate, error) {
	return types.CurrencyRate{
		Pair:      pair,
		Value:     fp.rate,
		UpdatedAt: fp.updateTime,
	}, nil
}
