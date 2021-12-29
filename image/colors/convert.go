package colors

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/grokify/mogo/errors/errorsutil"
	"golang.org/x/image/colornames"
)

var (
	rxColorHex  = regexp.MustCompile(`^(?:#|0x)?([0-9a-f]{2})([0-9a-f]{2})([0-9a-f]{2})$`)
	rxGoogleCol = regexp.MustCompile(`^google([0-9]+)$`)
)

// Parse returns a `color.RGBA` given a color name or hex color code.
func Parse(colorName string) (color.RGBA, error) {
	colorName = strings.ToLower(strings.TrimSpace(colorName))
	if col, ok := colornames.Map[colorName]; ok {
		return col, nil
	}
	colorHtml, err := ParseHex(colorName)
	if err == nil {
		return colorHtml, nil
	}
	colorGoog, err := ParseGoogle(colorName)
	if err == nil {
		return colorGoog, nil
	}
	return color.RGBA{}, fmt.Errorf("E_COLOR_NOT_FOUND [%s]", colorName)
}

// MustParse returns a `color.RGBA` given a hex color code or Google color string.
// It panics if the input string cannot be parsed.
func MustParse(input string) color.RGBA {
	c, err := Parse(input)
	if err != nil {
		panic(err)
	}
	return c
}

// ParseHex returns a `color.RGBA` given a hex color code.
func ParseHex(hexRGB string) (color.RGBA, error) {
	m := rxColorHex.FindStringSubmatch(strings.ToLower(strings.TrimSpace(hexRGB)))
	if len(m) == 0 {
		return color.RGBA{}, fmt.Errorf("E_COLOR_NOT_HEX_STRING [%s]", hexRGB)
	}
	rdecimal, errR := strconv.ParseUint(m[1], 16, 64)
	gdecimal, errG := strconv.ParseUint(m[2], 16, 64)
	bdecimal, errB := strconv.ParseUint(m[3], 16, 64)
	err := errorsutil.Join(false, errR, errG, errB)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("E_COLOR_NOT_HEX_PARSE [%s]", err.Error())
	}
	return color.RGBA{
		R: uint8(rdecimal),
		G: uint8(gdecimal),
		B: uint8(bdecimal),
		A: 0xff}, nil
}

func ParseGoogle(googString string) (color.RGBA, error) {
	m := rxGoogleCol.FindStringSubmatch(googString)
	if len(m) == 0 {
		return color.RGBA{},
			fmt.Errorf("E_COLOR_NOT_GOOG_MATCH [%s]", googString)
	}
	idxInt, err := strconv.Atoi(m[1])
	if err != nil {
		panic(err)
	}
	col := GoogleChartColorX(uint64(idxInt))
	return col, nil
}

func ColorRGBAToHex(c color.RGBA) string {
	return fmt.Sprintf("%02x%02x%02x", c.R, c.G, c.B)
}

// ColorToHex returns 6 byte hex code in lower case.
func ColorToHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("%02x%02x%02x", uint8(r), uint8(g), uint8(b))
}

func ColorAverageImage(i image.Image) color.Color {
	var r, g, b uint64

	bounds := i.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pr, pg, pb, _ := i.At(x, y).RGBA()

			r += uint64(pr * pr)
			g += uint64(pg * pg)
			b += uint64(pb * pg)
		}
	}

	d := uint64(bounds.Dy() * bounds.Dx())

	r /= d
	g /= d
	b /= d

	return color.RGBA{
		uint8(math.Sqrt(float64(r)) / 0x101),
		uint8(math.Sqrt(float64(g)) / 0x101),
		uint8(math.Sqrt(float64(b)) / 0x101),
		uint8(255)}
}

/*

https://jimsaunders.net/2015/05/22/manipulating-colors-in-go.html
https://sighack.com/post/averaging-rgb-colors-the-right-way

*/

func ColorAverage(c ...color.Color) color.Color {
	if len(c) == 0 {
		return color.Black
	}
	var r, g, b uint64

	for _, ci := range c {
		pr, pg, pb, _ := ci.RGBA()

		r += uint64(pr * pr)
		g += uint64(pg * pg)
		b += uint64(pb * pg)
	}

	d := uint64(len(c))

	r /= d
	g /= d
	b /= d

	return color.RGBA{
		uint8(math.Sqrt(float64(r)) / 0x101),
		uint8(math.Sqrt(float64(g)) / 0x101),
		uint8(math.Sqrt(float64(b)) / 0x101),
		uint8(255)}
}
