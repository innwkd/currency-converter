package converter

import (
	"sync"

	"fmt"

	"github.com/pkg/errors"
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
	var err error
	c.mu.Lock()
	defer c.mu.Unlock()

	// TODO check is rate allowed

	currency, exist := c.getRate(pair)
	if !exist {
		currency, err = c.updateRate(pair)
		if err != nil {
			// log
			return types.Conversion{}, errors.Wrapf(err, "can't update rate")
		}
	}

	return types.Conversion{
		Amount: amount,
		Result: amount.Mul(currency.Rate),
		Pair:   pair,
		Rate:   currency,
	}, nil
}

func (c *Converter) getRate(pair types.CurrencyPair) (types.CurrencyRate, bool) {
	rate, exists := c.ratesCache[pair]
	return rate, exists
}

func (c *Converter) updateRate(pair types.CurrencyPair) (types.CurrencyRate, error) {
	for _, provider := range c.providers {
		rate, err := provider.GetRate(pair)
		if err == nil {
			fmt.Println("pair updated")
			c.ratesCache[pair] = rate
			return rate, nil
		}

		// log warning
	}

	return types.CurrencyRate{}, errors.New("can't update rate")
}
