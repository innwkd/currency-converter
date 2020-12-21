package provider

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"currency-converter/app/types"

	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type exchangeRatesIOProvider struct {
	domain *url.URL
	name   string
}

func NewExchangeRatesIOProvider() *exchangeRatesIOProvider {
	domain, err := url.Parse("https://api.exchangeratesapi.io")
	if err != nil {
		logrus.WithError(err).Fatal("parsing url")
	}

	return &exchangeRatesIOProvider{domain: domain, name: "exchangeratesapi.io"}
}

type getRateResponse struct {
	Rates map[string]decimal.Decimal `json:"rates"`
}

func (e *exchangeRatesIOProvider) Name() string {
	return e.name
}

func (e *exchangeRatesIOProvider) GetRate(pair types.CurrencyPair) (types.CurrencyRate, error) {
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
		UpdatedAt: time.Now(),
	}, nil
}
