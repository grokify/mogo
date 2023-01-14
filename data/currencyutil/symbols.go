package currencyutil

import (
	"errors"
	"strings"
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
	return iso4217, errors.New("cannot find symbol")
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
