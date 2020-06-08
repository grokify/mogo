// base64 supports Base64 encoding and decoding.
package base64

import (
	"encoding/base64"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/grokify/gotilla/compress/gziputil"
	"github.com/grokify/gotilla/encoding"
	"github.com/pkg/errors"
)

var (
	rxCheck          = regexp.MustCompile(`^[0-9A-Za-z/\+]*=*$`)
	rxCheckNoPadding = regexp.MustCompile(`^[0-9A-Za-z/\+]*$`)
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
