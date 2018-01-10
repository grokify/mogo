package httputilmore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/grokify/gotilla/strconv/strconvutil"
)

const (
	// ContentTypeHeader is the content type header name
	ContentTypeHeader = "Content-Type"
	// ContentTypeValueJSONUTF8 represents the HTTP Content-Type header
	// value for UTF-8 encoded JSON.
	ContentTypeValueJSONUTF8 = "application/json; charset=utf-8"
	// HTTPS Scheme
	SchemeHTTPS = "https"
)

// GetWriteFile performs a HTTP GET request and saves the response body
// to the file path specified
func GetWriteFile(url string, filename string, perm os.FileMode) ([]byte, error) {
	_, bytes, err := GetResponseAndBytes(url)
	if err != nil {
		return bytes, err
	}
	err = ioutil.WriteFile(filename, bytes, perm)
	return bytes, err
}

// ResponseBody returns the body as a byte array
func ResponseBody(res *http.Response) ([]byte, error) {
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	return contents, nil
}

// ResponseBodyJSONMapIndent returns the body as a generic JSON dictionary
func ResponseBodyJSONMapIndent(res *http.Response, prefix string, indent string) ([]byte, error) {
	body, err := ResponseBody(res)
	if err != nil {
		return body, err
	}
	any := map[string]interface{}{}
	json.Unmarshal(body, &any)
	return json.MarshalIndent(any, prefix, indent)
}

// GetResponseAndBytes retreives a URL and returns the response body
// as a byte array in addition to the *http.Response.
func GetResponseAndBytes(url string) (*http.Response, []byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return resp, []byte{}, err
	}
	bytes, err := ResponseBody(resp)
	return resp, bytes, err
}

// UnmarshalResponseJSON unmarshal a `*http.Response` JSON body into
// a data pointer.
func UnmarshalResponseJSON(resp *http.Response, data interface{}) error {
	//bytes, err := ResponseBody(resp)
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, data)
}

// PrintRequestOut prints a http.Request using `httputil.DumpRequestOut`.
func PrintRequestOut(req *http.Request, includeBody bool) error {
	reqBytes, err := httputil.DumpRequestOut(req, includeBody)
	if err != nil {
		return err
	}
	fmt.Println(string(reqBytes))
	return nil
}

// PrintResponse prints a http.Response using `httputil.DumpResponse`.
func PrintResponse(resp *http.Response, includeBody bool) error {
	respBytes, err := httputil.DumpResponse(resp, includeBody)
	if err != nil {
		return err
	}
	fmt.Println(string(respBytes))
	return nil
}

// MergeHeader merges two http.Header adding the values of the second
// to the first.
func MergeHeader(base, extra http.Header, overwrite bool) http.Header {
	for k, vals := range extra {
		if overwrite {
			base.Del(k)
		}

		for _, v := range vals {
			v = strings.TrimSpace(v)
			if len(v) > 0 {
				base.Add(k, v)
			}
		}
	}
	return base
}

// RateLimitInfo is a structure for holding parsed rate limit info.
// It uses headers from the GitHub, RingCentral and Twitter APIs.
type RateLimitInfo struct {
	StatusCode          int
	RetryAfter          int
	XRateLimitLimit     int
	XRateLimitRemaining int
	XRateLimitReset     int
	XRateLimitWindow    int
}

// NewResponseRateLimitInfo returns a RateLimitInfo from a http.Response.
func NewResponseRateLimitInfo(resp *http.Response, useXrlHyphen bool) RateLimitInfo {
	rlstat := RateLimitInfo{
		StatusCode: resp.StatusCode,
		RetryAfter: strconvutil.AtoiWithDefault(resp.Header.Get("Retry-After"), 0)}

	if useXrlHyphen {
		rlstat.XRateLimitLimit = strconvutil.AtoiWithDefault(resp.Header.Get("X-Rate-Limit-Limit"), 0)
		rlstat.XRateLimitRemaining = strconvutil.AtoiWithDefault(resp.Header.Get("X-Rate-Limit-Remaining"), 0)
		rlstat.XRateLimitReset = strconvutil.AtoiWithDefault(resp.Header.Get("X-Rate-Limit-Reset"), 0)
		rlstat.XRateLimitWindow = strconvutil.AtoiWithDefault(resp.Header.Get("X-Rate-Limit-Window"), 0)
	} else {
		rlstat.XRateLimitLimit = strconvutil.AtoiWithDefault(resp.Header.Get("X-RateLimit-Limit"), 0)
		rlstat.XRateLimitRemaining = strconvutil.AtoiWithDefault(resp.Header.Get("X-RateLimit-Remaining"), 0)
		rlstat.XRateLimitReset = strconvutil.AtoiWithDefault(resp.Header.Get("X-RateLimit-Reset"), 0)
		rlstat.XRateLimitWindow = strconvutil.AtoiWithDefault(resp.Header.Get("X-RateLimit-Window"), 0)
	}
	return rlstat
}

type FnLogRateLimitInfo func(RateLimitInfo)
