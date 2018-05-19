package httputilmore

import (
	"encoding/json"
)

// ResponseInfo is a generic struct to handle response info.
type ResponseInfo struct {
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Message    string            `json:"message,omitempty"`
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
