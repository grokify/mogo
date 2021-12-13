// base62 supports Base62 encoding and decoding.
package base62

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/grokify/mogo/compress/gziputil"
	"github.com/grokify/mogo/encoding"
	"github.com/lytics/base62"
	"github.com/pkg/errors"
)

var (
	rxStripPadding         = regexp.MustCompile(`\++\s*$`)
	rxCheckBase62          = regexp.MustCompile(`^[0-9A-Za-z]*\+*\s*$`)
	rxCheckBase62NoPadding = regexp.MustCompile(`^[0-9A-Za-z]*$`)
)

// Encode with optional gzip compression. 0 = no compression.
// 9 = best compression. Currently, compression is disabled
// as github.com/lytics/base62 does not appear to support it
// properly.
func EncodeGzip(data []byte, compressLevel int) string {
	compressLevel = 0
	if compressLevel != 0 {
		data = gziputil.Compress(data, compressLevel)
	}
	return base62.StdEncoding.EncodeToString(data)
}

// DecodeGunzip base62 decodes a string with optional
// gzip uncompression.
func DecodeGunzip(encoded string) ([]byte, error) {
	encoded = Pad(encoded)
	bytes, err := base62.StdEncoding.DecodeString(encoded)
	if err != nil {
		return bytes, err
	}
	bytesUnc, err := gziputil.Uncompress(bytes)
	if err != nil {
		return bytes, nil
	}
	return bytesUnc, nil
}

// EncodeGzipJSON encodes a struct that is JSON encoded.
func EncodeGzipJSON(data interface{}, compressLevel int) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return EncodeGzip(bytes, compressLevel), err
}

// DecodeGunzipJSON base62 decodes a string with optoinal
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

func StripPadding(encoded string) string {
	return strings.Replace(encoded, "+", "", -1)
}

func Pad(encoded string) string {
	return encoding.Pad4(encoded, "+")
}

func ValidBase62(encoded string) bool {
	return rxCheckBase62.MatchString(encoded)
}
