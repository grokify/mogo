package bigint

import (
	"math/big"
	"testing"

	"github.com/grokify/mogo/encoding/basex"
)

var modtTests = []struct {
	a    int64
	n    int64
	want int64
}{
	{0, 2, 0},
	{1, 2, 1},
	{2, 2, 0},
	{3, 2, 1},
	{4, 2, 0},
	{12, 3, 0},
	{13, 3, 1},
}

func TestMod(t *testing.T) {
	for _, tt := range modtTests {
		got := ModInt64(tt.a, tt.n)
		if got != tt.want {
			t.Errorf("bigutil.ModIn64(%d,%d): want [%d], got [%d]", tt.a, tt.n, tt.want, got)
		}
		aBig := big.NewInt(tt.a)
		nBig := big.NewInt(tt.n)
		gotBig := Mod(aBig, nBig)
		if gotBig.Int64() != tt.want {
			t.Errorf("bigutil.Mod(%d,%d): want [%d], got [%d]", tt.a, tt.n, tt.want, gotBig.Int64())
		}
	}
}

var powBigIntTests = []struct {
	x    int64
	y    int64
	want string
}{
	{2, 0, "1"},
	{2, 1, "2"},
	{2, 2, "4"},
	{2, 3, "8"},
	{2, 4, "16"}}

func TestPowInt(t *testing.T) {
	for _, tt := range powBigIntTests {
		x := big.NewInt(tt.x)
		y := big.NewInt(tt.y)
		pow := Pow(x, y)

		if pow.String() != tt.want {
			t.Errorf("bigutil.PowInt(%v,%v): want [%v], got [%v]", tt.x, tt.y, tt.want, pow.String())
		}
	}
}

var intToBaseXStringTests = []struct {
	base int
	dec  int
	val  string
}{
	{16, 15, "f"},
	{32, 0, "0"},
	{32, 1, "1"},
	{32, 12, "c"}}

const base36Dictionary = "0123456789abcdefghijklmnopqrstuvwxyz"

func TestIntToBaseXString(t *testing.T) {
	for _, tt := range intToBaseXStringTests {
		try := Int64ToBaseX(int64(tt.dec), tt.base)

		if try != tt.val {
			t.Errorf("bigint.IntToBaseXString(%v,%v): want [%v], got [%v]",
				tt.dec, tt.base, tt.val, try)
		}
		try2 := string([]rune(base36Dictionary)[tt.dec])
		if try2 != tt.val {
			t.Errorf("bigint.IntToBaseXString(%v,%v): want [%s], got [%s]",
				tt.dec, tt.base, tt.val, try2)
		}
	}
}

var intToBaseXAlphabetTests = []struct {
	alphabet string
	dec      int64
	val      string
}{
	{basex.AlphabetBase16, 15, "f"},
	{basex.AlphabetBase16, 16, "10"},
	{basex.AlphabetBase16, 32, "20"},
	{basex.AlphabetBase32, 0, "A"},
	{basex.AlphabetBase32, 2, string(basex.AlphabetBase32[2])},
	{basex.AlphabetBase32, 31, "7"},
	{basex.AlphabetBase32, 32, "BA"},
	{basex.AlphabetBase32, 33, "BB"},
	{basex.AlphabetBase58, 58, "21"},
}

func TestIntToBaseXAlphabet(t *testing.T) {
	for _, tt := range intToBaseXAlphabetTests {
		enc, err := Int64ToBaseXAlphabet(tt.dec, tt.alphabet)
		if err != nil {
			t.Errorf("bigint.Int64ToBaseXAlphabet(%v,%v): error [%v]",
				tt.dec, tt.alphabet, err.Error())
		}
		if enc != tt.val {
			t.Errorf("bigint.Int64ToBaseXAlphabet(%v,%v): want [%v], got [%v]",
				tt.dec, tt.alphabet, tt.val, enc)
		}
		dec, err := BaseXAlphabetToInt64(enc, tt.alphabet)
		if err != nil {
			t.Errorf("bigint.BaseXAlphabetToInt64(%v,%v): error [%v]",
				enc, tt.alphabet, err.Error())
		}
		if dec != tt.dec {
			t.Errorf("bigint.BaseXAlphabetToInt64(%v,%v): want [%v], got [%v]",
				tt.dec, tt.alphabet, tt.val, dec)
		}
	}
}

var splitTests = []struct {
	x     int64
	scale uint
	y     int64
	z     int64
}{
	{12345, 2, 123, 45},
	{8888999, 3, 8888, 999},
}

func TestSplitInt64(t *testing.T) {
	for _, tt := range splitTests {
		y, z := SplitInt64(tt.x, tt.scale)
		if y != tt.y || z != tt.z {
			t.Errorf("bigutil.SplitInt64(%v,%v): want [%d,%d], got [%d,%d]", tt.x, tt.scale, tt.y, tt.z, y, z)
		}
	}
}

var parseScientificNotationTests = []struct {
	v    string
	want string
	acc  big.Accuracy
}{
	{"1.84467440737095E+19", "18446744073709500000", big.Exact},
}

func TestParseScientificNotation(t *testing.T) {
	for _, tt := range parseScientificNotationTests {
		i, acc, err := ParseScientificNotation(tt.v)
		if err != nil {
			t.Errorf("biting.ParseScientificNotation(\"%s\"): error (%s)", tt.v, err.Error())
		}
		got := i.String()
		if got != tt.want {
			t.Errorf("biting.ParseScientificNotation(\"%s\"): mismatch (int): want (%s) got (%s)", tt.v, tt.want, got)
		}
		if acc != tt.acc {
			t.Errorf("biting.ParseScientificNotation(\"%s\"): mismatch (accuracy): want (%s=d) got (%d)", tt.v, tt.acc, acc)
		}
	}
}
