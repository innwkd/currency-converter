package main

import (
	"time"

	"github.com/yddmat/currency-converter/app/api/http"
	"github.com/yddmat/currency-converter/app/converter"
	"github.com/yddmat/currency-converter/app/provider"
	"github.com/yddmat/currency-converter/app/types"
	"github.com/tarent/logrus"
)

func main() {
	bases := []types.CurrencyPair{
		{Base: "EUR", To: "USD"},
		{Base: "USD", To: "EUR"},
	}

	cacheDuration := time.Hour

	server := http.Server{
		Bases:         bases,
		Converter:     converter.NewConverter(provider.NewExchangeRatesIoProvider(), bases, cacheDuration),
		CacheDuration: cacheDuration,
	}

	logrus.Info("Starting server")
	logrus.WithError(server.Start()).Fatal("Server can't start")
}
