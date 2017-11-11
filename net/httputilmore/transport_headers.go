package httputilmore

import (
	"net/http"
)

// TransportWithHeaders implements http.RoundTripper. When set as Transport of http.Client, it adds HTTP headers to requests. No field is mandatory. Can be implemented with http.Client as: client.Transport = httputilmore.TransportWithHeaders{Transport:client.Transport, Header:myHeader}
type TransportWithHeaders struct {
	Transport http.RoundTripper
	Header    http.Header
	Override  bool
}

// THe default logging transport that wraps http.DefaultTransport.
var DefaultTransportWithHeaders = &TransportWithHeaders{
	Header:    http.Header{},
	Transport: http.DefaultTransport,
}

// RoundTrip is the core part of this module and implements http.RoundTripper.
// Executes HTTP request with request/response logging.
func (t TransportWithHeaders) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header = MergeHeader(req.Header, t.Header, t.Override)

	return t.transport().RoundTrip(req)
}

func (t TransportWithHeaders) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}

	return http.DefaultTransport
}
