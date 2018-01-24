package bytesutil

import (
	"bytes"
)

func TrimUTF8BOM(data []byte) []byte {
	return bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))
}
