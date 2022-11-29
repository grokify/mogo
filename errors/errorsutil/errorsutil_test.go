package errorsutil

import (
	"errors"
	"testing"
)

var wrapTests = []struct {
	err1String     string
	err2String     string
	err2wrapString string
}{
	{"foo", "bar", "bar: [foo]"},
	{"foobar", "bazqux", "bazqux: [foobar]"},
}

func TestWrap(t *testing.T) {
	for _, tt := range wrapTests {
		err1 := errors.New(tt.err1String)
		err2 := Wrap(err1, tt.err2String)
		if err2.Error() != tt.err2wrapString {
			t.Errorf("errorsutil.Wrap: want [%v] got [%v]", tt.err2wrapString, err2.Error())
		}
		unwrapped := errors.Unwrap(err2)
		if unwrapped.Error() != tt.err1String {
			t.Errorf("errorsutil.Wrap/Unwrap: want [%v] got [%v]", tt.err1String, unwrapped.Error())
		}
	}
}
