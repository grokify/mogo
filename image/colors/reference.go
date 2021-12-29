package colors

import (
	"image/color"

	"github.com/grokify/mogo/math/mathutil"
)

// ChartColor1 is the color palette for Google Charts as collected by
// Craig Davis here: https://gist.github.com/there4/2579834
var GoogleChartColorsHex = [...]string{
	"#3366CC",
	"#DC3912",
	"#FF9900",
	"#109618",
	"#990099",
	"#3B3EAC",
	"#0099C6",
	"#DD4477",
	"#66AA00",
	"#B82E2E",
	"#316395",
	"#994499",
	"#22AA99",
	"#AAAA11",
	"#6633CC",
	"#E67300",
	"#8B0707",
	"#329262",
	"#5574A6",
	"#3B3EAC",
}

func GetGoogleChartColors() []color.RGBA {
	rgbas := []color.RGBA{}
	for _, hex := range GoogleChartColorsHex {
		rgb, err := ParseHex(hex)
		if err != nil {
			panic(err)
		}
		rgbas = append(rgbas, rgb)
	}
	return rgbas
}

var GoogleChartColors = GetGoogleChartColors()

func GoogleChartColorX(index uint64) color.RGBA {
	_, remainder := mathutil.DivideInt64(int64(index),
		int64(len(GoogleChartColorsHex)))
	return GoogleChartColors[remainder]
}
