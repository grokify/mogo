package simpleclient

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/grokify/gotilla/net/httputilmore"
	"github.com/grokify/gotilla/net/urlutil"
)

var rxHttpUrl = regexp.MustCompile(`^(?i)https?://`)

type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    interface{}
	IsJSON  bool
}

func (req *Request) Inflate() {
	req.Method = strings.ToUpper(strings.TrimSpace(req.Method))
	if len(req.Method) == 0 {
		req.Method = http.MethodGet
	}
	if req.Headers == nil {
		req.Headers = map[string]string{}
	}
	if req.IsJSON {
		if _, ok := req.Headers[httputilmore.HeaderContentType]; !ok {
			req.Headers[httputilmore.HeaderContentType] = httputilmore.ContentTypeAppJsonUtf8
		}
	}
}

func (req *Request) BodyBytes() ([]byte, error) {
	if req.Body == nil {
		return []byte(""), nil
	} else if reqBodyBytesAssert, ok := req.Body.([]byte); ok {
		return reqBodyBytesAssert, nil
	} else if reqBodyString, ok := req.Body.(string); ok {
		return []byte(reqBodyString), nil
	}
	return json.Marshal(req.Body)
}

// SimpleClient provides a simple interface to making HTTP requests
// using `net/http`.
type SimpleClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewSimpleClient(httpClient *http.Client, baseURL string) SimpleClient {
	return SimpleClient{HTTPClient: httpClient, BaseURL: baseURL}
}

func (sc *SimpleClient) Get(reqURL string) (*http.Response, error) {
	return sc.Do(Request{Method: http.MethodGet, URL: reqURL})
}

func (sc *SimpleClient) Do(req Request) (*http.Response, error) {
	req.Inflate()
	bodyBytes, err := req.BodyBytes()
	if err != nil {
		return nil, err
	}
	reqURL := strings.TrimSpace(req.URL)
	if !rxHttpUrl.MatchString(reqURL) && len(reqURL) > 0 {
		reqURL = urlutil.JoinAbsolute(sc.BaseURL, reqURL)
	}
	return httputilmore.DoJSONSimple(
		sc.HTTPClient, req.Method, reqURL, req.Headers, bodyBytes)
}
