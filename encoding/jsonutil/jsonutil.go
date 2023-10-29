package jsonutil

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"reflect"
)

const FileExt = ".json"

var (
	MarshalPrefix = ""
	MarshalIndent = "    "
)

type mustMarhshalError struct {
	MustMarhshalError string `json:"must_marshal_error"`
}

func MarshalSimple(v any, prefix, indent string) ([]byte, error) {
	if prefix == "" && indent == "" {
		return json.Marshal(v)
	}
	return json.MarshalIndent(v, prefix, indent)
}

func MustMarshalSimple(v any, prefix, indent string) []byte {
	bytes, err := MarshalSimple(v, prefix, indent)
	if err != nil {
		panic(err)
	}
	return bytes
}

func MustMarshal(v any, embedError bool) []byte {
	bytes, err := json.Marshal(v)
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

func MustMarshalString(v any, embedError bool) string {
	return string(MustMarshal(v, embedError))
}

func MustMarshalIndent(v any, prefix, indent string, embedError bool) []byte {
	bytes, err := json.MarshalIndent(v, prefix, indent)
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

func MarshalBase64(v any) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func MustUnmarshal(b []byte, v any) {
	err := json.Unmarshal(b, v)
	if err != nil {
		panic(err.Error())
	}
}

// UnmarshalAny will unmarshal anything to `v`, including first marshalling anything
// that is not a byte array to a JSON byte array.
func UnmarshalAny(data, v any) error {
	var err error
	b, ok := data.([]byte)
	if !ok {
		b, err = json.Marshal(data)
		if err != nil {
			return err
		}
	}
	return json.Unmarshal(b, v)
}

func UnmarshalMSI(data map[string]any, v any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, v)
}

func UnmarshalReader(r io.Reader, v any) ([]byte, error) {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return bytes, err
	}
	return bytes, json.Unmarshal(bytes, v)
}

// UnmarshalStrict returns an error when the destination is a struct and the input contains object keys which do not match any non-ignored, exported fields in the destination.
func UnmarshalStrict(b []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(b))
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

func ReadFile(filename string, v any) ([]byte, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return bytes, err
	}
	return bytes, json.Unmarshal(bytes, v)
}

func WriteFile(filename string, v any, prefix, indent string, perm fs.FileMode) error {
	bytes, err := MarshalSimple(v, prefix, indent)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, bytes, perm)
}

func Equal(x, y io.Reader) (bool, error) {
	var ax, ay any
	d := json.NewDecoder(x)
	if err := d.Decode(&ax); err != nil {
		return false, err
	}
	d = json.NewDecoder(y)
	if err := d.Decode(&ay); err != nil {
		return false, err
	}
	return reflect.DeepEqual(ax, ay), nil
}

func EqualBytes(x, y []byte) (bool, error) {
	var ax, ay any
	if err := json.Unmarshal(x, &ax); err != nil {
		return false, err
	}
	if err := json.Unmarshal(y, &ay); err != nil {
		return false, err
	}
	return reflect.DeepEqual(ax, ay), nil
}

func EqualFiles(x, y string) (bool, error) {
	fx, err := os.Open(x)
	if err != nil {
		return false, err
	}
	defer fx.Close()
	fy, err := os.Open(y)
	if err != nil {
		return false, err
	}
	defer fy.Close()
	return Equal(fx, fy)
}
