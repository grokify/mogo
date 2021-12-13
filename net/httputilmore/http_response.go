package httputilmore

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/pkg/errors"
)

// ProxyResponse copies the information from a `*http.Response` to a
// `http.ResponseWriter`.
func ProxyResponse(w http.ResponseWriter, resp *http.Response) ([]byte, error) {
	if resp == nil {
		return nil, errors.New("E_NIL_HTTP_RESPONSE")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	_, err = w.Write(body)
	if err != nil {
		return body, err
	}
	headers := []string{
		HeaderContentEncoding,
		HeaderContentLanguage,
		HeaderContentTransferEncoding,
		HeaderContentType}
	for _, header := range headers {
		if len(resp.Header.Get(header)) > 0 {
			w.Header().Add(header, resp.Header.Get(header))
		}
	}
	w.WriteHeader(resp.StatusCode)
	return body, nil
}

func CondenseResponseNot2xxToError(resp *http.Response, err error, msg string) error {
	if err != nil {
		if len(msg) > 0 {
			return errors.Wrap(err, msg)
		} else {
			return err
		}
	} else if resp == nil {
		return errors.New("*http.Response_is_nil")
	} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
		if len(msg) > 0 {
			msg += ": "
		}
		more := []string{}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			more = append(more, err.Error())
		} else {
			more = append(more, string(body))
		}
		moreString := ""
		jsn, err := json.Marshal(more)
		if err != nil {
			moreString = err.Error()
		} else {
			moreString = string(jsn)
		}
		return fmt.Errorf("non_2xx_status_code [%d] [%s] [%s]", resp.StatusCode, msg, moreString)
	}
	return nil
}

// ResponseInfo is a generic struct to handle response info.
type ResponseInfo struct {
	Name       string            `json:"name,omitempty"` // to distinguish from other requests
	Method     string            `json:"method,omitempty"`
	URL        string            `json:"url,omitempty"`
	StatusCode int               `json:"statusCode,omitempty"`
	Time       time.Time         `json:"time,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
}

// ToJSON returns ResponseInfo as a JSON byte array, embedding json.Marshal
// errors if encountered.
func (resIn *ResponseInfo) ToJSON() []byte {
	bytes, err := json.Marshal(resIn)
	if err != nil {
		resIn2 := ResponseInfo{StatusCode: 500, Body: err.Error()}
		bytes, _ := json.Marshal(resIn2)
		return bytes
	}
	return bytes
}

func ResponseWriterWriteJSON(w http.ResponseWriter, statusCode int, body interface{}, prefix, indent string) {
	if w == nil {
		return
	}
	wroteStatus := false
	if body != nil {
		bytes, err := jsonutil.MarshalSimple(body, "", "  ")
		if err != nil {
			gr := ResponseInfo{
				Body: err.Error()}
			bytes, err = jsonutil.MarshalSimple(gr, "", "  ")
		}
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		_, err = w.Write(bytes)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			wroteStatus = true
		}
	}
	if !wroteStatus && statusCode != 0 {
		w.WriteHeader(statusCode)
	}
}
