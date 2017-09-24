package currencyutil

import (
	"errors"
	"strconv"
	"strings"

	"github.com/leekchan/accounting"
)

var SymbolMap = map[string]string{
	"AFN": "؋",
	"ARS": "$",
	"AWG": "ƒ",
	"BRL": "R$",
	"CAD": "$",
	"CLP": "$",
	"CNY": "¥",
	"CRC": "₡",
	"CUP": "₱",
	"EGP": "£",
	"EUR": "€",
	"FKP": "£",
	"GBP": "£",
	"ILS": "₪",
	"IRR": "﷼",
	"JPY": "¥",
	"KHR": "៛",
	"KLR": "₨",
	"KPW": "₩",
	"KRW": "₩",
	"MNT": "₮",
	"NGN": "₦",
	"NOK": "kr",
	"PHP": "₱",
	"PLN": "zł",
	"RUB": "₽",
	"SAR": "﷼",
	"THB": "฿",
	"UAH": "₴",
	"USD": "$",
	"UZS": "лв",
	"VND": "₫",
	"YER": "﷼"}

func Symbol(iso4217 string) (string, error) {
	iso4217 = strings.ToUpper(strings.TrimSpace(iso4217))
	if sym, ok := SymbolMap[iso4217]; ok {
		return sym, nil
	}
	return iso4217, errors.New("Cannot Find Symbol")
}

func SymbolPrefix(iso4217 string) string {
	iso4217 = strings.ToUpper(strings.TrimSpace(iso4217))
	sym, err := Symbol(iso4217)
	if err == nil {
		return sym
	}
	if len(iso4217) <= 2 {
		return iso4217
	}
	return strings.Join([]string{strings.ToUpper(iso4217), " "}, "")
}

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
