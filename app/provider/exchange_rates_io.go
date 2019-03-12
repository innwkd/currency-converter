package provider

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/tarent/logrus"
	"github.com/yddmat/currency-converter/app/types"
)

type ExchangeRatesIo struct {
	domain *url.URL
	name   string
}

func NewExchangeRatesIoProvider() *ExchangeRatesIo {
	domain, err := url.Parse("https://api.exchangeratesapi.io")
	if err != nil {
		logrus.WithError(err).Fatal("parsing url")
	}

	return &ExchangeRatesIo{domain: domain, name: "exchangeratesapi.io"}
}

type getRateResponse struct {
	Rates map[string]decimal.Decimal `json:"rates"`
}

func (e *ExchangeRatesIo) GetRate(pair types.CurrencyPair) (types.CurrencyRate, error) {
	resp := getRateResponse{}
	r, _, err := gorequest.
		New().
		Get(fmt.Sprintf(e.domain.String()+"/latest?base=%s&symbols=%s", pair.Base, pair.To)).
		EndStruct(&resp)

	if err != nil {
		return types.CurrencyRate{}, errors.Wrapf(err[0], "can't send request")
	}

	if r.StatusCode != http.StatusOK {
		return types.CurrencyRate{}, errors.Errorf("http request return %d response code", r.StatusCode)
	}

	if len(resp.Rates) != 1 {
		return types.CurrencyRate{}, errors.New("provider response format isn't valid")
	}

	rate, ok := resp.Rates[string(pair.To)]
	if !ok {
		return types.CurrencyRate{}, errors.New("can't parse rate response")
	}

	return types.CurrencyRate{
		Pair:      pair,
		Value:     rate,
		Provider:  e.name,
		UpdatedAt: time.Now(),
	}, nil
}
