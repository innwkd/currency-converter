package jrpc

import (
	"net/http"

	"github.com/shopspring/decimal"
	"github.com/yddmat/currency-converter/app/types"
)

type ConverterArgs struct {
	Pair   types.CurrencyPair `json:"pair"`
	Amount decimal.Decimal    `json:"amount"`
}

type Converter struct {
	Converter types.Converter
}

func (c *Converter) Convert(r *http.Request, args *ConverterArgs, reply *types.Conversion) (err error) {
	*reply, err = c.Converter.Convert(args.Pair, args.Amount)
	return err
}
