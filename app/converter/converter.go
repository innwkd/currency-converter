package converter

import (
	"sync"
	"time"

	"currency-converter/app/storage"
	"currency-converter/app/types"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type converter struct {
	provider      types.RateProvider
	bases         []types.CurrencyPair
	rateStorage   types.RateStorage
	cacheDuration time.Duration

	mu sync.RWMutex
}

func NewConverter(
	provider types.RateProvider,
	bases []types.CurrencyPair,
	rateStorage types.RateStorage,
	cacheDuration time.Duration,
) *converter {
	return &converter{
		provider:      provider,
		bases:         bases,
		rateStorage:   rateStorage,
		cacheDuration: cacheDuration,
	}
}

func (c *converter) Convert(pair types.CurrencyPair, amount decimal.Decimal) (types.Conversion, error) {
	if !c.baseAllowed(pair) {
		return types.Conversion{}, types.ConverterError("base is not allowed")
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
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.bases
}

func (c *converter) CacheDuration() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cacheDuration
}

func (c *converter) baseAllowed(pair types.CurrencyPair) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, base := range c.bases {
		if base == pair {
			return true
		}
	}

	return false
}
