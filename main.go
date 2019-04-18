package main

import (
	"time"

	"github.com/tarent/logrus"
	"github.com/yddmat/currency-converter/app/api/http/jrpc"
	"github.com/yddmat/currency-converter/app/api/http/rest"
	"github.com/yddmat/currency-converter/app/converter"
	"github.com/yddmat/currency-converter/app/provider"
	"github.com/yddmat/currency-converter/app/storage"
	"github.com/yddmat/currency-converter/app/types"
)

func main() {
	bases := []types.CurrencyPair{
		{Base: "EUR", To: "USD"},
		{Base: "USD", To: "EUR"},
	}

	converter := converter.NewConverter(provider.NewExchangeRatesIOProvider(), bases, storage.NewMemoryStorage(), time.Hour)
	server := rest.Server{Converter: converter, Stat: converter, Port: "8080"}

	go func() {
		logrus.Info("Starting rest server")
		logrus.WithError(server.Start()).Fatal("Rest server can't start")
	}()

	logrus.Info("Starting rpc server")
	rpc := jrpc.Server{Converter: converter, Port: "4444"}
	logrus.WithError(rpc.Start()).Fatal("RPC server can't start")
}
