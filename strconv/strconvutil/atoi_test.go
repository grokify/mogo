package strconvutil

import (
	"errors"
	"math"
	"testing"
)

func TestAtou32(t *testing.T) {
	tests := []struct {
		s       string
		want    uint32
		wantErr error
	}{
		{"0", 0, nil},
		{"42", 42, nil},
		{"4294967295", math.MaxUint32, nil},
		{"4294967296", 0, ErrValueIsOutOfRange},
		{"-1", 0, ErrValueIsNegative},
	}
	for _, tt := range tests {
		got, err := Atou32(tt.s)
		if !errors.Is(err, tt.wantErr) || got != tt.want {
			t.Errorf("Atou32(%q) = (%d, %v), want (%d, %v)", tt.s, got, err, tt.want, tt.wantErr)
		}
	}
	if _, err := Atou32("abc"); err == nil {
		t.Error(`Atou32("abc") error = nil, want error`)
	}
}
