package httputilmore

import (
	"context"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Endpoint struct {
	Method HTTPMethod
	URL    *url.URL
}

// ParseEndpoint returns an `Endpoint` upon parsing a string like "POST https://example.com".
// If no method is provided, `GET` is returned. If the string has more than two fields, the lsat
// field is ignored.
func ParseEndpoint(s string) (*Endpoint, error) {
	parts := strings.Fields(s)
	if len(parts) == 0 {
		return nil, errors.New("empty string cannot be parsed as RUL")
	} else if len(parts) == 1 {
		if u, err := url.Parse(s); err != nil {
			return nil, err
		} else {
			return &Endpoint{
				Method: MethodGet,
				URL:    u,
			}, nil
		}
	} else {
		m, err := ParseHTTPMethod(parts[0])
		if err != nil {
			return nil, err
		}
		e := Endpoint{
			Method: m,
		}
		if u, err := url.Parse(parts[1]); err != nil {
			return nil, err
		} else {
			e.URL = u
			return &e, nil
		}
	}
}

func ParseRequestEndpoint(r *http.Request) *Endpoint {
	if r == nil {
		return nil
	}
	m := strings.ToUpper((strings.TrimSpace(r.Method)))
	if m == "" {
		m = http.MethodGet
	}
	return &Endpoint{
		Method: HTTPMethod(m),
		URL:    r.URL,
	}
}

func ParseRequestMethodPath(r *http.Request) string {
	if r == nil {
		return ""
	}
	m := strings.ToUpper((strings.TrimSpace(r.Method)))
	path := r.URL.Path
	return m + " " + path
}

// CreateProxyRequest creates a proxy request given a mapping "POST /path" => "POST https://newurl"
func CreateProxyRequest(m map[string]string, r *http.Request) (*httputil.ProxyRequest, error) {
	if r == nil {
		return nil, nil
	}
	mp := ParseRequestMethodPath(r)
	if v, ok := m[mp]; ok {
		ep, err := ParseEndpoint(v)
		if err != nil {
			return nil, err
		}
		pr := httputil.ProxyRequest{
			In: r,
		}
		out := r.Clone(context.TODO())
		out.Method = string(ep.Method)
		out.URL = ep.URL
		return &pr, nil
	} else {
		return nil, nil
	}
}
