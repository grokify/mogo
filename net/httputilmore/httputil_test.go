package httputilmore

import (
	"testing"
)

var parseHeaderTests = []struct {
	v    string
	want [][]string
}{
	{"Foo: Bar", [][]string{{"Foo", "Bar"}}},
	{"Foo: Bar\nBaz: QUX", [][]string{{"Foo", "Bar"}, {"Baz", "QUX"}}},
	{"Foo: Bar\nBaz: QUX\nQUUX: quuz", [][]string{{"Foo", "Bar"}, {"Baz", "QUX"}, {"quux", "quuz"}}},
}

func TestParseHeader(t *testing.T) {
	for _, tt := range parseHeaderTests {
		header := ParseHeader(tt.v)

		for _, pair := range tt.want {
			got := header.Get(pair[0])
			want := pair[1]
			if got != want {
				t.Errorf("httputilmore.ParseHeader() Error: header [%v], want [%v], got [%v]",
					pair[0], want, got)
			}
		}
	}
}
