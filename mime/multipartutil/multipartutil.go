// Package multipartutil provides helper functionality for using multipart.Writer.
// Steps are to call NewMultipartBuilder(), write parts, call builder.Close(), and
// retrieve Content-Type header from builder.Writer.FormDataContentType().

package multipartutil

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"

	hum "github.com/grokify/gotilla/net/httputilmore"
)

func NewReaderBodyBytes(body []byte, boundary string) *multipart.Reader {
	return multipart.NewReader(bytes.NewReader(body), boundary)
}

func NewMultipartReaderForHttpResponse(resp *http.Response) (*multipart.Reader, error) {
	contentType := resp.Header.Get(hum.HeaderContentType)
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, err
	} else if !strings.HasPrefix(mediaType, "multipart/") {
		return nil, fmt.Errorf("MediaType is not multipart [%v]", mediaType)
	}
	if boundary, ok := params["boundary"]; !ok {
		return nil, fmt.Errorf("MIME Boundary not found in Content-Type header [%v]", contentType)
	} else {
		return multipart.NewReader(resp.Body, boundary), nil
	}
}

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

// WriteFieldString adds a text part.
func (builder *MultipartBuilder) WriteFieldString(partName string, data string) error {
	return builder.Writer.WriteField(partName, data)
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

// WriteFilepathPlus adds a file part given a filename with the Content Type and
// other associated headers as needed.
func (builder *MultipartBuilder) WriteFilePathPlus(partName, srcFilepath string, base64Encode bool) error {
	_, filename := filepath.Split(srcFilepath)
	ext := filepath.Ext(srcFilepath)
	mimeType := mime.TypeByExtension(ext)

	header := textproto.MIMEHeader{}
	header.Add(hum.HeaderContentDisposition,
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			partName, filename))
	header.Add(hum.HeaderContentType, mimeType)
	if base64Encode {
		header.Add(hum.HeaderContentTransferEncoding, "base64")
	}

	partWriter, err := builder.Writer.CreatePart(header)
	if err != nil {
		return err
	}

	fileBytes, err := ioutil.ReadFile(srcFilepath)
	if err != nil {
		return err
	}

	if base64Encode {
		str := base64.StdEncoding.EncodeToString(fileBytes)
		_, err = bytes.NewBuffer([]byte(str)).WriteTo(partWriter)
	} else {
		_, err = bytes.NewBuffer(fileBytes).WriteTo(partWriter)

	}
	return err
}

// WriteFilePath adds a file part given a filename.
func (builder *MultipartBuilder) WriteFilePath(partName, srcFilepath string) error {
	file, err := os.Open(srcFilepath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, filename := filepath.Split(srcFilepath)
	return builder.WriteFileReader(partName, filename, file)
}

// WriteFileHeader adds a file part given a part name and *multipart.FileHeader.
// See more at http://sanatgersappa.blogspot.com/2013/03/handling-multiple-file-uploads-in-go.html
// and https://gist.github.com/sanatgersappa/5127317#file-app-go
func (builder *MultipartBuilder) WriteFileHeader(partName string, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	return builder.WriteFileReader(partName, fileHeader.Filename, file)
}

// WriteFileReader adds a file part given a filename and `io.Reader`.
func (builder *MultipartBuilder) WriteFileReader(partName, filename string, src io.Reader) error {
	fw, err := builder.Writer.CreateFormFile(partName, filename)
	if err != nil {
		return err
	}
	_, err = io.Copy(fw, src)
	return err
}

// Close closes the `multipart.Writer`.
func (builder *MultipartBuilder) Close() error {
	return builder.Writer.Close()
}

// ContentType returns the content type for the `Content-Type` header.
func (builder *MultipartBuilder) ContentType() string {
	return builder.Writer.FormDataContentType()
}
