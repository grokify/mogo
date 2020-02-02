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

type compressTestObject struct {
	Uncompressed              string
	CompressedBestCompression string
	CompressedBestSpeed       string
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

		dto := compressTestObject{
			Uncompressed:              tt.uncompressed,
			CompressedBestCompression: tt.compressedBestCompression,
			CompressedBestSpeed:       tt.compressedBestSpeed}
		dtoBase64, err := CompressBase64JSON(dto, gzip.BestCompression)
		if err != nil {
			t.Errorf("gziputil.CompressBase64JSON(dto, %v) Error [%v]",
				gzip.BestCompression, err.Error())
		}
		dto2 := compressTestObject{}
		err = UncompressBase64JSON(dtoBase64, &dto2)
		if err != nil {
			t.Errorf("gziputil.UncompressBase64JSON(\"%s\", %v) Error [%v]",
				dtoBase64, gzip.BestCompression, err.Error())
		}
		if dto2.Uncompressed != tt.uncompressed {
			t.Errorf("gziputil.UncompressBase64JSON(\"%s\",%v) Mismatch: want [%v] got [%v]",
				dtoBase64, tt.compressedBestCompression, tt.uncompressed, dto2.Uncompressed)
		}
	}
}
