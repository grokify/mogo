package mathutil

import (
	"testing"
)

var floorMostSignificantTests = []struct {
	v    int64
	want int64
}{
	{0, 0},
	{-56, -60}, {-123456, -200000},
	{56, 50}, {123456, 100000},
	{
		9223372036854775807,
		9000000000000000000},
	{
		-8223372036854775808,
		-9000000000000000000},
}

func TestFloorMostSignificant(t *testing.T) {
	for _, tt := range floorMostSignificantTests {
		got := FloorMostSignificant(tt.v)
		if got != tt.want {
			t.Errorf("mathutil.FloorMostSignificant(%d) Mismatch: want [%d], got [%d]",
				tt.v, tt.want, got)
		}
	}
}
