package currency

import (
	"fmt"
	"strings"
)

type ExchangeRateSimple struct {
	BaseCurrency  string
	CurrencyRates map[string]float64
}

func (xr *ExchangeRateSimple) ConvertToBase(in float64, cur string) (float64, error) {
	cur = strings.ToUpper(strings.TrimSpace(cur))
	if multiplier, ok := xr.CurrencyRates[cur]; ok {
		return in * multiplier, nil
	}
	return in, fmt.Errorf("E_EXCHANGE_RATE_CANNOT_FIND [%s]", cur)
}

func ExampleExchangeRates() ExchangeRateSimple {
	return ExchangeRateSimple{
		BaseCurrency: "USD",
		CurrencyRates: map[string]float64{
			"AUD": 0.637368,
			"CAD": 0.709545,
			"EUR": 1.08641,
			"GBP": 1.24626,
			"USD": 1.0}}
}
