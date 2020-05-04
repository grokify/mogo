package colors

import (
	"fmt"
	"image/color"
	"regexp"
	"strconv"
	"strings"

	"github.com/grokify/gotilla/errors/errorsutil"
	"golang.org/x/image/colornames"
)

var (
	rxColorHex  = regexp.MustCompile(`^(?:#|0x)?([0-9a-f]{2})([0-9a-f]{2})([0-9a-f]{2})$`)
	rxGoogleCol = regexp.MustCompile(`^google([0-9]+)$`)
)

// Parse returns a `color.RGBA` given a color name
// or hex color code.
func Parse(colorName string) (color.RGBA, error) {
	colorName = strings.ToLower(strings.TrimSpace(colorName))
	if col, ok := colornames.Map[colorName]; ok {
		return col, nil
	}
	colorHtml, err := HexToColor(colorName)
	if err == nil {
		return colorHtml, nil
	}
	colorGoog, err := GoogleToColor(colorName)
	if err == nil {
		return colorGoog, nil
	}
	return color.RGBA{}, fmt.Errorf("E_COLOR_NOT_FOUND [%s]", colorName)
}

// Parse returns a `color.RGBA` given a hex color code.
func HexToColor(hexRGB string) (color.RGBA, error) {
	hexRGB = strings.ToLower(strings.TrimSpace(hexRGB))
	m := rxColorHex.FindStringSubmatch(hexRGB)
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

func GoogleToColor(googString string) (color.RGBA, error) {
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
