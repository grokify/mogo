package errorsutil

import (
	"strings"
	"testing"
	"time"
)

var errorWithLocationTests = []struct {
	fn           func() error
	errMsgSuffix string
}{
	{func() error { _, err := time.Parse(time.RFC3339, "Mon, 2024-12-31"); return err }, "github.com/grokify/mogo/errors/errorsutil/error_with_location_test.go:22"},
}

func TestErrorWithLocation(t *testing.T) {
	for _, tt := range errorWithLocationTests {
		try := tt.fn()
		if try == nil {
			panic("no error")
		}
		tryWithLocation := NewErrorWithLocation(try.Error())
		idx, idxCalc, ok := isSuffixOnly(tryWithLocation.Error(), tt.errMsgSuffix)
		if !ok {
			t.Errorf("errorsutil.NewErrorWithLocation(\"%s\"): mismatch want suffix (%s) got (%s), idx (%d) idxCalc (%d)", tryWithLocation.Error(), tt.errMsgSuffix, tryWithLocation.Error(), idx, idxCalc)
		}
	}
}

func isSuffixOnly(s, substr string) (int, int, bool) {
	idx := strings.Index(s, substr)
	idxCalc := len(s) - len(substr) - 1
	return idx, idxCalc, idx == idxCalc
}
