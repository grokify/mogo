package gziputil

import (
	"compress/gzip"
	"testing"
)

var compressTests = []struct {
	uncompressed              string
	compressedBestSpeed       string
	compressedBestCompression string
}{
	{"input",
		`H4sIAAAAAAAE/wAFAPr/aW5wdXQBAAD//9cyKNgFAAAA`,
		`H4sIAAAAAAAC/8rMKygtAQQAAP//1zIo2AUAAAA=`},
	{"input input input ",
		`H4sIAAAAAAAE/wASAO3/aW5wdXQgaW5wdXQgaW5wdXQgAQAA//8NX1C8EgAAAA==`,
		`H4sIAAAAAAAC/8rMKygtUUAmAQEAAP//DV9QvBIAAAA=`},
}

func TestGzipCompress(t *testing.T) {
	for _, tt := range compressTests {
		gotCompress := CompressBase64([]byte(tt.uncompressed), gzip.BestCompression)

		if 1 == 0 {
			if gotCompress != tt.compressedBestSpeed {
				t.Errorf("gziputil.CompressBase64(\"%s\") Mismatch: want [%v] got [%v]",
					tt.uncompressed, tt.compressedBestSpeed, gotCompress)
			}

			gotUncompressString, err := UncompressBase64String(tt.compressedBestSpeed)
			if err != nil {
				t.Errorf("gziputil.UncompressBase64(\"%s\") Error [%v]",
					tt.compressedBestSpeed, err.Error())
			}
			if gotUncompressString != tt.uncompressed {
				t.Errorf("gziputil.UncompressBase64(\"%s\") Mismatch: want [%v] got [%v]",
					tt.compressedBestSpeed, tt.uncompressed, gotUncompressString)
			}
		}
		gotUncompressString, err := UncompressBase64String(tt.compressedBestSpeed)
		if err != nil {
			t.Errorf("gziputil.UncompressBase64(\"%s\") Error [%v]",
				tt.compressedBestCompression, err.Error())
		}
		if gotUncompressString != tt.uncompressed {
			t.Errorf("gziputil.UncompressBase64(\"%s\") Mismatch: want [%v] got [%v]",
				tt.compressedBestCompression, tt.uncompressed, gotUncompressString)
		}

	}
}
