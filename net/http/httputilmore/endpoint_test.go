package httputilmore

import (
	"testing"
)

var parseEndpointTests = []struct {
	v          string
	wantMethod HTTPMethod
	wantURL    string
}{
	{"GET https://example.com", MethodGet, "https://example.com"},
	{"  PaTcH  https://example.com/ ", MethodPatch, "https://example.com/"},
}

func TestParseEndpoint(t *testing.T) {
	for _, tt := range parseEndpointTests {
		ep, err := ParseEndpoint(tt.v)
		if err != nil {
			t.Errorf("httputilmore.ParseEndpoint(\"%s\") Error: [%s]",
				tt.v, err.Error())
		}
		if ep.Method != tt.wantMethod {
			t.Errorf("httputilmore.ParseEndpoint(\"%s\") Fail Method: want [%s] got [%s]",
				tt.v, tt.wantMethod, ep.Method)
		}
		epURL := ep.URL.String()
		if epURL != tt.wantURL {
			t.Errorf("httputilmore.ParseEndpoint(\"%s\") Fail URL: want [%s] got [%s]",
				tt.v, tt.wantURL, epURL)
		}
	}
}
