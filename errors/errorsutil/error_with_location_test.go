package errorsutil

import (
	"errors"
	"strings"
	"testing"
)

var errorWithLocationTests = []struct {
	fn           func() error
	errMsgSuffix string
}{
	{func() error { return errors.New("test err") }, `mogo/errors/errorsutil/error_with_location_test.go:22)`},
}

func TestErrorWithLocation(t *testing.T) {
	for _, tt := range errorWithLocationTests {
		tryErr := tt.fn()
		if tryErr == nil {
			panic("no error")
		}
		tryWithLocation := NewErrorWithLocation(tryErr.Error())
		if !strings.HasSuffix(tryWithLocation.Error(), tt.errMsgSuffix) {
			t.Errorf("errorsutil.NewErrorWithLocation(\"%s\"): mismatch want suffix [%s] got [%s]", tryErr.Error(), tt.errMsgSuffix, tryWithLocation.Error())
		}
	}
}
