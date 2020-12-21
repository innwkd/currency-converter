package converter_test

import (
	"currency-converter/app/converter"
	"currency-converter/app/storage"
	"currency-converter/app/types"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/require"
)

func TestConverter_AllowedPair(t *testing.T) {
	cases := []struct {
		Pairs []types.CurrencyPair

		ExpectPanic bool
	}{
		{
			Pairs:       []types.CurrencyPair{},
			ExpectPanic: true,
		},
		{
			Pairs: []types.CurrencyPair{
				{Base: "USD", To: "EUR"},
				{Base: "EUR", To: "USD"},
			},
		},
		{
			Pairs: []types.CurrencyPair{
				{Base: "RUB", To: "EUR"},
			},
		},
	}

	for _, c := range cases {
		if c.ExpectPanic {
			require.Panics(t, func() {
				_ = converter.NewConverter(nil, c.Pairs, nil, 0)
			})
		} else {
			cnvrt := converter.NewConverter(nil, c.Pairs, nil, 0)
			require.Equal(t, c.Pairs, cnvrt.AllowedPair())
		}
	}
}

func TestConverter_Convert(t *testing.T) {
	requestPair := types.CurrencyPair{Base: "USD", To: "EUR"}

	mockedProviderRate := types.CurrencyRate{
		Pair:      requestPair,
		Value:     decimal.New(15, 0),
		UpdatedAt: time.Now(),
		Provider:  "mock",
	}

	storageMock := &types.MockRateStorage{}
	storageMock.
		On("Get", requestPair).
		Return(types.CurrencyRate{}, storage.ErrNotExists)
	storageMock.
		On("Set", requestPair, mockedProviderRate, time.Hour).
		Return(mockedProviderRate, nil)

	providerMock := &types.MockRateProvider{}
	providerMock.
		On("GetRate", requestPair).
		Return(mockedProviderRate, nil)
	providerMock.
		On("Name").
		Return("mock")

	converter := converter.NewConverter(providerMock, []types.CurrencyPair{requestPair}, storageMock, time.Hour)
	conversion, err := converter.Convert(requestPair, decimal.New(100, 0))
	require.Nil(t, err)
	require.Equal(t, types.Conversion{
		Result:       decimal.New(1500, 0),
		CurrencyRate: mockedProviderRate,
	}, conversion)
}
