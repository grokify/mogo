package mimeutil

import (
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/grokify/mogo/net/http/httputilmore"
)

const DefaultMIMEType = httputilmore.ContentTypeAppOctetStream

var ErrUnknownMediaType = errors.New("unknown media type")

// MustTypeByFilename follows the convention of
// `mime.TypeByExtension` by returning an empty
// string if type not found. If `useDefault` is
// set, a non-detected value is set to `application/octet-stream`
// which is the default for `http.DetectContentType`.
func MustTypeByFilename(nameOrExt string, useDefault bool) string {
	mt, err := TypeByFilename(nameOrExt)
	if err != nil || len(strings.TrimSpace(mt)) == 0 {
		if useDefault {
			return DefaultMIMEType
		} else {
			return ""
		}
	}
	return mt
}

// TypeByFilename detects a mimetype using `mime.TypeByExtension`.
func TypeByFilename(nameOrExt string) (string, error) {
	ext := strings.TrimSpace(nameOrExt)
	if strings.Contains(ext, ".") {
		m := regexp.MustCompile(`(.[^.]+)$`).FindStringSubmatch(ext)
		if len(m) < 2 {
			return "", errors.New("extension not found")
		}
		ext = m[1]
	} else {
		ext = "." + ext
	}
	mt := mime.TypeByExtension(ext)
	if len(mt) == 0 {
		return "", errors.New("type not detected")
	}
	return mt, nil
}

// MustTypeByFile follows the convention of
// `mime.TypeByExtension` by returning an empty
// string if type not found. If `useDefault` is
// set, a non-detected value is set to `application/octet-stream`
// which is the default for `http.DetectContentType`.
func MustTypeByFile(name string, useDefault bool) string {
	mt, err := TypeByFile(name)
	if err != nil || len(strings.TrimSpace(mt)) == 0 {
		if useDefault {
			return DefaultMIMEType
		} else {
			return ""
		}
	}
	return mt
}

// TypeByFile detects the media type by reading MIME type
// information of the file content. It relies on
// `http.DetectContentType“
func TypeByFile(name string) (string, error) {
	file, err := os.Open(name)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return TypeByReadSeeker(file, false)
}

// TypeByReadSeeker detects the media type by reading MIME type
// information of an `io.ReadSeeker`. It relies on
// `http.DetectContentType`
func TypeByReadSeeker(rs io.ReadSeeker, resetPointer bool) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	data := make([]byte, 512)
	_, err := rs.Read(data)
	if err != nil {
		return "", err
	}

	if resetPointer {
		// Reset the read pointer if necessary.
		_, err = rs.Seek(0, 0)
		if err != nil {
			return "", err
		}
	}

	// Returns `application/octet-stream` if media type is unknown.
	return http.DetectContentType(data), nil
}

// IsType checks to see if a media type corresponds to a type/subtype
func IsType(s, mediaType string) bool {
	if s == mediaType {
		return true
	}
	s = strings.ToLower(strings.TrimSpace(s))
	if s == mediaType {
		return true
	}
	if !strings.Contains(mediaType, ";") {
		if strings.Index(s, mediaType+";") == 0 {
			return true
		}
		return regexp.MustCompile(`^(?i)` + regexp.QuoteMeta(mediaType) + `\s*;`).MatchString(s)
	}
	return false
}
