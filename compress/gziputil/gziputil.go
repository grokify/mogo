package gziputil

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
)

// CompressWriter compresses a byte slide and writes the results
// to the supplied `io.Writer`. When writing to a file, a `*os.File`
// from `os.Create()` can be used as the `io.Writer`.
func CompressWriter(w io.Writer, data []byte) error {
	gw, err := gzip.NewWriterLevel(w, gzip.BestCompression)
	if err != nil {
		return err
	}
	defer gw.Close()
	_, err = gw.Write(data)
	return err
}

// Compress performs gzip compression on a byte slice.
func Compress(data []byte) []byte {
	buf := new(bytes.Buffer)
	CompressWriter(buf, data)
	return buf.Bytes()
}

// CompressBase64 performs gzip compression and then base64 encodes
// the data.
func CompressBase64(uncompressed []byte) string {
	compressed := Compress(uncompressed)
	return base64.StdEncoding.EncodeToString(compressed)
}

// CompressBase64JSON performs a JSON encoding, gzip compression and
// then base64 encodes the data.
func CompressBase64JSON(uncompressedData interface{}) (string, error) {
	uncompressedBytes, err := json.Marshal(uncompressedData)
	if err != nil {
		return "", err
	}
	return CompressBase64(uncompressedBytes), nil
}

// Uncompress gunzips a byte slice.
func Uncompress(compressed []byte) ([]byte, error) {
	gr, err := gzip.NewReader(bytes.NewBuffer(compressed))
	if err != nil {
		var byteSlice []byte
		return byteSlice, err
	}
	defer gr.Close()
	return ioutil.ReadAll(gr)
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
	var byteSlice []byte
	compressed, err := base64.StdEncoding.DecodeString(compressedB64)
	if err != nil {
		return byteSlice, err
	}
	return Uncompress(compressed)
}

// UncompressBase64JSON JSON encodes data, compresses it and then
// base 64 compresses the data.
func UncompressBase64JSON(compressedB64 string, data interface{}) error {
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
