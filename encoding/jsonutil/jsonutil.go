package jsonutil

import (
	"bytes"
	"encoding/json"
)

var (
	MarshalPrefix = ""
	MarshalIndent = "    "
)

func PrettyPrint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, MarshalPrefix, MarshalIndent)
	return out.Bytes(), err
}
