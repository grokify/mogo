package gziputil

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io"
)

func FixCompressLevel(level int) int {
	if level > 9 {
		level = 9
	} else if level < -1 {
		level = -1
	}
	return level
}

// CompressWriter compresses a byte slide and writes the results
// to the supplied `io.Writer`. When writing to a file, a `*os.File`
// from `os.Create()` can be used as the `io.Writer`.
func CompressWriter(w io.Writer, data []byte, level int) error {
	gw, err := gzip.NewWriterLevel(w, level)
	if err != nil {
		return err
	}
	defer gw.Close()
	_, err = gw.Write(data)
	return err
}

// Compress performs gzip compression on a byte slice.
func Compress(data []byte, level int) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := CompressWriter(buf, data, FixCompressLevel(level))
	return buf.Bytes(), err
}

// CompressBase64 performs gzip compression and then base64 encodes
// the data. Level includes `compress/gzip.BestSpeed`, `compress/gzip.BestCompression`,
// and `compress/gzip.DefaultCompression`.
func CompressBase64(data []byte, level int) (string, error) {
	compressed, err := Compress(data, level)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(compressed), nil
}

// CompressBase64JSON performs a JSON encoding, gzip compression and
// then base64 encodes the data.
func CompressBase64JSON(data any, level int) (string, error) {
	uncompressedBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return CompressBase64(uncompressedBytes, level)
}

// Uncompress gunzips a byte slice.
func Uncompress(compressed []byte) ([]byte, error) {
	gr, err := gzip.NewReader(bytes.NewBuffer(compressed))
	if err != nil {
		return make([]byte, 0), err
	}
	defer gr.Close()
	return io.ReadAll(gr)
}

// UncompressWriter gunzips a byte slice and writes the results
// to a `io.Writer`
func UncompressWriter(w io.Writer, compressed []byte) error {
	uncompressed, err := Uncompress(compressed)
	if err != nil {
		return err
	}
	_, err = w.Write(uncompressed)
	return err
}

// UncompressBase64 base 64 decodes an input string and then
// gunzips the results.
func UncompressBase64(compressedB64 string) ([]byte, error) {
	compressed, err := base64.StdEncoding.DecodeString(compressedB64)
	if err != nil {
		return make([]byte, 0), err
	}
	return Uncompress(compressed)
}

// UncompressBase64JSON JSON encodes data, compresses it and then
// base 64 compresses the data.
func UncompressBase64JSON(compressedB64 string, data any) error {
	uncompressed, err := UncompressBase64(compressedB64)
	if err != nil {
		return err
	}
	return json.Unmarshal(uncompressed, data)
}

// UncompressBase64String  base 64 decodes an input string and then
// gunzips the results, returning a decoded string.
func UncompressBase64String(compressedB64 string) (string, error) {
	byteSlice, err := UncompressBase64(compressedB64)
	return string(byteSlice), err
}
