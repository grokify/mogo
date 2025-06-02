package multipartutil

import (
	"bytes"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/grokify/mogo/net/http/httputilmore"
	hum "github.com/grokify/mogo/net/http/httputilmore"
)

// FileInfo represents a file for uploading.
type FileInfo struct {
	MIMEPartName string
	Filepath     string
}

// NewRequest returns a `*http.Request` for making a
// request using multipart/form-data. It supports simple strings
// and files. For more complex field requirements such as JSON
// body parts that require Content-Type headers and Base64
// encoding, use MultipartBuilder directly.
func NewHTTPRequest(method, url string, params url.Values, files []FileInfo) (*http.Request, error) {
	mb := NewMultipartBuilder()
	if err := mb.WriteURLValues(params); err != nil {
		return nil, err
	}
	for _, file := range files {
		if err := mb.WriteFilePathPlus(file.MIMEPartName, file.Filepath, true); err != nil {
			return nil, err
		}
	}
	if err := mb.Close(); err != nil {
		return nil, err
	} else if req, err := http.NewRequest(method, url, mb.Buffer); err != nil {
		return nil, err
	} else if ct, err := mb.ContentTypeHeader(httputilmore.ContentTypeMultipartFormData); err != nil {
		return nil, err
	} else {
		req.Header.Set(httputilmore.HeaderContentType, ct)
		return req, nil
	}
}

func NewReaderFromBodyBytes(body []byte, boundary string) *multipart.Reader {
	return multipart.NewReader(bytes.NewReader(body), boundary)
}

func NewReaderFromHTTPResponse(resp *http.Response) (*multipart.Reader, error) {
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
