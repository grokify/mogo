package imageutil

import (
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/image/draw"
)

var (
	red   = color.RGBA{255, 0, 0, 255}
	green = color.RGBA{0, 255, 0, 255}
	blue  = color.RGBA{0, 0, 255, 255}
	white = color.RGBA{255, 255, 255, 255}
)

// newTestGrid returns a 3x2 image where each pixel has a unique color whose
// R and G channels encode its x and y, so positions can be tracked.
func newTestGrid() *image.RGBA {
	levels := []uint8{0, 50, 100}
	img := image.NewRGBA(image.Rect(0, 0, 3, 2))
	for y := range 2 {
		for x := range 3 {
			img.SetRGBA(x, y, color.RGBA{levels[x], levels[y], 100, 255})
		}
	}
	return img
}

// newHalves returns a w x h image whose left half is left and right half is right.
func newHalves(w, h int, left, right color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(img, image.Rect(0, 0, w/2, h), image.NewUniform(left), image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(w/2, 0, w, h), image.NewUniform(right), image.Point{}, draw.Src)
	return img
}

func newSolid(w, h int, c color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(img, img.Bounds(), image.NewUniform(c), image.Point{}, draw.Src)
	return img
}

func assertSameImage(t *testing.T, name string, got, want *image.RGBA) {
	t.Helper()
	if got.Bounds() != want.Bounds() {
		t.Fatalf("%s: bounds mismatch: got %v, want %v", name, got.Bounds(), want.Bounds())
	}
	b := want.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if got.RGBAAt(x, y) != want.RGBAAt(x, y) {
				t.Fatalf("%s: pixel (%d,%d): got %v, want %v", name, x, y, got.RGBAAt(x, y), want.RGBAAt(x, y))
			}
		}
	}
}

func TestRotateQuarterTurns(t *testing.T) {
	src := newTestGrid()
	at := func(x, y int) color.RGBA { return src.RGBAAt(x, y) }

	tests := []struct {
		name string
		got  *image.RGBA
		// want[y][x] in output coordinates
		want [][]color.RGBA
	}{
		{"Rotate90", Rotate90(src), [][]color.RGBA{
			{at(0, 1), at(0, 0)},
			{at(1, 1), at(1, 0)},
			{at(2, 1), at(2, 0)},
		}},
		{"Rotate180", Rotate180(src), [][]color.RGBA{
			{at(2, 1), at(1, 1), at(0, 1)},
			{at(2, 0), at(1, 0), at(0, 0)},
		}},
		{"Rotate270", Rotate270(src), [][]color.RGBA{
			{at(2, 0), at(2, 1)},
			{at(1, 0), at(1, 1)},
			{at(0, 0), at(0, 1)},
		}},
	}
	for _, tt := range tests {
		wantH, wantW := len(tt.want), len(tt.want[0])
		if tt.got.Bounds() != image.Rect(0, 0, wantW, wantH) {
			t.Fatalf("%s: bounds: got %v, want %dx%d", tt.name, tt.got.Bounds(), wantW, wantH)
		}
		for y, row := range tt.want {
			for x, c := range row {
				if g := tt.got.RGBAAt(x, y); g != c {
					t.Errorf("%s: pixel (%d,%d): got %v, want %v", tt.name, x, y, g, c)
				}
			}
		}
	}
}

func TestRotateQuarterTurnsSubImage(t *testing.T) {
	// A sub-image has bounds that do not start at (0,0).
	big := image.NewRGBA(image.Rect(0, 0, 10, 10))
	grid := newTestGrid()
	draw.Draw(big, image.Rect(4, 5, 7, 7), grid, image.Point{}, draw.Src)
	sub, ok := big.SubImage(image.Rect(4, 5, 7, 7)).(*image.RGBA)
	if !ok {
		t.Fatal("SubImage did not return *image.RGBA")
	}
	assertSameImage(t, "Rotate90 sub-image", Rotate90(sub), Rotate90(grid))
}

func TestRotateMultiplesOf90UseLosslessPath(t *testing.T) {
	src := newTestGrid()
	tests := []struct {
		degrees float64
		want    *image.RGBA
	}{
		{0, src},
		{360, src},
		{-360, src},
		{90, Rotate90(src)},
		{450, Rotate90(src)},
		{-270, Rotate90(src)},
		{180, Rotate180(src)},
		{-180, Rotate180(src)},
		{270, Rotate270(src)},
		{-90, Rotate270(src)},
	}
	for _, tt := range tests {
		assertSameImage(t, "Rotate", Rotate(src, tt.degrees, white, nil), tt.want)
		assertSameImage(t, "RotateCrop", RotateCrop(src, tt.degrees, nil), tt.want)
	}
}

func TestRotateExpandsCanvas(t *testing.T) {
	src := newSolid(100, 50, red)
	got := Rotate(src, 45, white, draw.BiLinear)

	// 150 * cos(45°) = 106.07, rounded up.
	if got.Bounds() != image.Rect(0, 0, 107, 107) {
		t.Fatalf("bounds: got %v, want 107x107", got.Bounds())
	}
	b := got.Bounds()
	for _, p := range []image.Point{{0, 0}, {b.Max.X - 1, 0}, {0, b.Max.Y - 1}, {b.Max.X - 1, b.Max.Y - 1}} {
		if c := got.RGBAAt(p.X, p.Y); c != white {
			t.Errorf("corner %v: got %v, want background %v", p, c, white)
		}
	}
	if c := got.RGBAAt(53, 53); c != red {
		t.Errorf("center: got %v, want %v", c, red)
	}
}

func TestRotateNilBackgroundIsTransparent(t *testing.T) {
	got := Rotate(newSolid(40, 40, red), 30, nil, nil)
	if c := got.RGBAAt(0, 0); c.A != 0 {
		t.Errorf("corner: got %v, want transparent", c)
	}
}

func TestRotateDirection(t *testing.T) {
	// Left half red, right half blue. Turning clockwise moves the left
	// half to the top; turning counterclockwise moves it to the bottom.
	src := newHalves(100, 100, red, blue)
	for _, fn := range []struct {
		name   string
		rotate func(float64) *image.RGBA
	}{
		{"Rotate", func(d float64) *image.RGBA { return Rotate(src, d, white, nil) }},
		{"RotateCrop", func(d float64) *image.RGBA { return RotateCrop(src, d, nil) }},
	} {
		cw := fn.rotate(80)
		cx, cy := cw.Bounds().Dx()/2, cw.Bounds().Dy()/2
		if c := cw.RGBAAt(cx, cy-15); c != red {
			t.Errorf("%s clockwise: above center got %v, want %v", fn.name, c, red)
		}
		if c := cw.RGBAAt(cx, cy+15); c != blue {
			t.Errorf("%s clockwise: below center got %v, want %v", fn.name, c, blue)
		}

		ccw := fn.rotate(-80)
		cx, cy = ccw.Bounds().Dx()/2, ccw.Bounds().Dy()/2
		if c := ccw.RGBAAt(cx, cy-15); c != blue {
			t.Errorf("%s counterclockwise: above center got %v, want %v", fn.name, c, blue)
		}
		if c := ccw.RGBAAt(cx, cy+15); c != red {
			t.Errorf("%s counterclockwise: below center got %v, want %v", fn.name, c, red)
		}
	}
}

func TestRotateCropHasNoBackground(t *testing.T) {
	sizes := []image.Point{{100, 100}, {200, 80}, {80, 200}, {640, 480}}
	angles := []float64{1, 5, 15, 30, 45, 60, 89, 135, -10, -45, 200}
	for _, s := range sizes {
		src := newSolid(s.X, s.Y, green)
		for _, deg := range angles {
			got := RotateCrop(src, deg, draw.CatmullRom)
			cw, ch := RotateCropSize(float64(s.X), float64(s.Y), deg)
			want := image.Rect(0, 0, int(math.Floor(cw+rotateEpsilon)), int(math.Floor(ch+rotateEpsilon)))
			if got.Bounds() != want {
				t.Fatalf("%v at %v°: bounds got %v, want %v", s, deg, got.Bounds(), want)
			}
			b := got.Bounds()
			for y := b.Min.Y; y < b.Max.Y; y++ {
				for x := b.Min.X; x < b.Max.X; x++ {
					if c := got.RGBAAt(x, y); c.A != 255 {
						t.Fatalf("%v at %v°: pixel (%d,%d) not opaque: %v", s, deg, x, y, c)
					}
				}
			}
		}
	}
}

func TestRotateCropSize(t *testing.T) {
	tests := []struct {
		w, h, deg    float64
		wantW, wantH float64
	}{
		{100, 50, 0, 100, 50},
		{100, 50, 180, 100, 50},
		{100, 50, 90, 50, 100},
		// Square at 45°: side / sqrt(2).
		{100, 100, 45, 100 / math.Sqrt2, 100 / math.Sqrt2},
		// Wide image at 45°: limited by the short side.
		{200, 100, 45, 50 * math.Sqrt2, 50 * math.Sqrt2},
		{0, 100, 30, 0, 0},
	}
	for _, tt := range tests {
		gw, gh := RotateCropSize(tt.w, tt.h, tt.deg)
		if math.Abs(gw-tt.wantW) > 1e-6 || math.Abs(gh-tt.wantH) > 1e-6 {
			t.Errorf("RotateCropSize(%v, %v, %v): got %.4f x %.4f, want %.4f x %.4f",
				tt.w, tt.h, tt.deg, gw, gh, tt.wantW, tt.wantH)
		}
	}
}

func TestRotateJPEGFile(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "in.jpg")
	f, err := os.Create(in)
	if err != nil {
		t.Fatal(err)
	}
	if err := jpeg.Encode(f, newHalves(120, 80, red, blue), &jpeg.Options{Quality: 95}); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	img, format, err := ReadImageFile(in)
	if err != nil {
		t.Fatal(err)
	}
	if format != "jpeg" {
		t.Fatalf("format: got %q, want jpeg", format)
	}

	tests := []struct {
		name string
		img  image.Image
	}{
		{"rotate90.jpg", Rotate90(img)},
		{"rotate.jpg", Rotate(img, 15, white, draw.CatmullRom)},
		{"rotatecrop.jpg", RotateCrop(img, 15, draw.CatmullRom)},
	}
	for _, tt := range tests {
		out := filepath.Join(dir, tt.name)
		if err := (Image{Image: tt.img}).WriteJPEGFile(out, nil); err != nil {
			t.Fatalf("%s: %v", tt.name, err)
		}
		cfg, _, err := DecodeConfigFile(out)
		if err != nil {
			t.Fatalf("%s: %v", tt.name, err)
		}
		if cfg.Width != tt.img.Bounds().Dx() || cfg.Height != tt.img.Bounds().Dy() {
			t.Errorf("%s: size got %dx%d, want %v", tt.name, cfg.Width, cfg.Height, tt.img.Bounds().Size())
		}
	}
}
