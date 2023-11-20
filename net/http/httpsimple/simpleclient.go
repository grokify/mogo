package httpsimple

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/grokify/mogo/net/urlutil"
)

var rxHTTPURL = regexp.MustCompile(`^(?i)https?://`)

// Client provides a simple interface to making HTTP requests using `net/http`.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewSimpleClient(httpClient *http.Client, baseURL string) Client {
	return Client{HTTPClient: httpClient, BaseURL: baseURL}
}

func (sc *Client) Get(reqURL string) (*http.Response, error) {
	return sc.Do(Request{Method: http.MethodGet, URL: reqURL})
}

func (sc *Client) Do(req Request) (*http.Response, error) {
	req.Inflate()
	bodyBytes, err := req.BodyBytes()
	if err != nil {
		return nil, err
	}
	reqURL := strings.TrimSpace(req.URL)
	if len(sc.BaseURL) > 0 {
		if len(reqURL) == 0 {
			reqURL = sc.BaseURL
		} else if !rxHTTPURL.MatchString(reqURL) {
			reqURL = urlutil.JoinAbsolute(sc.BaseURL, reqURL)
		}
	}
	if len(req.Query) > 0 {
		goURL, err := urlutil.URLAddQueryString(reqURL, req.Query)
		if err != nil {
			return nil, err
		}
		reqURL = goURL.String()
	}
	if sc.HTTPClient == nil {
		sc.HTTPClient = &http.Client{}
	}
	return doSimple(sc.HTTPClient, req.Method, reqURL, req.Headers, bodyBytes)
}

func (sc *Client) DoUnmarshalJSON(req Request, resBody any) ([]byte, *http.Response, error) {
	resp, err := sc.Do(req)
	if err != nil {
		return []byte{}, nil, err
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return bytes, resp, err
	}
	err = json.Unmarshal(bytes, resBody)
	return bytes, resp, err
}
