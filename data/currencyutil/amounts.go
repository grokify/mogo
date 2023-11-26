package currencyutil

/*
import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grokify/mogo/strconv/strconvutil"
	"github.com/shopspring/decimal"
)

type Amounts []Amount

func (amts Amounts) Sum() (Amount, error) {
	sum := Amount{}
	for _, amt := range amts {
		err := sum.Add(amt)
		if err != nil {
			return sum, err
		}
	}
	return sum, nil
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
*/
