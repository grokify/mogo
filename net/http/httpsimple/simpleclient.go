package httpsimple

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/grokify/mogo/net/urlutil"
)

var rxHTTPURL = regexp.MustCompile(`^(?i)https?://`)

// Client provides a simple interface to making HTTP requests using `net/http`.
type Client struct {
	BaseURL    string
	Query      url.Values // Add If Not Exists
	HTTPClient *http.Client
}

func NewClient(httpClient *http.Client, baseURL string) Client {
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
	qry := urlutil.AppendValues(req.Query, sc.Query, false)
	if len(qry) > 0 {
		goURL, err := urlutil.URLStringAddQuery(reqURL, qry, false)
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

func doSimple(client *http.Client, httpMethod, reqURL string, headers map[string][]string, body []byte) (*http.Response, error) {
	reqURL = strings.TrimSpace(reqURL)
	if len(reqURL) == 0 {
		return nil, errors.New("requrl is required but not present")
	}
	if client == nil {
		client = &http.Client{}
	}
	httpMethod = strings.TrimSpace(httpMethod)
	if httpMethod == "" {
		return nil, errors.New("httpMethod is required but not present")
	}
	var req *http.Request
	var err error

	if len(body) == 0 {
		req, err = http.NewRequest(httpMethod, reqURL, nil)
	} else {
		req, err = http.NewRequest(httpMethod, reqURL, bytes.NewBuffer(body))
	}
	if err != nil {
		return nil, err
	}
	for k, vals := range headers {
		for _, v := range vals {
			req.Header.Set(k, v)
		}
	}

	return client.Do(req)
}

func Do(req Request) (*http.Response, error) {
	sc := Client{}
	return sc.Do(req)
}

func DoMore(req Request) ([]byte, *http.Response, error) {
	sc := Client{}
	resp, err := sc.Do(req)
	if err != nil {
		return []byte{}, resp, err
	}
	body, err := io.ReadAll(resp.Body)
	return body, resp, err
}
