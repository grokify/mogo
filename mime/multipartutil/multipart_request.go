package multipartutil

import (
	"net/http"
	"net/url"

	"github.com/grokify/simplego/net/httputilmore"
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
	err = mb.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, mb.Buffer)
	if err != nil {
		return nil, err
	}
	req.Header.Set(httputilmore.HeaderContentType, mb.ContentType())
	return req, nil
}
