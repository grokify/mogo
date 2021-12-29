package colors

import (
	"image/color"
	"strings"
	"testing"

	"golang.org/x/image/colornames"
)

var colorConvertTests = []struct {
	svgName string
	hexRaw  string
	hexHTML string
	hexNum  string
	hexInt  uint32
	R       uint8
	G       uint8
	B       uint8
}{
	{"black", "000000", "#000000", "0x000000", 0x000000, 0, 0, 0},
	{"white", "ffffff", "#ffffff", "0xffffff", 0xffffff, 255, 255, 255},
	{"Darkgrey", "a9a9a9", "#a9a9a9", "0xa9a9a9", 0xa9a9a9, 169, 169, 169},
	{"google0", "3366CC", "#3366CC", "0x3366CC", 0x3366CC, 51, 102, 204},
	{"google20", "3366CC", "#3366CC", "0x3366CC", 0x3366CC, 51, 102, 204},
	{"google40", "3366CC", "#3366CC", "0x3366CC", 0x3366CC, 51, 102, 204},
}

func TestConvertColors(t *testing.T) {
	for _, tt := range colorConvertTests {
		gotSVG, err := Parse(tt.svgName)
		if err != nil {
			t.Errorf("Parse('%s'): want [%v] error [%v]", tt.svgName, color.RGBA{
				tt.R, tt.G, tt.B, 0xff}, err.Error())
		} else if gotSVG.R != tt.R || gotSVG.G != tt.G || gotSVG.B != tt.B {
			t.Errorf("Parse('%s'): want [%v] got [%v]", tt.svgName, color.RGBA{
				tt.R, tt.G, tt.B, 0xff}, gotSVG)
		}
		gotHexRaw, err := Parse(tt.hexRaw)
		if err != nil {
			t.Errorf("Parse('%s'): want [%v] error [%v]", tt.hexRaw, color.RGBA{
				tt.R, tt.G, tt.B, 0xff}, err.Error())
		} else if gotHexRaw.R != tt.R || gotHexRaw.G != tt.G || gotHexRaw.B != tt.B {
			t.Errorf("Parse('%s'): want [%v] got [%v]", tt.hexRaw, color.RGBA{
				tt.R, tt.G, tt.B, 0xff}, gotHexRaw)
		}
		gotHexHTML, err := Parse(tt.hexHTML)
		if err != nil {
			t.Errorf("Parse('%s'): want [%v] error [%v]", tt.hexHTML, color.RGBA{
				tt.R, tt.G, tt.B, 0xff}, err.Error())
		} else if gotHexHTML.R != tt.R || gotHexHTML.G != tt.G || gotHexHTML.B != tt.B {
			t.Errorf("Parse('%s'): want [%v] got [%v]", tt.hexHTML, color.RGBA{
				tt.R, tt.G, tt.B, 0xff}, gotHexHTML)
		}
		gotHexNum, err := Parse(tt.hexNum)
		if err != nil {
			t.Errorf("Parse('%s'): want [%v] error [%v]", tt.hexNum, color.RGBA{
				tt.R, tt.G, tt.B, 0xff}, err.Error())
		} else if gotHexNum.R != tt.R || gotHexNum.G != tt.G || gotHexNum.B != tt.B {
			t.Errorf("Parse('%s'): want [%v] got [%v]", tt.hexNum, color.RGBA{
				tt.R, tt.G, tt.B, 0xff}, gotHexNum)
		}
		clr := color.RGBA{
			R: uint8(tt.R),
			G: uint8(tt.G),
			B: uint8(tt.B),
			A: 0xff}
		gotRGBAHexLc := ColorRGBAToHex(clr)
		if gotRGBAHexLc != strings.ToLower(tt.hexRaw) {
			t.Errorf("ColorRGBAToHex(%v): want [%v] got [%v]", clr, strings.ToLower(tt.hexRaw), gotRGBAHexLc)
		}
		gotHexLc := ColorToHex(clr)
		if gotHexLc != strings.ToLower(tt.hexRaw) {
			t.Errorf("ColorToHex(%v): want [%v] got [%v]", clr, strings.ToLower(tt.hexRaw), gotHexLc)
		}
	}
}

var averageTests = []struct {
	colors     []string
	averageHex string
}{
	{[]string{"black", "black"}, "000000"},
	{[]string{"white", "white"}, "ffffff"},
	{[]string{"black", "white"}, "b4b4b4"},
}

func TestAverage(t *testing.T) {
	for _, tt := range averageTests {
		ttColors := []color.Color{}
		for _, colorname := range tt.colors {
			if crgba, ok := colornames.Map[colorname]; ok {
				ttColors = append(ttColors, crgba)
			}
		}
		clrAvg := ColorAverage(ttColors...)
		clrAvgHex := ColorToHex(clrAvg)
		if clrAvgHex != tt.averageHex {
			t.Errorf("ColorToHex(...): want [%v] got [%v]", tt.averageHex, clrAvgHex)
		}
	}
}
