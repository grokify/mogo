package httputilmore

import (
	"errors"
	"io/ioutil"
	"net/http"
)

// ProxyResponse copies the information from a `*http.Response` to a
// `http.ResponseWriter`.
func ProxyResponse(w http.ResponseWriter, resp *http.Response) error {
	if resp == nil {
		return errors.New("E_NIL_HTTP_RESPONSE")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_, err = w.Write(body)
	if err != nil {
		return err
	}
	w.Header().Add(HeaderContentType, resp.Header.Get(HeaderContentType))
	w.WriteHeader(resp.StatusCode)
	return nil
}
