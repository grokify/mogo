package httputil

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
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
