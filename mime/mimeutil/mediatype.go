package mimeutil

import (
	"io"
	"net/http"
	"os"
)

// DetectContentTypeFile detects the media type by reading MIME type information of the file content.
func DetectContentTypeFile(name string) (string, error) {
	file, err := os.Open(name)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return DetectContentTypeReadSeeker(file)
}

// DetectContentTypeReadSeeker detects the media type by reading MIME type information of an `io.ReadSeeker`.
func DetectContentTypeReadSeeker(rs io.ReadSeeker) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	data := make([]byte, 512)
	_, err := rs.Read(data)
	if err != nil {
		return "", err
	}

	// Reset the read pointer if necessary.
	rs.Seek(0, 0)

	// Returns `application/octet-stream` if media type is unknown.
	return http.DetectContentType(data), nil
}
