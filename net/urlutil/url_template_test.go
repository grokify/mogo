package urlutil

import (
	"testing"
)

var parseURLTemplateTests = []struct {
	v    string
	want string
}{
	{"https://{customer}.example.com:{port}/v5", "/v5"},
	{"https://customer.example.com:0/{v5}", "/{v5}"},
	{"https://%7BcustomerId%7D.sexample.com:0/v8.0", "/v8.0"},
	{"https://{customer}.example.com:0/v8.0/", "/v8.0/"},
}

func TestParseURLTemplate(t *testing.T) {
	for _, tt := range parseURLTemplateTests {
		u1, err := ParseURLTemplate(tt.v)
		if err != nil {
			t.Errorf("urlutil.ParseURLTemplate() Error: input [%v], want [%v], got error [%v]",
				tt.v, tt.want, err.Error())
		}
		if u1.Path != tt.want {
			t.Errorf("urlutil.ParseURLTemplate() Failure: input [%v], want [%v], got [%v]",
				tt.v, tt.want, u1.Path)
		}
	}
}
