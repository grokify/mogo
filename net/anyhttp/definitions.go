package anyhttp

import (
	"io"
	"mime/multipart"

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
