package xmlutil

import (
	"encoding/xml"
	"io"
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
	if data, err := io.ReadAll(r); err != nil {
		return err
	} else {
		return xml.Unmarshal(data, v)
	}
}
