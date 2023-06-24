package lzfutil

import (
	"os"

	lzf "github.com/zhuyie/golzf"
)

// ReadFile reads a LZF compressed file, e.g. one compressed by PHP `lzf_compress`.
func ReadFile(file string) ([]byte, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return []byte{}, err
	}
	// *100 is an safe guess to ensure the file can be decoded.
	decompressed := make([]byte, len(bytes)*100)
	_, err = lzf.Decompress(bytes, decompressed)
	return decompressed, err
}
