package currency

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

const (
	USDAbbr           = "USD"
	USDSymbol         = "$"
	UnitsMillionsDesc = "millions"
	UnitsMillionsInt  = 1000000
)

var (
	ErrUnknownCurrencyPrefix = errors.New("unknown currency prefix")
	ErrUnknownCurrencySuffix = errors.New("unknown currency suffix")
	rxCurrencyPrefix         = regexp.MustCompile(`^(\D+)(\d.*)$`)
	rxCurrencySuffix         = regexp.MustCompile(`^(.*)([^\d.,].*)$`)
)

type ParseCurrencyOpts struct {
	Comma         string
	Decimal       string
	BillionsAbbr  []string
	MillionsAbbr  []string
	ThousandsAbbr []string
}

func NewParseCurrencyOpts() ParseCurrencyOpts {
	return ParseCurrencyOpts{
		BillionsAbbr:  []string{},
		MillionsAbbr:  []string{},
		ThousandsAbbr: []string{}}
}

func ParseCurrency(opts *ParseCurrencyOpts, s string) (string, float64, error) {
	try := strings.TrimSpace(s)
	if len(try) == 0 {
		return "", 0, nil
	}
	if opts == nil {
		newopts := NewParseCurrencyOpts()
		opts = &newopts
	}
	var currency string
	var numeric string
	var units string
	mPrefix := rxCurrencyPrefix.FindStringSubmatch(try)
	if len(mPrefix) > 0 {
		prefix := mPrefix[1]
		numeric = mPrefix[2]
		if len(prefix) == 3 {
			currency = strings.ToUpper(prefix)
		} else if len(prefix) == 1 {
			if prefix == "$" {
				currency = USDAbbr
			}
		} else {
			return "", 0, ErrUnknownCurrencyPrefix
		}
	}
	mSuffix := rxCurrencySuffix.FindStringSubmatch(numeric)
	if len(mSuffix) > 0 {
		numeric = mSuffix[1]
		suffix := mSuffix[2]
		for _, mtry := range opts.MillionsAbbr {
			if suffix == mtry {
				units = UnitsMillionsDesc
			}
		}
		if len(strings.TrimSpace(units)) == 0 {
			return "", 0, ErrUnknownCurrencySuffix
		}
	}

	val, err := strconv.ParseFloat(numeric, 64)
	if err != nil {
		return currency, 0, err
	}

	return currency, val, nil
}
