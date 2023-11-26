package currencyutil

import (
	"golang.org/x/exp/slices"
)

// activate when rebuilding this func.
// func CurrencyCodesAll() []string { return []string{} }

// CurrencyCodes0D returns zero decimal currencies. Listed at: https://stripe.com/docs/currencies
func CurrencyCodes0D() []string {
	return []string{
		CurrencyBIF,
		CurrencyCLP,
		CurrencyDJF,
		CurrencyGNF,
		CurrencyJPY,
		CurrencyKMF,
		CurrencyKRW,
		CurrencyMGA,
		CurrencyPYG,
		CurrencyRWF,
		CurrencyUGX,
		CurrencyVND,
		CurrencyVUV,
		CurrencyXAF,
		CurrencyXOF,
		CurrencyXPF}
}

// CurrencyCodes3D returns three decimal currencies. Listed at: https://stripe.com/docs/currencies
func CurrencyCodes3D() []string {
	return []string{
		CurrencyBHD,
		CurrencyJOD,
		CurrencyKWD,
		CurrencyOMR,
		CurrencyTND}
}

func CurrencyCodeKnown(value string) bool {
	return slices.Index(CurrencyCodesAll(), value) >= 0
}
