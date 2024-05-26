package mimeutil

import (
	"testing"

	"github.com/grokify/mogo/net/http/httputilmore"
)

var isTypeTests = []struct {
	v       string
	tryType string
	isType  bool
}{
	{"IMAGE/PNG", httputilmore.ContentTypeImagePNG, true},
	{"IMAGE/PNG;", httputilmore.ContentTypeImagePNG, true},
	{"IMAGE/PNG ;", httputilmore.ContentTypeImagePNG, true},
	{"XIMAGE/PNG ;", httputilmore.ContentTypeImagePNG, false},
	{httputilmore.ContentTypeAppXMLUtf8, httputilmore.ContentTypeAppXML, true},
}

func TestIsType(t *testing.T) {
	for _, tt := range isTypeTests {
		isType := IsType(tt.tryType, tt.v)
		if isType != tt.isType {
			t.Errorf("mimeutil.IsType(\"%s\", \"%s\") Fail: want [%v] got [%v]",
				tt.v, tt.tryType, tt.isType, isType)
		}
	}
}
