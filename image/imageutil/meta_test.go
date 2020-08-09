package imageutil

import (
	"path/filepath"
	"testing"
)

var widthHeightTests = []struct {
	file   string
	width  int
	height int
}{
	{"transparent_640x480.png", 640, 480},
}

// TestImageWidthHeight tests reading image width and height.
func TestImageWidthHeight(t *testing.T) {
	for _, tt := range widthHeightTests {
		img, err := ReadImage(filepath.Join("testdata", tt.file))
		if err != nil {
			t.Errorf("imageutil.ReadImage(\"%s\") Error :[%v]", tt.file, err.Error())
		}
		width, height := RectangleWidthHeight(img.Bounds())
		if width != tt.width {
			t.Errorf("imageutil.RectangleWidthHeight(\"...\") Mismatch: want [%v], got [%v] File [%v]", tt.width, width, tt.file)
		}
		if height != tt.height {
			t.Errorf("imageutil.RectangleWidthHeight(\"...\") Mismatch: want [%v], got [%v] File [%v]", tt.height, height, tt.file)
		}
	}
}
