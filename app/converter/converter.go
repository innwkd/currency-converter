package converter

import (
	"time"

	"currency-converter/app/storage"
	"currency-converter/app/types"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type converter struct {
	provider      types.RateProvider
	pairs         []types.CurrencyPair
	rateStorage   types.RateStorage
	cacheDuration time.Duration
}

func NewConverter(
	provider types.RateProvider,
	pairs []types.CurrencyPair,
	rateStorage types.RateStorage,
	cacheDuration time.Duration,
) *converter {
	if len(pairs) < 1 {
		panic("need at least 1 currency pair")
	}

	return &converter{
		provider:      provider,
		pairs:         pairs,
		rateStorage:   rateStorage,
		cacheDuration: cacheDuration,
	}
}

var BaseIsNotAllowed = types.ConverterError("base is not allowed")

func (c *converter) Convert(pair types.CurrencyPair, amount decimal.Decimal) (types.Conversion, error) {
	if !c.baseAllowed(pair) {
		return types.Conversion{}, BaseIsNotAllowed
	}

	rate, err := c.rateStorage.Get(pair)
	if err != nil {
		if !storage.IsNotExists(err) {
			return types.Conversion{}, errors.Wrapf(err, "can't get info about rate from storage")
		}

		rate, err = c.provider.GetRate(pair)
		if err != nil {
			return types.Conversion{}, errors.Wrapf(err, "can't get rate from provider")
		}

		rate.Provider = c.provider.Name()
		rate, err = c.rateStorage.Set(pair, rate, c.cacheDuration)
		if err != nil {
			return types.Conversion{}, errors.Wrapf(err, "can't save new rate")
		}
	}

	return types.Conversion{
		Result:       amount.Mul(rate.Value),
		CurrencyRate: rate,
	}, nil
}

func (c *converter) CachedRates() ([]types.CurrencyRate, error) {
	return c.rateStorage.GetAll()
}

func (c *converter) AllowedPair() []types.CurrencyPair {
	return c.pairs
}

func (c *converter) CacheDuration() time.Duration {
	return c.cacheDuration
}

func (c *converter) baseAllowed(pair types.CurrencyPair) bool {
	for _, base := range c.pairs {
		if base == pair {
			return true
		}
	}

	return false
}
