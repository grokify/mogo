package httputilmore

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

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

func ConsolidateErrorRespCodeGte300(resp *http.Response, err error, msg string) error {
	if err != nil {
		return errors.Wrap(err, msg)
	} else if resp.StatusCode >= 300 {
		if len(msg) > 0 {
			msg += ": "
		}
		return fmt.Errorf("%sStatusCode [%v]", msg, resp.StatusCode)
	}
	return nil
}
