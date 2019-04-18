package converter

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/yddmat/currency-converter/app/storage"
	"github.com/yddmat/currency-converter/app/types"
)

type Converter struct {
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
) *Converter {
	return &Converter{
		provider:      provider,
		bases:         bases,
		rateStorage:   rateStorage,
		cacheDuration: cacheDuration,
	}
}

func (c *Converter) Convert(pair types.CurrencyPair, amount decimal.Decimal) (types.Conversion, error) {
	var err error
	c.mu.RLock()

	if !c.baseAllowed(pair) {
		return types.Conversion{}, types.ConverterError("base is not allowed")
	}

	rate, err := c.rateStorage.Get(pair)
	if err != nil {
		if !storage.IsNotExists(err) {
			return types.Conversion{}, errors.Wrapf(err, "can't get info about rate from storage")
		} else {
			rate, err = c.provider.GetRate(pair)
			if err != nil {
				return types.Conversion{}, errors.Wrapf(err, "can't get rate from provider")
			}

			rate, err = c.rateStorage.Set(pair, rate, c.cacheDuration)
			if err != nil {
				return types.Conversion{}, errors.Wrapf(err, "can't save new rate")
			}
		}
	}

	return types.Conversion{
		Result:       amount.Mul(rate.Value),
		CurrencyRate: rate,
	}, nil
}

func (c *Converter) CachedRates() ([]types.CurrencyRate, error) {
	return c.rateStorage.GetAll()
}

func (c *Converter) AllowedPair() []types.CurrencyPair {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.bases
}

func (c *Converter) CacheDuration() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cacheDuration
}

func (c *Converter) baseAllowed(pair types.CurrencyPair) bool {
	for _, base := range c.bases {
		if base == pair {
			return true
		}
	}

	return false
}
