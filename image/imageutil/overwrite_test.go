package imageutil

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// TestWriteFileOverwriteTruncates writes a large image and then a small one
// to the same path, and checks that no bytes from the large image remain.
func TestWriteFileOverwriteTruncates(t *testing.T) {
	noisy := func(w, h int) image.Image {
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		for y := range h {
			for x := range w {
				img.Set(x, y, color.RGBA{uint8(x * 7 % 256), uint8(y * 13 % 256), uint8((x + y) % 256), 255})
			}
		}
		return img
	}
	big, small := noisy(300, 300), noisy(10, 10)

	tests := []struct {
		name   string
		write  func(path string, img image.Image) error
		encode func(img image.Image) ([]byte, error)
	}{
		{
			"jpeg",
			func(p string, img image.Image) error { return Image{Image: img}.WriteJPEGFile(p, nil) },
			func(img image.Image) ([]byte, error) {
				var b bytes.Buffer
				err := Image{Image: img}.WriteJPEG(&b, nil)
				return b.Bytes(), err
			},
		},
		{
			"png",
			func(p string, img image.Image) error { return Image{Image: img}.WritePNGFile(p) },
			func(img image.Image) ([]byte, error) {
				var b bytes.Buffer
				err := png.Encode(&b, img)
				return b.Bytes(), err
			},
		},
		{
			"gif",
			func(p string, img image.Image) error {
				g, err := Image{Image: img}.GIF()
				if err != nil {
					return err
				}
				return WriteGIFFile(p, g, 0600)
			},
			func(img image.Image) ([]byte, error) {
				g, err := Image{Image: img}.GIF()
				if err != nil {
					return nil, err
				}
				var b bytes.Buffer
				err = gif.EncodeAll(&b, g)
				return b.Bytes(), err
			},
		},
	}
	for _, tt := range tests {
		p := filepath.Join(t.TempDir(), "out."+tt.name)
		if err := tt.write(p, big); err != nil {
			t.Fatalf("%s: write big: %v", tt.name, err)
		}
		if err := tt.write(p, small); err != nil {
			t.Fatalf("%s: write small: %v", tt.name, err)
		}
		got, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("%s: %v", tt.name, err)
		}
		want, err := tt.encode(small)
		if err != nil {
			t.Fatalf("%s: %v", tt.name, err)
		}
		if !bytes.Equal(got, want) {
			t.Errorf("%s: file is %d bytes, want %d (old bytes left behind)", tt.name, len(got), len(want))
		}
	}
}
