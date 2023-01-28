package httputilmore

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strings"

	"github.com/grokify/mogo/strconv/strconvutil"
)

/*
// UnmarshalResponseJSON is EOL, use jsonutil.UnmarshalIoReader()
func UnmarshalResponseJSON(resp *http.Response, data any) ([]byte, error) {
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return bytes, err
	}
	return bytes, json.Unmarshal(bytes, data)
}*/

// ResponseBodyJSONMapIndent returns the body as a generic JSON dictionary
func ResponseBodyJSONMapIndent(res *http.Response, prefix string, indent string) ([]byte, error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return body, err
	}
	any := map[string]any{}
	err = json.Unmarshal(body, &any)
	if err != nil {
		return body, err
	}
	return json.MarshalIndent(any, prefix, indent)
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

// ParseHeader converts a raw strign to a header struct.
func ParseHeader(s string) http.Header {
	h := http.Header{}
	lines := strings.Split(s, "\n")
	rx := regexp.MustCompile(`^([^\s+]+):\s*(.*)$`)
	for _, line := range lines {
		m := rx.FindStringSubmatch(line)
		if len(m) == 3 {
			key := strings.TrimSpace(m[1])
			val := strings.TrimSpace(m[2])
			if len(key) > 0 {
				h.Add(key, val)
			}
		}
	}
	return h
}

// MergeHeader merges two http.Header adding the values of the second
// to the first.
func MergeHeader(base, more http.Header, overwrite bool) http.Header {
	if base == nil {
		base = http.Header{}
	}
	if more == nil {
		return base
	}
	for k, vals := range more {
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

func ParseMultipartFormDataBoundaryFromHeader(contentType string) string {
	rx := regexp.MustCompile(`^multipart/form-data.+boundary="?([^;"]+)`)
	m := rx.FindStringSubmatch(contentType)
	if len(m) > 0 {
		return m[1]
	}
	return ""
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
		RetryAfter: strconvutil.AtoiOrDefault(resp.Header.Get("Retry-After"), 0)}

	if useXrlHyphen {
		rlstat.XRateLimitLimit = strconvutil.AtoiOrDefault(resp.Header.Get("X-Rate-Limit-Limit"), 0)
		rlstat.XRateLimitRemaining = strconvutil.AtoiOrDefault(resp.Header.Get("X-Rate-Limit-Remaining"), 0)
		rlstat.XRateLimitReset = strconvutil.AtoiOrDefault(resp.Header.Get("X-Rate-Limit-Reset"), 0)
		rlstat.XRateLimitWindow = strconvutil.AtoiOrDefault(resp.Header.Get("X-Rate-Limit-Window"), 0)
	} else {
		rlstat.XRateLimitLimit = strconvutil.AtoiOrDefault(resp.Header.Get("X-RateLimit-Limit"), 0)
		rlstat.XRateLimitRemaining = strconvutil.AtoiOrDefault(resp.Header.Get("X-RateLimit-Remaining"), 0)
		rlstat.XRateLimitReset = strconvutil.AtoiOrDefault(resp.Header.Get("X-RateLimit-Reset"), 0)
		rlstat.XRateLimitWindow = strconvutil.AtoiOrDefault(resp.Header.Get("X-RateLimit-Window"), 0)
	}
	return rlstat
}

type FnLogRateLimitInfo func(RateLimitInfo)
