package currencyutil

import (
	"errors"
	"strings"

	"github.com/grokify/mogo/errors/errorsutil"
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

func MustSymbol(iso4217 string) string {
	s, err := Symbol(iso4217)
	if err != nil {
		panic(err)
	}
	return s
}

func MustSymbolOrDefault(iso4217, defSymbol string) string {
	iso4217 = strings.TrimSpace(iso4217)
	if iso4217 == "" {
		return defSymbol
	}
	s, err := Symbol(iso4217)
	if err != nil {
		panic(err)
	}
	return s
}

var ErrSymbolNotFound = errors.New("symbol not found")

func Symbol(iso4217 string) (string, error) {
	iso4217 = strings.ToUpper(strings.TrimSpace(iso4217))
	if sym, ok := SymbolMap[iso4217]; ok {
		return sym, nil
	}
	return iso4217, errorsutil.Wrapf(ErrSymbolNotFound, "cannot find symbol (%s)", iso4217)
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
