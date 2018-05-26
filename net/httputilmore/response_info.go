package httputilmore

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// ResponseInfo is a generic struct to handle response info.
type ResponseInfo struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Message    string            `json:"body"`
}

// ToJson returns ResponseInfo as a JSON byte array, embedding json.Marshal
// errors if encountered.
func (resIn *ResponseInfo) ToJson() []byte {
	bytes, err := json.Marshal(resIn)
	if err != nil {
		resIn2 := ResponseInfo{StatusCode: 500, Message: err.Error()}
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
