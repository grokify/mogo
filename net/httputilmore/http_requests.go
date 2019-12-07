package httputilmore

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/grokify/gotilla_bak/net/httputilmore"
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

func GetJsonSimple(requrl string, header http.Header, data interface{}) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, requrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header = header
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	err = UnmarshalResponseJSON(resp, data)
	return resp, err
}

// PostJsonSimple performs a HTTP POST request converting a body interface{} to
// JSON and adding the appropriate JSON Content-Type header.
func PostJsonSimple(requrl string, body interface{}) (*http.Response, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return &http.Response{}, err
	}

	req, err := http.NewRequest(http.MethodPost, requrl, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return &http.Response{}, err
	}
	req.Header.Set(HeaderContentType, ContentTypeAppJsonUtf8)

	client := &http.Client{}
	return client.Do(req)
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

func SendWwwFormUrlEncodedSimple(method, urlStr string, data url.Values) (*http.Response, error) {
	req, err := http.NewRequest(
		method,
		urlStr,
		strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		return nil, err
	}
	req.Header.Add(httputilmore.HeaderContentType, httputilmore.ContentTypeAppFormUrlEncoded)
	req.Header.Add(httputilmore.HeaderContentLength, strconv.Itoa(len(data.Encode())))
	client := &http.Client{}
	return client.Do(req)
}
