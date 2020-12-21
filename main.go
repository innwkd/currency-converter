package main

import (
	"currency-converter/app/api/http/jrpc"
	"currency-converter/app/api/http/rest"
	"currency-converter/app/converter"
	"currency-converter/app/provider"
	"currency-converter/app/storage"
	"currency-converter/app/types"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.WithError(err).Fatal("Loading .env file")
	}

	if os.Getenv("DEBUG") == "1" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	bases := []types.CurrencyPair{
		{Base: "EUR", To: "USD"},
		{Base: "USD", To: "EUR"},
	}

	cacheDuration, err := strconv.Atoi(os.Getenv("APP_CACHE_DURATION_MIN"))
	if err != nil {
		logrus.WithError(err).Fatal("Parsing cache duration min")
	}

	con := converter.NewConverter(
		buildProvider(),
		bases,
		buildRateStorage(),
		time.Duration(cacheDuration)*time.Minute,
	)
	server := rest.Server{Converter: con, Stat: con, Port: os.Getenv("APP_REST_PORT")}

	go func() {
		logrus.Info("Starting rest server")
		logrus.WithError(server.Start()).Fatal("Rest server can't start")
	}()

	logrus.Info("Starting rpc server")
	rpc := jrpc.Server{Converter: con, Port: os.Getenv("APP_RPC_PORT")}
	logrus.WithError(rpc.Start()).Fatal("RPC server can't start")
}

func buildRateStorage() types.RateStorage {
	switch os.Getenv("STORAGE_TYPE") {
	case "memory":
		return storage.NewMemoryStorage()
	case "redis":
		redisAddr := fmt.Sprintf("%s:%s", os.Getenv("STORAGE_REDIS_HOST"), os.Getenv("STORAGE_REDIS_PORT"))
		return storage.NewRedisStorage(redis.NewClient(&redis.Options{Addr: redisAddr}))
	default:
		panic("Unknown store type. Available types are: 'memory', 'redis'")
	}
}

func buildProvider() types.RateProvider {
	switch os.Getenv("APP_RATE_PROVIDER") {
	case "fake":
		return provider.NewFakeProvider(decimal.New(10, 0), time.Now())
	case "exchangeratesapi.io":
		return provider.NewExchangeRatesIOProvider()
	default:
		panic("Unknown provider type. Available types are: 'fake', 'exchangeratesapi.io'")
	}
}
