package fmtutil

import (
	"io"
	"os"
	"testing"

	"github.com/grokify/mogo/pointer"
)

var fprintxTests = []struct {
	w     io.Writer
	str   *string
	args  []any
	wantN int
}{
	{os.Stdout, pointer.Pointer("testprint1"), []any{}, 10},
	{nil, pointer.Pointer("testprint1"), []any{}, 0},
	{os.Stdout, nil, []any{"testprint2"}, 10},
	{nil, nil, []any{"testprint2"}, 0},
}

func TestFprintxIf(t *testing.T) {
	for _, tt := range fprintxTests {
		if tt.str != nil {
			n, err := FprintfIf(tt.w, *tt.str, tt.args...)
			if err != nil {
				t.Errorf("fmtutil.FprintfIf(%v, %v, %v) Error: (%s)",
					tt.w, tt.str, tt.args, err.Error())
			} else if n != tt.wantN {
				t.Errorf("fmtutil.FprintfIf(%v, %v, %v) Mismatch: want (%d) got (%d)",
					tt.w, tt.str, tt.args, tt.wantN, n)
			}
		} else {
			n, err := FprintIf(tt.w, tt.args...)
			if err != nil {
				t.Errorf("fmtutil.FprintfIf(%v, %v) Error: (%s)",
					tt.w, tt.args, err.Error())
			} else if n != tt.wantN {
				t.Errorf("fmtutil.FprintfIf(%v, %v) Mismatch: want (%d) got (%d)",
					tt.w, tt.args, tt.wantN, n)
			}
		}
	}
}
