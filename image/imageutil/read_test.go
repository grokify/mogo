package imageutil

import (
	"testing"
)

var readImageFileTests = []struct {
	filename   string
	formatName string
}{
	{"read_testdata/gopher_color.png", FormatNamePNG},
	{"read_testdata/gopher_color.jpg", FormatNameJPG},
	{"read_testdata/gopher_color.webp", FormatNameWEBP},
}

func TestReadImageFile(t *testing.T) {
	for _, tt := range readImageFileTests {
		_, formatName, err := ReadImageFile(tt.filename)
		if err != nil {
			t.Errorf("imageutil.ReadImageFile(\"%s\") Error: [%v]", tt.filename, err.Error())
		}
		if formatName != tt.formatName {
			t.Errorf("imageutil.ReadImageFile(\"%s\") Format Want [%s], Got [%v]", tt.filename, tt.formatName, formatName)
		}
	}
}
