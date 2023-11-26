package internal

import (
	_ "embed"
)

//go:embed currencies.tsv
var CurrenciesDataRaw string
