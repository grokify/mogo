package imageutil

import (
	"image"
	"image/color"
	"math"

	"golang.org/x/image/draw"
	"golang.org/x/image/math/f64"
)

// rotateEpsilon is the tolerance used when comparing angles and sizes, so
// that floating point noise (e.g. sin(180°) != 0) does not change results.
const rotateEpsilon = 1e-9

// Rotate90 rotates an image 90° clockwise without any loss of quality.
// The output width and height are the input height and width.
func Rotate90(src image.Image) *image.RGBA {
	return rotateRight(src, 1)
}

// Rotate180 rotates an image 180° without any loss of quality.
func Rotate180(src image.Image) *image.RGBA {
	return rotateRight(src, 2)
}

// Rotate270 rotates an image 270° clockwise (90° counterclockwise) without
// any loss of quality. The output width and height are the input height
// and width.
func Rotate270(src image.Image) *image.RGBA {
	return rotateRight(src, 3)
}

// rotateRight rotates src by quarterTurns * 90° clockwise by moving pixels,
// so there is no interpolation. The output bounds start at (0, 0).
func rotateRight(src image.Image, quarterTurns int) *image.RGBA {
	in := ImageToRGBA(src)
	b := in.Bounds()
	w, h := b.Dx(), b.Dy()

	quarterTurns = ((quarterTurns % 4) + 4) % 4
	dw, dh := w, h
	if quarterTurns%2 == 1 {
		dw, dh = h, w
	}
	out := image.NewRGBA(image.Rect(0, 0, dw, dh))

	for y := range h {
		for x := range w {
			var dx, dy int
			switch quarterTurns {
			case 0:
				dx, dy = x, y
			case 1:
				dx, dy = h-1-y, x
			case 2:
				dx, dy = w-1-x, h-1-y
			case 3:
				dx, dy = y, w-1-x
			}
			si := in.PixOffset(b.Min.X+x, b.Min.Y+y)
			di := out.PixOffset(dx, dy)
			copy(out.Pix[di:di+4], in.Pix[si:si+4])
		}
	}
	return out
}

// Rotate rotates an image by any number of degrees, like Photoshop's
// "Image Rotation > Arbitrary". Positive degrees rotate clockwise and
// negative degrees rotate counterclockwise.
//
// The canvas grows to hold the whole rotated image, and the exposed corners
// are filled with bgColor. If bgColor is nil the corners are transparent,
// which encodes as black in formats without alpha such as JPEG.
//
// interp sets how pixels are sampled, e.g. draw.BiLinear (used when nil)
// or draw.CatmullRom for higher quality. Multiples of 90° use the lossless
// Rotate90, Rotate180 and Rotate270 and ignore interp and bgColor.
func Rotate(src image.Image, degrees float64, bgColor color.Color, interp draw.Interpolator) *image.RGBA {
	if turns, ok := quarterTurns(degrees); ok {
		return rotateRight(src, turns)
	}
	w, h := float64(src.Bounds().Dx()), float64(src.Bounds().Dy())
	sin, cos := math.Sincos(degreesToRadians(degrees))
	sin, cos = math.Abs(sin), math.Abs(cos)
	dw := ceilSize(w*cos + h*sin)
	dh := ceilSize(w*sin + h*cos)

	dst := image.NewRGBA(image.Rect(0, 0, dw, dh))
	if bgColor != nil {
		draw.Draw(dst, dst.Bounds(), image.NewUniform(bgColor), image.Point{}, draw.Src)
	}
	rotateInto(dst, src, degrees, interp)
	return dst
}

// RotateCrop rotates an image by any number of degrees, then crops it to the
// largest upright rectangle that holds only image pixels, so no background
// is exposed. This is like straightening a photo in Photoshop with the crop
// tool's "Delete Cropped Pixels". Positive degrees rotate clockwise and
// negative degrees rotate counterclockwise.
//
// The output is smaller than the input for any angle that is not a multiple
// of 90°. interp works as in Rotate. Multiples of 90° use the lossless
// Rotate90, Rotate180 and Rotate270.
func RotateCrop(src image.Image, degrees float64, interp draw.Interpolator) *image.RGBA {
	if turns, ok := quarterTurns(degrees); ok {
		return rotateRight(src, turns)
	}
	w, h := float64(src.Bounds().Dx()), float64(src.Bounds().Dy())
	cw, ch := RotateCropSize(w, h, degrees)
	// Floor so every output pixel lies inside the rotated image.
	dst := image.NewRGBA(image.Rect(0, 0, int(math.Floor(cw+rotateEpsilon)), int(math.Floor(ch+rotateEpsilon))))
	rotateInto(dst, src, degrees, interp)
	return dst
}

// RotateCropSize returns the width and height of the largest upright
// rectangle that fits inside a w x h rectangle rotated by degrees. It is
// the size RotateCrop uses before rounding down to whole pixels.
func RotateCropSize(w, h, degrees float64) (float64, float64) {
	if w <= 0 || h <= 0 {
		return 0, 0
	}
	sin, cos := math.Sincos(degreesToRadians(degrees))
	sin, cos = math.Abs(sin), math.Abs(cos)

	long, short := w, h
	widthIsLonger := w >= h
	if !widthIsLonger {
		long, short = h, w
	}

	if short <= 2*sin*cos*long || math.Abs(sin-cos) < rotateEpsilon {
		// The crop touches the two long sides of the rotated rectangle.
		x := 0.5 * short
		if widthIsLonger {
			return x / sin, x / cos
		}
		return x / cos, x / sin
	}
	// The crop touches all four sides of the rotated rectangle.
	cos2 := cos*cos - sin*sin
	return (w*cos - h*sin) / cos2, (h*cos - w*sin) / cos2
}

// rotateInto draws src rotated clockwise by degrees about its center onto
// the center of dst.
func rotateInto(dst *image.RGBA, src image.Image, degrees float64, interp draw.Interpolator) {
	if interp == nil {
		interp = draw.BiLinear
	}
	sr := src.Bounds()
	db := dst.Bounds()
	sin, cos := math.Sincos(degreesToRadians(degrees))

	scx := float64(sr.Min.X+sr.Max.X) / 2
	scy := float64(sr.Min.Y+sr.Max.Y) / 2
	dcx := float64(db.Min.X+db.Max.X) / 2
	dcy := float64(db.Min.Y+db.Max.Y) / 2

	// Map source to destination: move the source center to the origin,
	// rotate, then move to the destination center. With y pointing down,
	// this matrix turns the image clockwise for positive angles.
	s2d := f64.Aff3{
		cos, -sin, dcx - (cos*scx - sin*scy),
		sin, cos, dcy - (sin*scx + cos*scy),
	}
	interp.Transform(dst, s2d, src, sr, draw.Over, nil)
}

// quarterTurns reports whether degrees is a multiple of 90 and, if so, how
// many clockwise quarter turns (0-3) it equals.
func quarterTurns(degrees float64) (int, bool) {
	q := degrees / 90
	r := math.Round(q)
	if math.Abs(q-r) > rotateEpsilon {
		return 0, false
	}
	return ((int(math.Mod(r, 4)) % 4) + 4) % 4, true
}

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// ceilSize rounds a size up to whole pixels, ignoring floating point noise.
func ceilSize(v float64) int {
	return int(math.Ceil(v - rotateEpsilon))
}
