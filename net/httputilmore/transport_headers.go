package httputilmore

import (
	"net/http"
)

// TransportRequestModifier implements http.RoundTripper.
// When set as Transport of http.Client, it adds HTTP headers and or query
// string parameters to requests.
// No field is mandatory. Can be implemented with http.Client as:
// client.Transport = httputilmore.TransportRequestModifier{
// Transport:client.Transport, Header:myHeader}
type TransportRequestModifier struct {
	Transport http.RoundTripper
	Header    http.Header
	Query     map[string][]string // use raw struct to avoid param capitalization with `http.Header`
	Override  bool
}

// RoundTrip adds the additional headers per request implements http.RoundTripper.
func (t TransportRequestModifier) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header = MergeHeader(req.Header, t.Header, t.Override)

	if len(t.Query) > 0 {
		q := req.URL.Query()
		for k, vals := range t.Query {
			for _, v := range vals {
				if t.Override || len(q.Get(k)) == 0 {
					q.Add(k, v)
				}
			}
		}
		req.URL.RawQuery = q.Encode()
	}

	return t.transport().RoundTrip(req)
}

func (t TransportRequestModifier) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}

	return http.DefaultTransport
}
