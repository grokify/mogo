package currencyutil

/*
import (
	"strconv"
	"strings"

	"github.com/leekchan/accounting"
)

func FormatMoney(symbol string, amount float64, precision int) string {
	ac := accounting.Accounting{Symbol: SymbolPrefix(symbol), Precision: precision}
	return ac.FormatMoney(amount)
}

type Formatter struct {
	Accounting           accounting.Accounting
	SymbolPrefix         string
	DefaultISOCode       string
	PrefixCodeDefault    bool
	PrefixCodeNonDefault bool
}

func NewFormatter(symbol string) Formatter {
	if len(symbol) == 0 {
		symbol = "$"
	}
	return Formatter{
		Accounting:   accounting.Accounting{Symbol: symbol, Precision: 2},
		SymbolPrefix: symbol}
}

func (f *Formatter) FormatMoney(symbol string, amount float64) string {
	symbolPrefix := SymbolPrefix(symbol)

	if symbolPrefix != f.SymbolPrefix {
		f.Accounting = accounting.Accounting{Symbol: symbolPrefix, Precision: 2}
	}
	return f.Accounting.FormatMoney(amount)
}

func (f *Formatter) FormatMoneyString(symbol string, amountString string) (string, error) {
	symbolPrefix := SymbolPrefix(symbol)
	amountFloat, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		return strings.Join([]string{symbolPrefix, amountString}, ""), err
	}
	if symbolPrefix != f.SymbolPrefix {
		f.Accounting = accounting.Accounting{Symbol: symbolPrefix, Precision: 2}
	}
	return f.Accounting.FormatMoney(amountFloat), nil
}
*/
