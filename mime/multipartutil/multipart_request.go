package multipartutil

import (
	"net/http"
	"net/url"

	"github.com/grokify/mogo/net/http/httputilmore"
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
func NewRequest(method, url string, params url.Values, files []FileInfo) (*http.Request, error) {
	mb := NewMultipartBuilder()
	err := mb.WriteURLValues(params)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		err := mb.WriteFilePathPlus(file.MIMEPartName, file.Filepath, true)
		if err != nil {
			return nil, err
		}
	}
	if err = mb.Close(); err != nil {
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
