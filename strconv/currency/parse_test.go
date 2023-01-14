package currency

import (
	"testing"

	"github.com/shopspring/decimal"
)

var parseTests = []struct {
	v    string
	want Amount
}{
	{"USD $1,500.00", Amount{
		Unit:  CurrencyUSD,
		Value: decimal.NewFromInt(1500)}},
	{"GBP 1,500.00", Amount{
		Unit:  CurrencyGBP,
		Value: decimal.NewFromInt(1500)}},
	{"C $1,500.00", Amount{
		Unit:  CurrencyCAD,
		Value: decimal.NewFromInt(1500)}},
}

// TestParse tests parsing curency.
func TestParse(t *testing.T) {
	for _, tt := range parseTests {
		try, err := ParseAmount(tt.v)
		if err != nil {
			t.Errorf("currency.Parse(\"%v\") Error: [%v]", tt.v, err.Error())
		}
		if !try.Value.Equal(tt.want.Value) {
			t.Errorf("timeutil.SliceMinMax(\"%v\") Mismatch: on Value want [%v], got [%v]", tt.v,
				tt.want.Value.String(), try.Value.String())
		}
		if try.Unit != tt.want.Unit {
			t.Errorf("timeutil.SliceMinMax(\"%v\") Mismatch on Unit: want [%v], got [%v]", tt.v,
				tt.want.Unit, try.Unit)
		}
	}
}
