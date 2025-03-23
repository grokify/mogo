package jsonraw

import (
	"encoding/json"
	"fmt"
)

type unescapeWrap struct {
	Raw string `json:"raw"`
}

// Unescape is designed to unescape a stringified JSON. It is typically used
// after a stringified JSON has been embedded as a value in an wrapper JSON object.
// When using this, do not include outer quotes.
func Unescape(b []byte, prefix, indent string) ([]byte, error) {
	wrapped := fmt.Sprintf("{\"raw\":\"%s\"}", string(b))
	w := &unescapeWrap{}
	if err := json.Unmarshal([]byte(wrapped), w); err != nil {
		return nil, err
	} else if formatted, err := IndentBytes([]byte(w.Raw), prefix, indent); err != nil {
		return nil, err
	} else {
		return formatted, nil
	}
}
