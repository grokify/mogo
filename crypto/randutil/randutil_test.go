package randutil

import (
	"math"
	"strings"
	"testing"
	"time"

	"github.com/grokify/mogo/encoding/basex"
	"github.com/grokify/mogo/type/number"
)

// callTimeout is generous for any non-pathological range: rejection sampling
// rejects fewer than half of draws, so a hang means a retry-forever bug.
const callTimeout = 5 * time.Second

type port uint16 // a named type permitted by the number.Integer constraint

// draw calls CryptoRandIntInRange, failing the test instead of hanging if the
// call does not return.
func draw[T number.Integer](t *testing.T, min, max T) T {
	t.Helper()
	type result struct {
		v   T
		err error
	}
	ch := make(chan result, 1)
	go func() {
		v, err := CryptoRandIntInRange(min, max)
		ch <- result{v, err}
	}()
	select {
	case r := <-ch:
		if r.err != nil {
			t.Fatalf("CryptoRandIntInRange(%v, %v) error: %v", min, max, r.err)
		}
		return r.v
	case <-time.After(callTimeout):
		t.Fatalf("CryptoRandIntInRange(%v, %v) did not return within %v", min, max, callTimeout)
	}
	panic("unreachable")
}

// checkRange draws n values and fails if any falls outside [min, max]. For
// spans of at most 16 values it also requires every value to appear; with
// n = 2000 the chance of missing one by luck is below 1e-55.
func checkRange[T number.Integer](t *testing.T, min, max T, n int) {
	t.Helper()
	seen := map[T]bool{}
	for range n {
		v := draw(t, min, max)
		if v < min || v > max {
			t.Fatalf("CryptoRandIntInRange(%v, %v) = %v, out of range", min, max, v)
		}
		seen[v] = true
	}
	// Compare the distance, not span = distance+1, which overflows to 0 for
	// full 64-bit ranges.
	if diff := uint64(max) - uint64(min); diff < 16 && uint64(len(seen)) != diff+1 {
		t.Errorf("CryptoRandIntInRange(%v, %v): saw %d distinct values in %d draws, want %d", min, max, len(seen), n, diff+1)
	}
}

func TestCryptoRandIntInRange(t *testing.T) {
	t.Run("int small", func(t *testing.T) { checkRange(t, 0, 9, 2000) })
	t.Run("int negative", func(t *testing.T) { checkRange(t, -5, 5, 2000) })
	t.Run("int32 negative", func(t *testing.T) { checkRange[int32](t, -3, 3, 2000) })
	t.Run("uint8 top", func(t *testing.T) { checkRange[uint8](t, 250, 255, 2000) })
	t.Run("int8 wide", func(t *testing.T) { checkRange[int8](t, -100, 100, 2000) })
	t.Run("int8 full", func(t *testing.T) { checkRange[int8](t, math.MinInt8, math.MaxInt8, 2000) })
	t.Run("uint8 full", func(t *testing.T) { checkRange[uint8](t, 0, math.MaxUint8, 2000) })
	t.Run("int16 wide", func(t *testing.T) { checkRange[int16](t, -30000, 30000, 2000) })
	t.Run("uint16 full", func(t *testing.T) { checkRange[uint16](t, 0, math.MaxUint16, 2000) })
	t.Run("uint32 full", func(t *testing.T) { checkRange[uint32](t, 0, math.MaxUint32, 2000) })
	t.Run("int64 full", func(t *testing.T) { checkRange[int64](t, math.MinInt64, math.MaxInt64, 2000) })
	t.Run("uint64 full", func(t *testing.T) { checkRange[uint64](t, 0, math.MaxUint64, 2000) })
	t.Run("named type", func(t *testing.T) { checkRange[port](t, 8080, 8083, 2000) })
}

func TestCryptoRandIntInRangeSingleValue(t *testing.T) {
	if v := draw(t, 7, 7); v != 7 {
		t.Errorf("CryptoRandIntInRange(7, 7) = %d, want 7", v)
	}
	if v := draw[int8](t, math.MinInt8, math.MinInt8); v != math.MinInt8 {
		t.Errorf("CryptoRandIntInRange(MinInt8, MinInt8) = %d, want %d", v, math.MinInt8)
	}
	if v := draw[uint64](t, math.MaxUint64, math.MaxUint64); v != math.MaxUint64 {
		t.Errorf("CryptoRandIntInRange(MaxUint64, MaxUint64) = %d, want %d", v, uint64(math.MaxUint64))
	}
}

func TestCryptoRandIntInRangeMinGreaterThanMax(t *testing.T) {
	if _, err := CryptoRandIntInRange(5, 4); err == nil {
		t.Error("CryptoRandIntInRange(5, 4) error = nil, want error")
	}
	if _, err := CryptoRandIntInRange[int8](1, -1); err == nil {
		t.Error("CryptoRandIntInRange[int8](1, -1) error = nil, want error")
	}
}

// TestCryptoRandIntInRangeUniform guards against modulo bias: with 20000
// draws over 10 values each bucket expects 2000 (sd ~42), so a 15% tolerance
// is about 7 standard deviations.
func TestCryptoRandIntInRangeUniform(t *testing.T) {
	const n, buckets = 20000, 10
	counts := make([]int, buckets)
	for range n {
		counts[draw(t, 0, buckets-1)]++
	}
	want := n / buckets
	for i, c := range counts {
		if math.Abs(float64(c-want)) > 0.15*float64(want) {
			t.Errorf("value %d drawn %d times, want about %d", i, c, want)
		}
	}
}

func TestIntn(t *testing.T) {
	for range 1000 {
		if v := Intn(3); v < 0 || v >= 3 {
			t.Fatalf("Intn(3) = %d, want [0, 3)", v)
		}
	}
	defer func() {
		if recover() == nil {
			t.Error("Intn(0) did not panic")
		}
	}()
	Intn(0)
}

func TestInt63(t *testing.T) {
	for range 1000 {
		if v := Int63(); v < 0 {
			t.Fatalf("Int63() = %d, want non-negative", v)
		}
	}
}

func TestFloat64(t *testing.T) {
	for range 1000 {
		f, err := Float64()
		if err != nil {
			t.Fatal(err)
		}
		if f < 0 || f >= 1 {
			t.Fatalf("Float64() = %v, want [0, 1)", f)
		}
		if g := MustFloat64(); g < 0 || g >= 1 {
			t.Fatalf("MustFloat64() = %v, want [0, 1)", g)
		}
	}
}

func TestRandString(t *testing.T) {
	s, err := RandString("ab", 64)
	if err != nil {
		t.Fatal(err)
	}
	if len(s) != 64 || strings.Trim(s, "ab") != "" {
		t.Errorf(`RandString("ab", 64) = %q, want 64 chars of "ab"`, s)
	}

	s, err = RandString("", 32)
	if err != nil {
		t.Fatal(err)
	}
	if len(s) != 32 || strings.Trim(s, basex.AlphabetBase16) != "" {
		t.Errorf(`RandString("", 32) = %q, want 32 base16 chars`, s)
	}

	if s, err = RandString("ab", 0); err != nil || s != "" {
		t.Errorf(`RandString("ab", 0) = (%q, %v), want ("", nil)`, s, err)
	}
}

func TestNewSeedInt64Crypto(t *testing.T) {
	for range 100 {
		v, err := NewSeedInt64Crypto()
		if err != nil {
			t.Fatal(err)
		}
		if v < 0 {
			t.Fatalf("NewSeedInt64Crypto() = %d, want non-negative", v)
		}
	}
}

func TestNewMathRandCryptoSource(t *testing.T) {
	if NewMathRandCryptoSource() == nil {
		t.Fatal("NewMathRandCryptoSource() = nil")
	}
}
