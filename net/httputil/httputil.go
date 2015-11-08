package httputil

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

func GetWriteFile(url string, filepath string, perm os.FileMode) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	bytes, err := ResponseBody(resp)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath, bytes, perm)
	return err
}

func ResponseBody(res *http.Response) ([]byte, error) {
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	} else {
		return contents, nil
	}
}

func ResponseBodyJsonMapIndent(res *http.Response, prefix string, indent string) ([]byte, error) {
	body, err := ResponseBody(res)
	if err != nil {
		return body, err
	}
	any := map[string]interface{}{}
	json.Unmarshal(body, &any)
	return json.MarshalIndent(any, prefix, indent)
}
