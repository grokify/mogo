package anyhttp

import (
	"io"
	"mime/multipart"
	"net/url"

	hum "github.com/grokify/gotilla/net/httputilmore"
)

type Request interface {
	Method() []byte
	ParseForm() error
	AllArgs() Args
	QueryArgs() Args
	PostArgs() Args
	MultipartForm() (*multipart.Form, error)
}

type Args interface {
	GetBytes(key string) []byte
	GetBytesSlice(key string) [][]byte
	GetString(key string) string
	GetStringSlice(key string) []string
}

type Response interface {
	SetStatusCode(int)
	SetContentType(string)
	SetBodyBytes([]byte) (int, error)
	SetBodyStream(bodyStream io.Reader, bodySize int) error
}

func WriteSimpleJson(w Response, status int, message string) {
	w.SetStatusCode(status)
	w.SetContentType(hum.ContentTypeAppJsonUtf8)
	resInfo := hum.ResponseInfo{
		StatusCode: status,
		Message:    message}
	w.SetBodyBytes(resInfo.ToJson())
}

type MapStringString map[string]string

func (m MapStringString) Get(key string) string {
	if val, ok := m[key]; ok {
		return val
	}
	return ""
}

func (m MapStringString) GetSlice(key string) []string {
	return []string{m.Get(key)}
}

type ArgsMapStringString struct{ Raw MapStringString }

func NewArgsMapStringString(args MapStringString) ArgsMapStringString {
	return ArgsMapStringString{Raw: args}
}

func (args ArgsMapStringString) GetBytes(key string) []byte { return []byte(args.Raw.Get(key)) }
func (args ArgsMapStringString) GetBytesSlice(key string) [][]byte {
	output := make([][]byte, 1)
	output[0] = args.GetBytes(key)
	return output
}
func (args ArgsMapStringString) GetString(key string) string        { return args.Raw.Get(key) }
func (args ArgsMapStringString) GetStringSlice(key string) []string { return args.Raw.GetSlice(key) }

type ArgsUrlValues struct{ Raw url.Values }

func NewArgsUrlValues(args url.Values) ArgsUrlValues {
	return ArgsUrlValues{Raw: args}
}

func (args ArgsUrlValues) GetBytes(key string) []byte { return []byte(args.Raw.Get(key)) }
func (args ArgsUrlValues) GetBytesSlice(key string) [][]byte {
	newSlice := [][]byte{}
	if slice, ok := args.Raw[key]; ok {
		for _, item := range slice {
			newSlice = append(newSlice, []byte(item))
		}
	}
	return newSlice
}

func (args ArgsUrlValues) GetString(key string) string { return args.Raw.Get(key) }
func (args ArgsUrlValues) GetStringSlice(key string) []string {
	if slice, ok := args.Raw[key]; ok {
		return slice
	}
	return []string{}
}
