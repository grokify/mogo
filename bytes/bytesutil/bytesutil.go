package bytesutil

import (
	"bytes"
)

const UTF8BOM = "\xef\xbb\xbf"

func TrimUTF8BOM(data []byte) []byte {
	return bytes.TrimPrefix(data, []byte(UTF8BOM))
}
