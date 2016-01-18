package timeutil

import (
	"testing"
)

var dt6ForDt14Tests = []struct {
	v    int64
	want int32
}{
	{int64(20060102150405), int32(20060102)}}

func TestDt6ForDt14(t *testing.T) {
	for _, tt := range dt6ForDt14Tests {
		got := Dt6ForDt14(tt.v)
		if got != tt.want {
			t.Errorf("Dt6ForDt14(%v): want %v, got %v", tt.v, tt.want, got)
		}
	}
}
