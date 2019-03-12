package converter

import (
	"sync"

	"github.com/shopspring/decimal"
	"github.com/yddmat/currency-converter/app/types"
)

type Converter struct {
	ratesCache map[types.CurrencyPair]types.CurrencyRate
	providers  []types.RateProvider

	// mu protecting rates and providers
	mu sync.RWMutex
}

func NewConverter(providers ...types.RateProvider) *Converter {
	return &Converter{ratesCache: make(map[types.CurrencyPair]types.CurrencyRate), providers: providers}
}

func (c *Converter) Convert(pair types.CurrencyPair, amount decimal.Decimal) (types.Conversion, error) {
	conversion := types.Conversion{
		Amount: amount,
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// TODO check is rate allowed

	var cached, updated bool
	conversion.CurrencyRate, cached = c.getRate(pair)
	if !cached {
		conversion.CurrencyRate, updated = c.updateRate(pair)
		if !updated {
			// log
			return types.Conversion{}, types.ConverterError("error while updating missing rate")
		}
	}

	conversion.Result = amount.Mul(conversion.CurrencyRate.Value)
	return conversion, nil
}

func (c *Converter) getRate(pair types.CurrencyPair) (types.CurrencyRate, bool) {
	rate, exists := c.ratesCache[pair]
	return rate, exists
}

func (c *Converter) updateRate(pair types.CurrencyPair) (types.CurrencyRate, bool) {
	for _, provider := range c.providers {
		rate, err := provider.GetRate(pair)
		if err == nil {
			c.ratesCache[pair] = rate
			return rate, true
		}

		// log warning
	}

	return types.CurrencyRate{}, false
}
