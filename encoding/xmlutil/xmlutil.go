package xmlutil

import (
	"encoding/xml"
	"io"

	"github.com/grokify/mogo/io/ioutilmore"
)

func MarshalIndent(v any, prefix, indent string, addDoctype bool) ([]byte, error) {
	data, err := xml.MarshalIndent(v, prefix, indent)
	if err != nil || !addDoctype {
		return data, err
	}
	out := []byte(xml.Header)
	return append(out, data...), nil
}

func UnmarshalReader(r io.Reader, v any) error {
	data, err := ioutilmore.ReaderToBytes(r)
	if err != nil {
		return err
	}
	return xml.Unmarshal(data, v)
}
