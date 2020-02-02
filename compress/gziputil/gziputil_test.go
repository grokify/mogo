package gziputil

import (
	"testing"
)

var compressTests = []struct {
	uncompressed string
	compressed   string
}{
	{"input", `H4sIAAAAAAAE/wAFAPr/aW5wdXQBAAD//9cyKNgFAAAA`},
	{"input input input ", `H4sIAAAAAAAE/wASAO3/aW5wdXQgaW5wdXQgaW5wdXQgAQAA//8NX1C8EgAAAA==`},
}

func TestGzipCompress(t *testing.T) {
	for _, tt := range compressTests {
		gotCompress := CompressBase64([]byte(tt.uncompressed))

		if 1 == 0 {
			if gotCompress != tt.compressed {
				t.Errorf("gziputil.CompressBase64(\"%s\") Mismatch: want [%v] got [%v]",
					tt.uncompressed, tt.compressed, gotCompress)
			}

			gotUncompressString, err := UncompressBase64String(tt.compressed)
			if err != nil {
				t.Errorf("gziputil.UncompressBase64(\"%s\") Error [%v]",
					tt.compressed, err.Error())
			}
			if gotUncompressString != tt.uncompressed {
				t.Errorf("gziputil.UncompressBase64(\"%s\") Mismatch: want [%v] got [%v]",
					tt.compressed, tt.uncompressed, gotUncompressString)
			}
		}
		gotUncompressString, err := UncompressBase64String(tt.compressed)
		if err != nil {
			t.Errorf("gziputil.UncompressBase64(\"%s\") Error [%v]",
				tt.compressed, err.Error())
		}
		if gotUncompressString != tt.uncompressed {
			t.Errorf("gziputil.UncompressBase64(\"%s\") Mismatch: want [%v] got [%v]",
				tt.compressed, tt.uncompressed, gotUncompressString)
		}

	}
}
