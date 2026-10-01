package imageutil

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"testing"
	"time"
)

func solid(c color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 8, 6))
	draw.Draw(img, img.Rect, image.NewUniform(c), image.Point{}, draw.Src)
	return img
}

func TestGIFRoundTrip(t *testing.T) {
	colors := []color.RGBA{{255, 0, 0, 255}, {0, 255, 0, 255}, {0, 0, 255, 255}}
	imgs := []image.Image{}
	for _, c := range colors {
		imgs = append(imgs, solid(c))
	}
	g := BuildGIFAnimationSimple(nil, 10, imgs, nil)

	var buf bytes.Buffer
	if err := gif.EncodeAll(&buf, g); err != nil {
		t.Fatal(err)
	}
	dec, err := ReadGIF(&buf)
	if err != nil {
		t.Fatal(err)
	}

	info, err := GIFInfoFromGIF(dec)
	if err != nil {
		t.Fatal(err)
	}
	if info.Width != 8 || info.Height != 6 || info.FrameCount != 3 {
		t.Fatalf("unexpected info: %+v", info)
	}
	if info.Duration != 300*time.Millisecond {
		t.Fatalf("duration: got %v want 300ms", info.Duration)
	}

	frames, err := GIFFrames(dec)
	if err != nil {
		t.Fatal(err)
	}
	if len(frames) != 3 {
		t.Fatalf("frames: got %d want 3", len(frames))
	}
	for i, want := range colors {
		if got := frames[i].RGBAAt(4, 3); got != want {
			t.Errorf("frame %d pixel: got %v want %v", i, got, want)
		}
	}
}

func TestGIFFramesDisposal(t *testing.T) {
	pal := color.Palette{color.Transparent, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}}
	full := image.NewPaletted(image.Rect(0, 0, 4, 4), pal)
	for i := range full.Pix {
		full.Pix[i] = 1 // red
	}
	patch := image.NewPaletted(image.Rect(1, 1, 3, 3), pal)
	for i := range patch.Pix {
		patch.Pix[i] = 2 // blue
	}
	tests := []struct {
		name     string
		disposal byte
		want     color.RGBA // pixel (1,1) in 3rd frame, which draws nothing new there
	}{
		{"none keeps patch", gif.DisposalNone, color.RGBA{0, 0, 255, 255}},
		{"background clears patch", gif.DisposalBackground, color.RGBA{}},
		{"previous restores red", gif.DisposalPrevious, color.RGBA{255, 0, 0, 255}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// frame 2 disposal is under test; frame 3 is a 1px transparent patch elsewhere
			dot := image.NewPaletted(image.Rect(0, 0, 1, 1), pal)
			g := &gif.GIF{
				Image:    []*image.Paletted{full, patch, dot},
				Delay:    []int{1, 1, 1},
				Disposal: []byte{gif.DisposalNone, tt.disposal, gif.DisposalNone},
				Config:   image.Config{Width: 4, Height: 4, ColorModel: pal},
			}
			frames, err := GIFFrames(g)
			if err != nil {
				t.Fatal(err)
			}
			if got := frames[1].RGBAAt(1, 1); got != (color.RGBA{0, 0, 255, 255}) {
				t.Fatalf("frame 1 should show patch, got %v", got)
			}
			if got := frames[2].RGBAAt(1, 1); got != tt.want {
				t.Errorf("frame 2 pixel: got %v want %v", got, tt.want)
			}
		})
	}
}

func TestGIFFramesEmpty(t *testing.T) {
	if _, err := GIFFrames(nil); err == nil {
		t.Error("expected error for nil gif")
	}
	if _, err := GIFInfoFromGIF(&gif.GIF{}); err == nil {
		t.Error("expected error for empty gif")
	}
}

func TestWriteGIFFramesPNG(t *testing.T) {
	g := BuildGIFAnimationSimple(nil, 5, []image.Image{solid(color.White), solid(color.Black)}, nil)
	names, err := WriteGIFFramesPNG(g, t.TempDir(), "f")
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 2 {
		t.Fatalf("got %d files want 2", len(names))
	}
}
