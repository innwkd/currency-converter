package converter

import (
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"github.com/yddmat/currency-converter/app/types"
	"github.com/pkg/errors"
)

type Converter struct {
	ratesCache    map[types.CurrencyPair]types.CurrencyRate
	cacheDuration time.Duration
	provider      types.RateProvider
	bases         []types.CurrencyPair

	mu sync.RWMutex
}

func NewConverter(provider types.RateProvider, bases []types.CurrencyPair, cacheDuration time.Duration) *Converter {
	return &Converter{
		ratesCache:    make(map[types.CurrencyPair]types.CurrencyRate),
		cacheDuration: cacheDuration,
		provider:      provider,
		bases:         bases,
	}
}

func (c *Converter) Convert(pair types.CurrencyPair, amount decimal.Decimal) (types.Conversion, error) {
	var err error
	c.mu.RLock()

	if !c.baseAllowed(pair) {
		return types.Conversion{}, types.ConverterError("base is not allowed")
	}

	rate, cached := c.ratesCache[pair]
	c.mu.RUnlock()

	if !cached || rate.UpdatedAt.Add(c.cacheDuration).Before(time.Now()) {
		rate, err = c.provider.GetRate(pair)
		if err != nil {
			return types.Conversion{}, errors.Wrapf(err, "can't update missing rate")
		}

		c.mu.Lock()
		c.ratesCache[pair] = rate
		c.mu.Unlock()
	}

	return types.Conversion{
		Result:       amount.Mul(rate.Value),
		CurrencyRate: rate,
	}, nil
}

func (c *Converter) CachedRates() []types.CurrencyRate {
	c.mu.RLock()
	defer c.mu.RUnlock()

	rates := make([]types.CurrencyRate, 0)
	for _, cache := range c.ratesCache {
		rates = append(rates, cache)
	}

	return rates
}

func (c *Converter) baseAllowed(pair types.CurrencyPair) bool {
	for _, base := range c.bases {
		if base == pair {
			return true
		}
	}

	return false
}
