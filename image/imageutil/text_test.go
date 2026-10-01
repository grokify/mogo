package imageutil

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
)

func TestDrawText(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 40))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	pt := image.Pt(10, 25)
	DrawText(img, "Hi", pt, color.Black, nil)

	// basicfont.Face7x13 glyphs are 7px wide with an 11px ascent, so "Hi"
	// occupies roughly [10,24) x [14,27).
	textArea := image.Rect(pt.X, pt.Y-11, pt.X+14, pt.Y+2)
	var drawn int
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			if img.RGBAAt(x, y) == (color.RGBA{255, 255, 255, 255}) {
				continue
			}
			if !image.Pt(x, y).In(textArea) {
				t.Fatalf("pixel (%d,%d) changed outside text area %v", x, y, textArea)
			}
			drawn++
		}
	}
	if drawn == 0 {
		t.Fatal("DrawText() drew no pixels")
	}
}
