package tlsutil

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"

	"github.com/grokify/mogo/errors/errorsutil"
	"golang.org/x/net/context/ctxhttp"
)

// SupportsTLSVersion returns an error if a connection cannot be made and a nil
// if the connection is successful.
func SupportsTLSVersion(ctx context.Context, tlsVersion TLSVersion, url string) (*int, error) {
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: uint16(tlsVersion),
			MaxVersion: uint16(tlsVersion),
		},
	}}

	if resp, err := ctxhttp.Get(ctx, client, url); err != nil {
		return nil, errorsutil.Wrapf(err, "tls version not supported (%s)", tlsVersion.String())
	} else {
		defer resp.Body.Close()
		return &resp.StatusCode, nil
	}
}

func HTTPResponseTLSVersion(r *http.Response) (TLSVersion, error) {
	if r == nil {
		return 0, errors.New("http.Response cannot be nil")
	} else if r.TLS == nil {
		return 0, errors.New("http.Response.TLS cannot be nil")
	} else {
		return TLSVersion(r.TLS.Version), nil
	}
}

func CheckURLs(urls []string) HTTPSVersionCheckResponse {
	res := NewHTTPSVersionCheckResponse()
	for _, url := range urls {
		res.Results = append(res.Results, CheckURL(url))
	}
	return res
}

func CheckURL(url string) URLResults {
	res := URLResults{
		URL:              url,
		TLSVersionChecks: []TLSVersionCheck{}}
	tlsVersions := TLSVersions()
	for _, tlsVersion := range tlsVersions {
		tlsVersionCheck := TLSVersionCheck{
			TLSVersion: tlsVersion.String(),
		}
		if statusCode, err := SupportsTLSVersion(context.Background(), tlsVersion, url); err != nil {
			tlsVersionCheck.Message = err.Error()
		} else {
			tlsVersionCheck.Supported = true
			tlsVersionCheck.HTTPStatusCode = statusCode
		}
		res.TLSVersionChecks = append(res.TLSVersionChecks, tlsVersionCheck)
	}
	return res
}

type HTTPSVersionCheckResponse struct {
	Results []URLResults `json:"results"`
}

func NewHTTPSVersionCheckResponse() HTTPSVersionCheckResponse {
	return HTTPSVersionCheckResponse{
		Results: []URLResults{},
	}
}

type URLResults struct {
	URL              string            `json:"url"`
	TLSVersionChecks []TLSVersionCheck `json:"tlsVersionChecks"`
}

type TLSVersionCheck struct {
	TLSVersion     string `json:"tlsVersion"`
	HTTPStatusCode *int   `json:"httpStatusCode"`
	Supported      bool   `json:"supported"`
	Message        string `json:"message"`
}
