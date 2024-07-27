package gziputil

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
)

func FixCompressLevel(level int) int {
	if level > gzip.BestCompression {
		level = gzip.BestCompression
	} else if level < gzip.DefaultCompression {
		level = gzip.DefaultCompression
	}
	return level
}

// Compress compresses a byte slide and writes the results
// to the supplied `io.Writer`. When writing to a file, a `*os.File`
// from `os.Create()` can be used as the `io.Writer`.
func Compress(w io.Writer, r io.Reader, level int) error {
	gw, err := gzip.NewWriterLevel(w, FixCompressLevel(level))
	if err != nil {
		return err
	}
	defer gw.Close()
	if bytes, err := io.ReadAll(r); err != nil {
		return err
	} else {
		_, err = gw.Write(bytes)
		return err
	}
}

// CompressBytes performs gzip compression on a byte slice.
func CompressBytes(data []byte, level int) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := Compress(buf, bytes.NewReader(data), level)
	return buf.Bytes(), err
}

// CompressToBase64String performs gzip compression and then base64 encodes
// the data. Level includes `compress/gzip.BestSpeed`, `compress/gzip.BestCompression`,
// and `compress/gzip.DefaultCompression`.
func CompressToBase64String(data []byte, level int) (string, error) {
	if compressed, err := CompressBytes(data, level); err != nil {
		return "", err
	} else {
		return base64.StdEncoding.EncodeToString(compressed), nil
	}
}

// CompressFileToBase64String compresses a file to base64 encoded string.
func CompressFileToBase64String(name string, level int) (string, error) {
	if b, err := os.ReadFile(name); err != nil {
		return "", err
	} else {
		return CompressToBase64String(b, level)
	}
}

// CompressBase64JSON performs a JSON encoding, gzip compression and
// then base64 encodes the data.
func CompressBase64JSON(data any, level int) (string, error) {
	uncompressedBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return CompressToBase64String(uncompressedBytes, level)
}
