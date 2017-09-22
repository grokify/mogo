package httputilmore

import (
	"fmt"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	DIAL_TIMEOUT int64 = 5
	TLS_TIMEOUT  int   = 5
	HTTP_TIMEOUT int   = 10
)

// The default Go HTTP client never times out.
// This HTTP client provides default and updatable timeouts
// More at: https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779#.ymd655pgz
func NewHttpClient() *http.Client {
	dial_timeout, _ := time.ParseDuration(fmt.Sprintf("%vs", DIAL_TIMEOUT))

	netTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: dial_timeout}).Dial,
		TLSHandshakeTimeout: 5 * time.Second}

	netClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: netTransport}

	return netClient
}

func GetRequestRateLimited(client *http.Client, reqURL string) (*http.Response, error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return &http.Response{}, err
	}
	return DoRequestRateLimited(client, req)
}

func DoRequestRateLimited(client *http.Client, req *http.Request) (*http.Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}

	rlstat := NewResponseRateLimitInfo(resp, true)

	if rlstat.XRateLimitRemaining <= 0 {
		log.WithFields(log.Fields{
			"action":                 "http_rate_limited",
			"status_code":            rlstat.StatusCode,
			"retry-after":            rlstat.RetryAfter,
			"x-rate-limit-remaining": rlstat.XRateLimitRemaining,
			"x-rate-limit-window":    rlstat.XRateLimitWindow,
		}).Info("Request has been rated limited.")
		time.Sleep(time.Duration(rlstat.XRateLimitWindow) * time.Second)
		return resp, nil
	} else if rlstat.StatusCode == 429 {
		log.WithFields(log.Fields{
			"action":                 "http_rate_limited",
			"status_code":            rlstat.StatusCode,
			"retry-after":            rlstat.RetryAfter,
			"x-rate-limit-remaining": rlstat.XRateLimitRemaining,
			"x-rate-limit-window":    rlstat.XRateLimitWindow,
		}).Info("Request has been rated limited.")
		if rlstat.RetryAfter > 0 {
			time.Sleep(time.Duration(rlstat.RetryAfter) * time.Second)
		} else if rlstat.XRateLimitWindow > 0 {
			time.Sleep(time.Duration(rlstat.XRateLimitWindow) * time.Second)
		} else {
			time.Sleep(time.Duration(60) * time.Second)
		}
		return client.Do(req)
	}
	return resp, nil
}
