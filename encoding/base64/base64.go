// base64 supports Base64 encoding and decoding.
package base64

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"regexp"
	"strings"

	"github.com/grokify/mogo/compress/gziputil"
	"github.com/grokify/mogo/encoding"
	"github.com/grokify/mogo/errors/errorsutil"
)

// Decode decodes a byte array to provide an interface like `base64/DecodeString`.
func Decode(enc []byte) ([]byte, error) {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(string(enc))))
	if n, err := base64.StdEncoding.Decode(dst, []byte(string(enc))); err != nil {
		return []byte{}, err
	} else {
		return dst[:n], nil
	}
}

/*
func DecodeOld(input []byte) ([]byte, error) {
	var output []byte
	n, err := base64.StdEncoding.Decode(output, input)
	return output[:n], err
}
*/

const (
	// RxCheckMore is from https://stackoverflow.com/a/8571649/1908967
	RxCheckMore      = `^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)?$`
	RxCheckSimple    = `^[0-9A-Za-z/\+]*=*$`
	RxCheckNoPadding = `^[0-9A-Za-z/\+]*$`
)

var (
	rxCheckMore      = regexp.MustCompile(RxCheckMore)
	rxCheck          = regexp.MustCompile(RxCheckSimple)
	rxCheckNoPadding = regexp.MustCompile(RxCheckNoPadding)
)

// Encode with optional gzip compression. 0 = no compression.
// 9 = best compression.
func EncodeGzip(src []byte, compressLevel int) (string, error) {
	var err error
	if compressLevel != 0 {
		src, err = gziputil.CompressBytes(src, compressLevel)
		if err != nil {
			return "", err
		}
	}
	return base64.StdEncoding.EncodeToString(src), nil
}

// DecodeGunzip base64 decodes a string with optional gzip uncompression.
func DecodeGunzip(enc string) ([]byte, error) {
	enc = Pad(enc)
	if bytes, err := base64.StdEncoding.DecodeString(enc); err != nil {
		return bytes, err
	} else if bytesUnc, err := gziputil.UncompressBytes(bytes); err != nil {
		return bytes, nil
	} else {
		return bytesUnc, nil
	}
}

func IsValid(enc []byte) bool {
	return rxCheckMore.Match(enc)
}

func IsValidString(enc string) bool {
	return rxCheckMore.MatchString(enc)
}

func StripPadding(enc string) string {
	return strings.Replace(enc, "=", "", -1)
}

func Pad(enc string) string {
	return encoding.Pad4(enc, "=")
}

// EncodeGzipJSON encodes a struct that is JSON encoded.
func EncodeGzipJSON(data any, compressLevel int) (string, error) {
	if bytes, err := json.Marshal(data); err != nil {
		return "", err
	} else {
		return EncodeGzip(bytes, compressLevel)
	}
}

// DecodeGunzipJSON base64 decodes a string with optoinal
// gunzip uncompression and then unmarshals the data to a
// struct.
func DecodeGunzipJSON(enc string, output any) error {
	enc = strings.TrimSpace(enc)
	if strings.Index(enc, "{") == 0 || strings.Index(enc, "[") == 0 {
		return json.Unmarshal([]byte(enc), output)
	} else if bytes, err := DecodeGunzip(enc); err != nil {
		return errorsutil.Wrap(err, "DecodeGunzipJSON.DecodeGunzip")
	} else {
		return json.Unmarshal(bytes, output)
	}
}

// ReadAll provides an interface like `io.ReadAll`
// with optional base64 decoding. It is useful for
// decoding `*http.Response.Body`.
func ReadAll(r io.Reader) ([]byte, error) {
	if bytes, err := io.ReadAll(r); err != nil {
		return bytes, err
	} else if IsValid(bytes) {
		return Decode(bytes)
	} else {
		return bytes, err
	}
}
