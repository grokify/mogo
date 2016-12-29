package httputil

import (
	"fmt"
	"net"
	"net/http"
	"time"
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
