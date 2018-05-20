package multipartutil

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"

	hum "github.com/grokify/gotilla/net/httputilmore"
)

// MultipartBuilder is a multipart helper.
type MultipartBuilder struct {
	Buffer *bytes.Buffer
	Writer *multipart.Writer
}

// NewMultipartBuilder instantiates a new MultipartBuilder.
func NewMultipartBuilder() MultipartBuilder {
	builder := MultipartBuilder{}
	var b bytes.Buffer
	builder.Buffer = &b
	builder.Writer = multipart.NewWriter(&b)
	return builder
}

// WriteFieldAsJSON adds a JSON part.
func (builder *MultipartBuilder) WriteFieldAsJSON(partName string, data interface{}, base64Encode bool) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	header := textproto.MIMEHeader{}
	header.Add(hum.HeaderContentDisposition, fmt.Sprintf(`form-data; name="%v"`, partName))
	header.Add(hum.HeaderContentType, hum.ContentTypeAppJsonUtf8)
	if base64Encode {
		header.Add(hum.HeaderContentTransferEncoding, "base64")
	}

	partWriter, err := builder.Writer.CreatePart(header)
	if err != nil {
		return err
	}

	if base64Encode {
		str := base64.StdEncoding.EncodeToString(jsonBytes)
		_, err = bytes.NewBuffer([]byte(str)).WriteTo(partWriter)
	} else {
		_, err = bytes.NewBuffer(jsonBytes).WriteTo(partWriter)
	}
	return err
}

// WriteFile adds a file part.
func (builder *MultipartBuilder) WriteFile(partName, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	fw, err := builder.Writer.CreateFormFile(partName, file)
	if err != nil {
		return err
	}
	_, err = io.Copy(fw, f)
	return err
}

func (builder *MultipartBuilder) Close() error {
	return builder.Writer.Close()
}
