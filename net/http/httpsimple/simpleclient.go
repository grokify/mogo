package httpsimple

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/grokify/mogo/io/ioutil"
	"github.com/grokify/mogo/net/urlutil"
	"golang.org/x/net/context/ctxhttp"
)

var rxHTTPURL = regexp.MustCompile(`^(?i)https?://`)

// Client provides a simple interface to making HTTP requests using `net/http`.
type Client struct {
	BaseURL    string     // TODO: See if we can do this within Transport: https://gist.github.com/epelc/cc286ad0fd7878fb176a89f2af1177b6
	Query      url.Values // Add If Not Exists
	HTTPClient *http.Client
}

func NewClient(httpClient *http.Client, baseURL string) Client {
	return Client{HTTPClient: httpClient, BaseURL: baseURL}
}

func (sc *Client) Get(ctx context.Context, reqURL string) (*http.Response, error) {
	return sc.Do(ctx, Request{Method: http.MethodGet, URL: reqURL})
}

func (sc *Client) Do(ctx context.Context, req Request) (*http.Response, error) {
	req.Inflate()
	var bodyReader io.Reader
	if req.BodyType == BodyTypeFile && ioutil.IsReader(req.Body) {
		if reqBodyReader, ok := req.Body.(io.Reader); ok {
			bodyReader = reqBodyReader
		} else {
			panic("cannot cast `io.Reader` as `io.Reader`")
		}
	} else {
		if bodyBytes, err := req.BodyBytes(); err != nil {
			return nil, err
		} else {
			bodyReader = bytes.NewReader(bodyBytes)
		}
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
	return doSimple(ctx, sc.HTTPClient, req.Method, reqURL, req.Headers, bodyReader)
}

func (sc *Client) DoUnmarshalJSON(ctx context.Context, req Request, resBody any) ([]byte, *http.Response, error) {
	resp, err := sc.Do(ctx, req)
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

func doSimple(ctx context.Context, client *http.Client, httpMethod, reqURL string, headers map[string][]string, body io.Reader) (*http.Response, error) {
	// func doSimple(client *http.Client, httpMethod, reqURL string, headers map[string][]string, body []byte) (*http.Response, error) {
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
	//var req *http.Request
	//var err error

	//if len(body) == 0 {
	req, err := http.NewRequest(httpMethod, reqURL, body)
	//} else {
	//	req, err = http.NewRequest(httpMethod, reqURL, bytes.NewBuffer(body))
	//}
	if err != nil {
		return nil, err
	}
	for k, vals := range headers {
		for _, v := range vals {
			req.Header.Set(k, v)
		}
	}

	return ctxhttp.Do(ctx, client, req)
}
