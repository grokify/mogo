package jsonutil

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/derekstavis/go-qs"
)

var (
	MarshalPrefix = ""
	MarshalIndent = "    "
)

type mustMarhshalError struct {
	MustMarhshalError string `json:"must_marshal_error"`
}

func MarshalSimple(v interface{}, prefix, indent string) ([]byte, error) {
	if prefix == "" && indent == "" {
		return json.Marshal(v)
	}
	return json.MarshalIndent(v, prefix, indent)
}

func MustMarshalSimple(v interface{}, prefix, indent string) []byte {
	bytes, err := MarshalSimple(v, prefix, indent)
	if err != nil {
		panic(err)
	}
	return bytes
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

func MarshalBase64(i interface{}) (string, error) {
	data, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func UnmarshalIoReader(r io.Reader, i interface{}) ([]byte, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return b, err
	}
	return b, json.Unmarshal(b, i)
}

// UnmarshalRailsQS unmarshals a Rails query string to a Go struct.
func UnmarshalRailsQS(railsQueryString string, i interface{}) error {
	query, err := qs.Unmarshal(railsQueryString)
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(query)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, i)
}
