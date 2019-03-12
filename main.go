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

	converter := converter.NewConverter(provider.NewExchangeRatesIOProvider(), bases, time.Hour)
	server := http.Server{Converter: converter, Stat: converter}

	logrus.Info("Starting server")
	logrus.WithError(server.Start()).Fatal("Server can't start")
}
