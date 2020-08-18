package httputilmore

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

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

func ConsolidateErrorRespCodeGte300(resp *http.Response, err error, msg string) error {
	if err != nil {
		return errors.Wrap(err, msg)
	} else if resp.StatusCode >= 300 {
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
		return fmt.Errorf("%sStatusCode [%v] %s", msg, resp.StatusCode, moreString)
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

// ToJson returns ResponseInfo as a JSON byte array, embedding json.Marshal
// errors if encountered.
func (resIn *ResponseInfo) ToJson() []byte {
	bytes, err := json.Marshal(resIn)
	if err != nil {
		resIn2 := ResponseInfo{StatusCode: 500, Body: err.Error()}
		bytes, _ := json.Marshal(resIn2)
		return bytes
	}
	return bytes
}
