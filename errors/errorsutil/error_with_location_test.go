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
	{func() error { _, err := time.Parse(time.RFC3339, "Mon, 2024-12-31"); return err }, `github.com/grokify/mogo/errors/errorsutil/error_with_location_test.go:22)`},
}

func TestErrorWithLocation(t *testing.T) {
	for _, tt := range errorWithLocationTests {
		try := tt.fn()
		if try == nil {
			panic("no error")
		}
		tryWithLocation := NewErrorWithLocation(try.Error())
		if !strings.HasSuffix(tryWithLocation.Error(), tt.errMsgSuffix) {
			t.Errorf("errorsutil.NewErrorWithLocation(\"%s\"): mismatch want suffix [%s] got [%s]", try.Error(), tt.errMsgSuffix, tryWithLocation.Error())
		}
	}
}
