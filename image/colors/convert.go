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

func ColorToHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("%02x%02x%02x", uint8(r), uint8(g), uint8(b))
}

func ColorAverageImage(i image.Image) color.Color {
	min := i.Bounds().Min
	max := i.Bounds().Max
	r, b, g, n := 0.0, 0.0, 0.0, 0.0
	for x := min.X; x < max.X; x++ {
		for y := min.Y; y < max.Y; y++ {
			c := i.At(x, y)
			ri, gi, bi, _ := c.RGBA()
			r += float64(uint8(ri) * uint8(ri))
			g += float64(uint8(gi) * uint8(gi))
			b += float64(uint8(bi) * uint8(bi))
			n++
		}
	}
	return color.RGBA{
		R: uint8(math.Sqrt(r / n)),
		G: uint8(math.Sqrt(g / n)),
		B: uint8(math.Sqrt(b / n)),
		A: 255}
}

func ColorAverage(c ...color.Color) color.Color {
	if len(c) == 0 {
		return color.Black
	}
	r, b, g := 0.0, 0.0, 0.0
	for _, ci := range c {
		ri, gi, bi, _ := ci.RGBA()
		r += float64(uint8(ri) * uint8(ri))
		g += float64(uint8(gi) * uint8(gi))
		b += float64(uint8(bi) * uint8(bi))
	}
	n := float64(len(c))
	return color.RGBA{
		R: uint8(math.Sqrt(r / n)),
		G: uint8(math.Sqrt(g / n)),
		B: uint8(math.Sqrt(b / n)),
		A: 255}
}
