package jsonutil

import (
	"bytes"
	"encoding/json"
)

var (
	MarshalPrefix = ""
	MarshalIndent = "    "
)

type mustMarhshalError struct {
	MustMarhshalError string `json:"must_marshal_error"`
}

func MustMarshal(i interface{}, embedError bool) []byte {
	bytes, err := json.Marshal(i)
	if err != nil {
		if embedError {
			e := mustMarhshalError{
				MustMarhshalError: err.Error(),
			}
			bytes, err := json.Marshal(e)
			if err != nil {
				panic(err)
			}
			return bytes
		}
		panic(err)
	}
	return bytes
}

func MustMarshalString(i interface{}, embedError bool) string {
	return string(MustMarshal(i, embedError))
}

func MustMarshalIndent(i interface{}, prefix, indent string, embedError bool) []byte {
	bytes, err := json.MarshalIndent(i, prefix, indent)
	if err != nil {
		panic(err)
	}
	return bytes
}

func PrettyPrint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, MarshalPrefix, MarshalIndent)
	return out.Bytes(), err
}
