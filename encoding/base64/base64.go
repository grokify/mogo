// base64 supports Base64 encoding and decoding.
package base64

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/grokify/simplego/compress/gziputil"
	"github.com/grokify/simplego/encoding"
	"github.com/pkg/errors"
)

// Decode decodes a byte array to provide an interface
// like `base64/DecodeString`.
func Decode(input []byte) ([]byte, error) {
	var output []byte
	n, err := base64.StdEncoding.Decode(output, input)
	return output[:n], err
}

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
func EncodeGzip(data []byte, compressLevel int) string {
	if compressLevel != 0 {
		data = gziputil.Compress(data, compressLevel)
	}
	return base64.StdEncoding.EncodeToString(data)
}

// DecodeGunzip base64 decodes a string with optional
// gzip uncompression.
func DecodeGunzip(encoded string) ([]byte, error) {
	encoded = Pad(encoded)
	bytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return bytes, err
	}
	bytesUnc, err := gziputil.Uncompress(bytes)
	if err != nil {
		return bytes, nil
	}
	return bytesUnc, nil
}

func IsValid(input []byte) bool {
	return rxCheckMore.Match(input)
}

func IsValidString(input string) bool {
	return rxCheckMore.MatchString(input)
}

func StripPadding(str string) string {
	return strings.Replace(str, "=", "", -1)
}

func Pad(encoded string) string {
	return encoding.Pad4(encoded, "=")
}

// EncodeGzipJSON encodes a struct that is JSON encoded.
func EncodeGzipJSON(data interface{}, compressLevel int) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return EncodeGzip(bytes, compressLevel), err
}

// DecodeGunzipJSON base64 decodes a string with optoinal
// gunzip uncompression and then unmarshals the data to a
// struct.
func DecodeGunzipJSON(encoded string, output interface{}) error {
	encoded = strings.TrimSpace(encoded)
	if strings.Index(encoded, "{") == 0 || strings.Index(encoded, "[") == 0 {
		return json.Unmarshal([]byte(encoded), output)
	}
	bytes, err := DecodeGunzip(encoded)
	if err != nil {
		return errors.Wrap(err, "DecodeGunzipJSON.DecodeGunzip")
	}
	return json.Unmarshal(bytes, output)
}

// ReadAll provides an interface like `ioutil.ReadAll`
// with optional base64 decoding. It is useful for
// decoding `*http.Response.Body`.
func ReadAll(r io.Reader) ([]byte, error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return bytes, err
	}
	if IsValid(bytes) {
		return Decode(bytes)
	}
	return bytes, err
}
