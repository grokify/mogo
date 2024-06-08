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

/*
// Uncompress gunzips a byte slice.
func Uncompress(compressed []byte) ([]byte, error) {
	gr, err := gzip.NewReader(bytes.NewBuffer(compressed))
 := UncompressBase64(compressedB64)
	return string(byteSlice), err
}
*/
