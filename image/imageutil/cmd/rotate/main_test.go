package main

import (
	"image"
	"image/color"
	"testing"

	"golang.org/x/image/draw"
)

var (
	red   = color.RGBA{R: 255, A: 255}
	blue  = color.RGBA{B: 255, A: 255}
	black = color.RGBA{A: 255}
)

// halves returns a w x h image whose left half is red and right half is blue.
func halves(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(img, image.Rect(0, 0, w/2, h), image.NewUniform(red), image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(w/2, 0, w, h), image.NewUniform(blue), image.Point{}, draw.Src)
	return img
}

func colorAt(img image.Image, x, y int) color.RGBA {
	c, ok := color.RGBAModel.Convert(img.At(img.Bounds().Min.X+x, img.Bounds().Min.Y+y)).(color.RGBA)
	if !ok {
		panic("color.RGBAModel did not return color.RGBA")
	}
	return c
}

func size(img image.Image) image.Point { return img.Bounds().Size() }

func TestProcessQuarterTurns(t *testing.T) {
	tests := []struct {
		degrees float64
		// color of the top-left pixel after rotating
		topLeft color.RGBA
	}{
		{90, red},   // left half moves to the top
		{-90, blue}, // right half moves to the top
		{180, blue},
	}
	for _, tt := range tests {
		got := process(halves(100, 40), options{degrees: tt.degrees, bg: black, opaque: true})
		want := image.Pt(40, 100)
		if tt.degrees == 180 {
			want = image.Pt(100, 40)
		}
		if size(got) != want {
			t.Errorf("%v°: size = %v, want %v", tt.degrees, size(got), want)
		}
		if c := colorAt(got, 0, 0); c != tt.topLeft {
			t.Errorf("%v°: top-left = %v, want %v", tt.degrees, c, tt.topLeft)
		}
	}
}

func TestProcessExpandFillsBackground(t *testing.T) {
	got := process(halves(100, 100), options{degrees: 30, bg: black, interp: draw.BiLinear, opaque: true})
	if size(got).X <= 100 || size(got).Y <= 100 {
		t.Errorf("size = %v, want larger than 100x100", size(got))
	}
	if c := colorAt(got, 0, 0); c != black {
		t.Errorf("corner = %v, want %v", c, black)
	}
}

func TestProcessCropHasNoBackground(t *testing.T) {
	got := process(halves(100, 100), options{degrees: 30, crop: true, bg: black, interp: draw.BiLinear, opaque: true})
	if size(got).X >= 100 || size(got).Y >= 100 {
		t.Errorf("size = %v, want smaller than 100x100", size(got))
	}
	b := got.Bounds()
	for _, p := range []image.Point{{0, 0}, {b.Dx() - 1, 0}, {0, b.Dy() - 1}, {b.Dx() - 1, b.Dy() - 1}} {
		if c := colorAt(got, p.X, p.Y); c == black {
			t.Errorf("corner %v shows background", p)
		}
	}
}

func TestProcessFlattensTransparencyForJPEG(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 10, 10)) // fully transparent
	got := process(src, options{degrees: 90, bg: black, opaque: true})
	if c := colorAt(got, 5, 5); c != black {
		t.Errorf("opaque output = %v, want %v", c, black)
	}
	got = process(src, options{degrees: 90, bg: black, opaque: false})
	if c := colorAt(got, 5, 5); c.A != 0 {
		t.Errorf("PNG output = %v, want transparent", c)
	}
}

func TestParseInterpolator(t *testing.T) {
	for _, alg := range []string{"nearestneighbor", "approxbilinear", "bilinear", "catmullrom"} {
		if _, err := parseInterpolator(alg); err != nil {
			t.Errorf("%s: %v", alg, err)
		}
	}
	if _, err := parseInterpolator("bogus"); err == nil {
		t.Error("bogus: want error")
	}
}
