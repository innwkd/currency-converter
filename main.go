package main

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/yddmat/currency-converter/app/api/http"
	"github.com/yddmat/currency-converter/app/converter"
	"github.com/yddmat/currency-converter/app/provider"
)

func main() {
	fakeRate, _ := decimal.NewFromString("26.6")

	server := http.Server{
		Converter: converter.NewConverter(&provider.FakeProvider{Rate: fakeRate, UpdateTime: time.Unix(5, 0)}),
	}

	server.Start()
}
