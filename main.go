package main

import (
	"time"

	"github.com/go-redis/redis"
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

	con := converter.NewConverter(provider.NewExchangeRatesIOProvider(), bases, storage.NewRedisStorage(redis.NewClient(&redis.Options{Addr: "redis:6379"})), time.Hour)
	server := rest.Server{Converter: con, Stat: con, Port: "8080"}

	go func() {
		logrus.Info("Starting rest server")
		logrus.WithError(server.Start()).Fatal("Rest server can't start")
	}()

	logrus.Info("Starting rpc server")
	rpc := jrpc.Server{Converter: con, Port: "4444"}
	logrus.WithError(rpc.Start()).Fatal("RPC server can't start")
}
