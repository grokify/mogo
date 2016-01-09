package httputil

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

// GetWriteFile performs a HTTP GET request and saves the response body
// to the file path specified

func GetWriteFile(url string, filename string, perm os.FileMode) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	bytes, err := ResponseBody(resp)
	if err != nil {
		return []byte{}, err
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
	} else {
		return contents, nil
	}
}

// ResponseBodyJsonMapIndent returns the body as a generic JSON dictionary

func ResponseBodyJsonMapIndent(res *http.Response, prefix string, indent string) ([]byte, error) {
	body, err := ResponseBody(res)
	if err != nil {
		return body, err
	}
	any := map[string]interface{}{}
	json.Unmarshal(body, &any)
	return json.MarshalIndent(any, prefix, indent)
}
