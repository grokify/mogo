package httputilmore

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/grokify/gotilla/encoding/jsonutil"
	"github.com/pkg/errors"
)

// GetWriteFile gets the conents of a URL and stores the body in
// the desired filename location.
func GetWriteFile(client *http.Client, url, filename string) (*http.Response, error) {
	resp, err := client.Get(url)
	if err != nil {
		return resp, errors.Wrap(err, "httputilmore.GetStoreURL.client.Get()")
	}
	defer resp.Body.Close()
	dir, file := filepath.Split(filename)
	if len(strings.TrimSpace(dir)) > 0 {
		os.Chdir(dir)
	}
	f, err := os.Create(file)
	if err != nil {
		return resp, errors.Wrap(err, "httputilmore.GetStoreURL.os.Create()")
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		err = errors.Wrap(err, "httputilmore.GetStoreURL.io.Copy()")
	}
	return resp, err
}

// GetWriteFile performs a HTTP GET request and saves the response body
// to the file path specified. It reeads the entire file into memory
// which is not ideal for large files.
func GetWriteFileSimple(url string, filename string, perm os.FileMode) ([]byte, error) {
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
	_, err = jsonutil.UnmarshalIoReader(resp.Body, data)
	return resp, err
}

func PatchJsonBytes(client *http.Client, requrl string, headers map[string]string, body []byte) (*http.Response, error) {
	return RequestJsonBody(client, http.MethodPatch, requrl, headers, body)
}

func PostJsonBytes(client *http.Client, requrl string, headers map[string]string, body []byte) (*http.Response, error) {
	return RequestJsonBody(client, http.MethodPost, requrl, headers, body)
}

func PutJsonBytes(client *http.Client, requrl string, headers map[string]string, body []byte) (*http.Response, error) {
	return RequestJsonBody(client, http.MethodPut, requrl, headers, body)
}

func RequestJsonBody(client *http.Client, httpMethod, requrl string, headers map[string]string, body []byte) (*http.Response, error) {
	requrl = strings.TrimSpace(requrl)
	if len(requrl) == 0 {
		return nil, errors.New("E_NO_REQUEST_URL")
	}
	if client == nil {
		client = &http.Client{}
	}
	httpMethod = strings.TrimSpace(httpMethod)
	if httpMethod == "" {
		httpMethod = http.MethodPost
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
	for k, v := range headers {
		k = strings.TrimSpace(k)
		kMatch := strings.ToLower(k)
		if kMatch == strings.ToLower(HeaderContentType) {
			continue
		}
		req.Header.Set(k, v)
	}
	if len(body) > 0 {
		req.Header.Set(HeaderContentType, ContentTypeAppJsonUtf8)
	}

	return client.Do(req)
}

func DoJSON(client *http.Client, httpMethod, reqURL string, headers map[string]string, reqBody interface{}, resBody interface{}) (*http.Response, error) {
	bodyBytes := []byte("")
	var err error
	if reqBody != nil {
		bodyBytes, err = json.Marshal(reqBody)
	}
	if err != nil {
		return nil, err
	}

	resp, err := RequestJsonBody(client, httpMethod, reqURL, headers, bodyBytes)
	if err != nil || resBody == nil {
		return resp, err
	}
	if resBody != nil {
		_, err = jsonutil.UnmarshalIoReader(resp.Body, resBody)
	}
	return resp, err
}

/*
func PostJsonMarshal(client *http.Client, requrl string, headers map[string]string, body interface{}) (*http.Response, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return PostJsonBytes(client, requrl, headers, bodyBytes)
}

// PostJsonSimple performs a HTTP POST request converting a body interface{} to
// JSON and adding the appropriate JSON Content-Type header.
func PostJsonSimple(requrl string, body interface{}) (*http.Response, error) {
	return PostJsonMarshal(nil, requrl, map[string]string{}, body)
}
*/

// GetResponseAndBytes retreives a URL and returns the response body
// as a byte array in addition to the *http.Response.
func GetResponseAndBytes(url string) (*http.Response, []byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return resp, []byte{}, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
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
	req.Header.Add(HeaderContentType, ContentTypeAppFormUrlEncoded)
	req.Header.Add(HeaderContentLength, strconv.Itoa(len(data.Encode())))
	client := &http.Client{}
	return client.Do(req)
}

func Delete(client *http.Client, url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}
