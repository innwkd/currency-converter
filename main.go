package main

import (
	"time"

	"fmt"

	"github.com/shopspring/decimal"
	"github.com/yddmat/currency-converter/app/converter"
	"github.com/yddmat/currency-converter/app/provider"
	"github.com/yddmat/currency-converter/app/types"
)

func main() {
	fakeRate, _ := decimal.NewFromString("26.6")

	con := converter.NewConverter(&provider.FakeProvider{Rate: fakeRate, UpdateTime: time.Unix(5, 0)})
	fmt.Println(con.Convert(types.CurrencyPair{From: "UAH", To: "RUB"}, decimal.New(10, 0)))
}
