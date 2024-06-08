package gziputil

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io"
)

func UncompressBytes(b []byte) ([]byte, error) {
	return UncompressToBytes(bytes.NewBuffer(b))
}

// UncompressBytes gunzips a byte slice.
func UncompressToBytes(r io.Reader) ([]byte, error) {
	if gr, err := gzip.NewReader(r); err != nil {
		return make([]byte, 0), err
	} else {
		defer gr.Close()
		return io.ReadAll(gr)
	}
}

// Uncompress gunzips a byte slice and writes the results to a `io.Writer`
func Uncompress(w io.Writer, r io.Reader) error {
	if uncompressed, err := UncompressToBytes(r); err != nil {
		return err
	} else {
		_, err = w.Write(uncompressed)
		return err
	}
}

// UncompressBase64String base 64 decodes an input string and then gunzips the results.
// Base64 strings start with `H4sI` to `H4sIAAAAAAAAA`.
func UncompressBase64String(compressedB64 string) ([]byte, error) {
	if compressed, err := base64.StdEncoding.DecodeString(compressedB64); err != nil {
		return make([]byte, 0), err
	} else {
		return UncompressBytes(compressed)
	}
}

// UncompressBase64JSON JSON encodes data, compresses it and then base 64 compresses the data.
func UncompressBase64JSON(compressedB64 string, data any) error {
	if uncompressed, err := UncompressBase64String(compressedB64); err != nil {
		return err
	} else {
		return json.Unmarshal(uncompressed, data)
	}
}

/*
// UncompressBase64String  base 64 decodes an input string and then
// gunzips the results, returning a decoded string.
func UncompressBase64String(compressedB64 string) (string, error) {
	byteSlice, err := UncompressBase64(compressedB64)
	return string(byteSlice), err
}
*/

/*
func GunzipBase64String(s string) ([]byte, error) {
	b, err := base64.Decode([]byte(s))
	if err != nil {
		return b, err
	} else {
		return GunzipBytes(b)
	}
}

func GunzipBytes(b []byte) ([]byte, error) {
	return Gunzip(bytes.NewReader(b))
}

func Gunzip(r io.Reader) ([]byte, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	var uncompressed bytes.Buffer
	if _, err = io.Copy(&uncompressed, gr); err != nil {
		return nil, err
	} else {
		return uncompressed.Bytes(), nil
	}
}
*/
