package simpleclient

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
)

func DoSimple(client *http.Client, httpMethod, requrl string, headers map[string][]string, body []byte) (*http.Response, error) {
	requrl = strings.TrimSpace(requrl)
	if len(requrl) == 0 {
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
		req, err = http.NewRequest(httpMethod, requrl, nil)
	} else {
		req, err = http.NewRequest(httpMethod, requrl, bytes.NewBuffer(body))
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

func Do(req SimpleRequest) (*http.Response, error) {
	sc := SimpleClient{}
	return sc.Do(req)
}

func DoMore(req SimpleRequest) ([]byte, *http.Response, error) {
	sc := SimpleClient{}
	resp, err := sc.Do(req)
	if err != nil {
		return []byte{}, resp, err
	}
	body, err := io.ReadAll(resp.Body)
	return body, resp, err
}
