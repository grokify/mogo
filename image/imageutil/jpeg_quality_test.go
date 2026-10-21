package imageutil

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"testing"
)

// TestWriteJPEGNilOptionsUsesDefaultQuality checks that nil options encode
// at jpeg.DefaultQuality, not the lowest quality.
func TestWriteJPEGNilOptionsUsesDefaultQuality(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := range 64 {
		for x := range 64 {
			img.Set(x, y, color.RGBA{uint8(x * 4 % 256), uint8(y * 4 % 256), 128, 255})
		}
	}
	var want bytes.Buffer
	if err := jpeg.Encode(&want, img, &jpeg.Options{Quality: jpeg.DefaultQuality}); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		opt  *JPEGEncodeOptions
	}{
		{"nil options", nil},
		{"nil jpeg.Options", &JPEGEncodeOptions{}},
	}
	for _, tt := range tests {
		got, err := Image{Image: img}.JPEGBytes(tt.opt)
		if err != nil {
			t.Fatalf("%s: %v", tt.name, err)
		}
		if !bytes.Equal(got, want.Bytes()) {
			t.Errorf("%s: output (%d bytes) does not match quality %d (%d bytes)",
				tt.name, len(got), jpeg.DefaultQuality, want.Len())
		}
	}
}
