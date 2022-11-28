package jsonutil

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
)

const FileExt = ".json"

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

// IndentBytes converts a JSON byte array into a prettified byte array.
func IndentBytes(b []byte, prefix, indent string) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, prefix, indent)
	if err != nil {
		return []byte{}, err
	}
	return out.Bytes(), nil
}

// IndentReader returns a byte slice of indented JSON given an `io.Reader`.
// It is useful to use with `http.Response.Body` which is an `io.ReadCloser`.
func IndentReader(r io.Reader, prefix, indent string) ([]byte, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return b, err
	}
	return IndentBytes(b, prefix, indent)
}

func MarshalBase64(i interface{}) (string, error) {
	data, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func MustUnmarshal(data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	if err != nil {
		panic(err.Error())
	}
}

func UnmarshalMSI(data map[string]interface{}, v interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, v)
}

func UnmarshalReader(r io.Reader, v interface{}) ([]byte, error) {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return bytes, err
	}
	return bytes, json.Unmarshal(bytes, v)
}

func UnmarshalStrict(data []byte, v interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// PrintReaderIndent returns an indented JSON byte array given an `io.Reader`.
func PrintReaderIndent(r io.Reader, prefix, indent string) ([]byte, error) {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return bytes, err
	}
	outBytes, err := IndentBytes(bytes, prefix, indent)
	if err != nil {
		return bytes, err
	}
	_, err = fmt.Println(string(outBytes))
	return outBytes, err
}

func ReadFile(filename string, v interface{}) ([]byte, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return bytes, err
	}
	return bytes, json.Unmarshal(bytes, v)
}

func WriteFile(filename string, v interface{}, prefix, indent string, perm fs.FileMode) error {
	bytes, err := MarshalSimple(v, prefix, indent)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, bytes, perm)
}
