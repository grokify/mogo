package jsonutil

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	// jsoniter "github.com/json-iterator/go"
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	EmptyArray  = "[]"
	EmptyObject = "{}"
	FileExt     = ".json"
)

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
	} else {
		return json.MarshalIndent(v, prefix, indent)
	}
}

func MustMarshalSimple(v any, prefix, indent string) []byte {
	if b, err := MarshalSimple(v, prefix, indent); err != nil {
		panic(err)
	} else {
		return b
	}
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

// MarshalOrDefault returns the supplied default value in the event
// of an error. It also returns the error for processing.
func MarshalOrDefault(v any, def []byte) ([]byte, error) {
	if b, err := json.Marshal(v); err != nil {
		return def, err
	} else {
		return b, err
	}
}

// MarshalSlice marshals any slice, using `fmt.Sprintf()` for non-strings.
func MarshalSlice[T any](v []T, stripBrackets bool) ([]byte, error) {
	var out []string
	for _, vi := range v {
		via := any(vi)
		if s, ok := via.(string); ok {
			out = append(out, s)
		} else {
			out = append(out, fmt.Sprintf("%v", vi))
		}
	}
	if b, err := json.Marshal(out); err != nil {
		return []byte{}, err
	} else if !stripBrackets {
		return b, nil
	} else {
		return b[1 : len(b)-1], nil
	}
}

func MustMarshalOrDefault(v any, def []byte) []byte {
	if b, err := json.Marshal(v); err != nil {
		return def
	} else {
		return b
	}
}

func MustMarshalIndent(v any, prefix, indent string, embedError bool) []byte {
	if b, err := json.MarshalIndent(v, prefix, indent); err != nil {
		panic(err)
	} else {
		return b
	}
}

// ValidateQuick checks to see if the first and last bytes match `{}[]`.
// Set `fuzzy` to true to trim spaces from the beginning and end. It is not
// Design to provide full validation but quick decision making on whether
// to attempt JSON parsing.
func ValidateQuick(b []byte, fuzzy bool) bool {
	if fuzzy {
		s := strings.TrimSpace(string(b))
		if (strings.Index(s, "{") == 0 && strings.HasSuffix(s, "}")) ||
			(strings.Index(s, "[") == 0 && strings.HasSuffix(s, "]")) {
			return true
		} else {
			return false
		}
	} else {
		if len(b) < 2 {
			return false
		} else if (b[0] == 91 && b[len(b)-1] == 93) ||
			(b[0] == 123 && b[len(b)-1] == 125) {
			return true
		} else {
			return false
		}
	}
}

func MarshalBase64(v any) (string, error) {
	if data, err := json.Marshal(v); err != nil {
		return "", err
	} else {
		return base64.StdEncoding.EncodeToString(data), nil
	}
}

func MustUnmarshal(b []byte, v any) {
	if err := json.Unmarshal(b, v); err != nil {
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
	if b, err := json.Marshal(data); err != nil {
		return err
	} else {
		return json.Unmarshal(b, v)
	}
}

func UnmarshalReader(r io.Reader, v any) ([]byte, error) {
	if b, err := io.ReadAll(r); err != nil {
		return b, err
	} else {
		return b, json.Unmarshal(b, v)
	}
}

// UnmarshalStrict returns an error when the destination is a struct and the input contains object keys which do not match any non-ignored, exported fields in the destination.
func UnmarshalStrict(b []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

func UnmarshalFile(filename string, v any) error {
	if f, err := os.Open(filename); err != nil {
		return err
	} else {
		defer f.Close()
		decr := json.NewDecoder(f)
		return decr.Decode(v)
	}
}

func UnmarshalFileWithBytes(filename string, v any) ([]byte, error) {
	if b, err := os.ReadFile(filename); err != nil {
		return b, err
	} else {
		return b, json.Unmarshal(b, v)
	}
}

func MarshalFile(filename string, v any, prefix, indent string, perm fs.FileMode) error {
	if f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm); err != nil {
		return err
	} else {
		defer f.Close()
		encr := json.NewEncoder(f)
		encr.SetIndent(prefix, indent)
		return encr.Encode(v)
	}
}
