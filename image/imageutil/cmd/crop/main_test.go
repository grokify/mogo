package main

import (
	"image"
	"image/color"
	"testing"
)

// quadrants returns a w x h image whose left half is red and right half is
// blue, so tests can tell which part of the source an output kept.
func quadrants(w, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			c := color.RGBA{R: 255, A: 255}
			if x >= w/2 {
				c = color.RGBA{B: 255, A: 255}
			}
			img.Set(x, y, c)
		}
	}
	return img
}

func defaults() options {
	return options{alignX: "center", alignY: "center"}
}

func size(img image.Image) image.Point {
	return image.Pt(img.Bounds().Dx(), img.Bounds().Dy())
}

// colorAt returns the color at (x, y) relative to the image's bounds.
func colorAt(img image.Image, x, y int) color.RGBA {
	c, ok := color.RGBAModel.Convert(img.At(img.Bounds().Min.X+x, img.Bounds().Min.Y+y)).(color.RGBA)
	if !ok {
		panic("color.RGBAModel did not return color.RGBA")
	}
	return c
}

var (
	red  = color.RGBA{R: 255, A: 255}
	blue = color.RGBA{B: 255, A: 255}
)

func TestProcessCrop(t *testing.T) {
	o := defaults()
	o.width, o.alignX = 40, "left"
	outs, err := process(quadrants(100, 50), o)
	if err != nil {
		t.Fatal(err)
	}
	if got := size(outs[0].img); got != image.Pt(40, 50) {
		t.Errorf("size = %v, want (40,50)", got)
	}
	if got := colorAt(outs[0].img, 39, 0); got != red {
		t.Errorf("left-aligned crop kept %v at right edge, want red", got)
	}

	o = defaults()
	o.height, o.alignY = 20, "bottom"
	outs, err = process(quadrants(100, 50), o)
	if err != nil {
		t.Fatal(err)
	}
	if got := size(outs[0].img); got != image.Pt(100, 20) {
		t.Errorf("size = %v, want (100,20)", got)
	}
}

func TestProcessSquare(t *testing.T) {
	o := defaults()
	o.square = "smaller"
	outs, err := process(quadrants(100, 50), o)
	if err != nil {
		t.Fatal(err)
	}
	if got := size(outs[0].img); got != image.Pt(50, 50) {
		t.Errorf("smaller: size = %v, want (50,50)", got)
	}

	o.square, o.bg = "larger", color.RGBA{G: 255, A: 255}
	outs, err = process(quadrants(100, 50), o)
	if err != nil {
		t.Fatal(err)
	}
	if got := size(outs[0].img); got != image.Pt(100, 100) {
		t.Errorf("larger: size = %v, want (100,100)", got)
	}
	if got := colorAt(outs[0].img, 0, 0); got != o.bg {
		t.Errorf("larger: padding = %v, want %v", got, o.bg)
	}
}

func TestProcessSplit(t *testing.T) {
	o := defaults()
	o.split = "x"
	outs, err := process(quadrants(100, 50), o)
	if err != nil {
		t.Fatal(err)
	}
	if len(outs) != 2 || outs[0].suffix != "_left" || outs[1].suffix != "_right" {
		t.Fatalf("split x outputs = %+v, want _left and _right", outs)
	}
	for i, want := range []color.RGBA{red, blue} {
		if got := size(outs[i].img); got != image.Pt(50, 50) {
			t.Errorf("%s size = %v, want (50,50)", outs[i].suffix, got)
		}
		if got := colorAt(outs[i].img, 25, 25); got != want {
			t.Errorf("%s color = %v, want %v", outs[i].suffix, got, want)
		}
	}

	o.split = "y"
	outs, err = process(quadrants(100, 50), o)
	if err != nil {
		t.Fatal(err)
	}
	if len(outs) != 2 || outs[0].suffix != "_top" || size(outs[1].img) != image.Pt(100, 25) {
		t.Errorf("split y outputs = %+v, want _top/_bottom of 100x25", outs)
	}
}

func TestProcessFlattensTransparency(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 10, 10)) // fully transparent
	o := defaults()
	o.bg = color.RGBA{G: 255, A: 255}
	outs, err := process(src, o)
	if err != nil {
		t.Fatal(err)
	}
	if got := colorAt(outs[0].img, 5, 5); got != o.bg {
		t.Errorf("transparent pixel = %v, want background %v", got, o.bg)
	}

	o.bg = nil // defaults to white
	if outs, err = process(src, o); err != nil {
		t.Fatal(err)
	}
	if got, want := colorAt(outs[0].img, 5, 5), (color.RGBA{255, 255, 255, 255}); got != want {
		t.Errorf("transparent pixel = %v, want white", got)
	}
}

func TestProcessInvalid(t *testing.T) {
	for name, mod := range map[string]func(*options){
		"negative width": func(o *options) { o.width = -1 },
		"bad align-x":    func(o *options) { o.alignX = "top" },
		"bad align-y":    func(o *options) { o.alignY = "left" },
		"bad square":     func(o *options) { o.square = "round" },
		"bad split":      func(o *options) { o.split = "z" },
	} {
		o := defaults()
		mod(&o)
		if _, err := process(quadrants(10, 10), o); err == nil {
			t.Errorf("%s: process() error = nil, want error", name)
		}
	}
}

func TestOutputName(t *testing.T) {
	for _, tt := range []struct{ name, suffix, want string }{
		{"out.jpg", "", "out.jpg"},
		{"out.jpg", "_left", "out_left.jpg"},
		{"dir/scan.jpeg", "_top", "dir/scan_top.jpeg"},
		{"noext", "_right", "noext_right"},
	} {
		if got := outputName(tt.name, tt.suffix); got != tt.want {
			t.Errorf("outputName(%q, %q) = %q, want %q", tt.name, tt.suffix, got, tt.want)
		}
	}
}
