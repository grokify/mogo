package currencyutil

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/grokify/mogo/strconv/strconvutil"
	"github.com/shopspring/decimal"
)

const (
	USDSymbol          = "$"
	UnitsBillionsDesc  = "billions"
	UnitsBillionsInt   = 1000000000
	UnitsMillionsDesc  = "millions"
	UnitsMillionsInt   = 1000000
	UnitsThousandsDesc = "thousands"
	UnitsThousandsInt  = 1000
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
				currency = CurrencyUSD
			}
		} else {
			return "", 0, ErrUnknownCurrencyPrefix
		}
	}
	mSuffix := rxCurrencySuffix.FindStringSubmatch(numeric)
	if len(mSuffix) > 0 {
		numeric = mSuffix[1]
		suffix := mSuffix[2]
		for _, btry := range opts.BillionsAbbr {
			if suffix == btry {
				units = UnitsBillionsDesc
				break
			}
		}
		if len(units) == 0 {
			for _, mtry := range opts.MillionsAbbr {
				if suffix == mtry {
					units = UnitsMillionsDesc
					break
				}
			}
		}
		if len(units) == 0 {
			for _, ttry := range opts.ThousandsAbbr {
				if suffix == ttry {
					units = UnitsThousandsDesc
					break
				}
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

	switch strings.ToLower(strings.TrimSpace(units)) {
	case UnitsBillionsDesc:
		val *= float64(UnitsBillionsInt)
	case UnitsMillionsDesc:
		val *= float64(UnitsMillionsInt)
	case UnitsThousandsDesc:
		val *= float64(UnitsThousandsInt)
	}

	return currency, val, nil
}

type Amount struct {
	Value decimal.Decimal
	Unit  string
}

var (
	// rxUSD   = regexp.MustCompile(`^USD?\s+\$([0-9,]+\.[0-9]{2})$`)
	rxCUR   = regexp.MustCompile(`^([A-Z]{1,3})\s+([^0-9])?([0-9,]+\.[0-9]{2})$`)
	rxComma = regexp.MustCompile(`,`)
)

func NewAmountFloat(u string, v float64) (Amount, error) {
	amt := Amount{
		Value: decimal.NewFromFloat(v)}
	can, err := ParseCurrencyUnit(u, "")
	if err != nil {
		return amt, err
	}
	amt.Unit = can
	return amt, nil
}

func MustNewAmountInt(u string, v int64, exp int32) Amount {
	amt, err := NewAmountInt(u, v, exp)
	if err != nil {
		panic(err)
	}
	return amt
}

func NewAmountInt(u string, v int64, exp int32) (Amount, error) {
	amt := Amount{
		Value: decimal.New(v, exp)}
	can, err := ParseCurrencyUnit(u, "")
	if err != nil {
		return amt, err
	}
	amt.Unit = can
	return amt, nil
}

func ParseAmount(value string) (Amount, error) {
	value = strings.TrimSpace(value)
	m := rxCUR.FindStringSubmatch(value)
	pr := Amount{}
	if len(m) > 0 {
		digits := m[3]
		digits = rxComma.ReplaceAllString(digits, "")
		n, err := decimal.NewFromString(digits)
		if err != nil {
			return pr, err
		}
		pr.Value = n
		cur, err := ParseCurrencyUnit(m[1], m[2])
		if err != nil {
			return pr, err
		}
		pr.Unit = cur
		return pr, nil
	}
	return pr, fmt.Errorf("currency not found [%s]", value)
}

/*
func (amt *Amount) ValueInt() (int64, error) {

}
*/

func (amt *Amount) Add(a Amount) error {
	if a.Value.IsZero() {
		return nil
	}
	amt.Unit = strings.ToUpper(strings.TrimSpace(amt.Unit))
	a.Unit = strings.ToUpper(strings.TrimSpace(a.Unit))
	if amt.Unit != a.Unit {
		if amt.Value.IsZero() && amt.Unit == "" {
			amt.Unit = a.Unit
		} else {
			return fmt.Errorf("mismatch currency units have (%s %v) adding (%s %v)", amt.Unit, amt.Value, a.Unit, a.Value)
		}
	}
	amt.Value = amt.Value.Add(a.Value)
	return nil
}

func (amt *Amount) Equal(a Amount) bool {
	return amt.Value.Equal(a.Value) && amt.Unit == a.Unit
}

var rxIntPrefix = regexp.MustCompile(`^[0-9]+`)

func (amt Amount) MustStringFixed(places int32, defaultUnit string) string {
	i := amt.Value.IntPart()
	is := strconvutil.Commify(i)
	if places < 1 {
		return MustSymbol(amt.Unit) + is
	}
	exp := amt.Value.StringFixed(places)
	exp = rxIntPrefix.ReplaceAllString(exp, "")
	unit := strings.TrimSpace(amt.Unit)
	if unit == "" {
		unit = defaultUnit
	}
	return MustSymbol(unit) + is + exp
}

func ParseCurrencyUnit(abbr, symbol string) (string, error) {
	abbr = CurrencyCanonical(abbr)
	if CurrencyCodeKnown(abbr) {
		return abbr, nil
	}
	return "", fmt.Errorf("could not parse currency [%s]", abbr)
}

func CurrencyCanonical(value string) string {
	value = strings.TrimSpace(strings.ToUpper(value))
	abbrs := map[string]string{
		"AU": CurrencyAUD,
		"C":  CurrencyCAD,
		"US": CurrencyUSD,
	}
	can, ok := abbrs[value]
	if ok {
		return can
	}
	return value
}
