package number

import (
	"math"
	"testing"
)

var uint32Tests = []struct {
	v       int
	want    uint32
	wantErr bool
}{
	{v: -1, want: 0, wantErr: true},
	{v: 100, want: 100, wantErr: false},
	{v: int(^uint32(0)), want: math.MaxUint32, wantErr: false},
	{v: math.MaxInt, want: math.MaxUint32, wantErr: true},
	{v: math.MaxInt64, want: math.MaxUint32, wantErr: true},
	{v: 4294967295 + 1, want: math.MaxUint32, wantErr: true},
}

func TestUint32(t *testing.T) {
	for _, tt := range uint32Tests {
		got, err := Uint32(tt.v)
		if err != nil {
			if !tt.wantErr {
				t.Errorf("number.Uint32(%d) error (%s)",
					tt.v, err.Error())
			} else {
				continue
			}
		} else if tt.want != got {
			t.Errorf("number.Uint32(%d) want (%d) got (%d)",
				tt.v, tt.want, got)
		}
	}
}
