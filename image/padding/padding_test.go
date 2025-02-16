package padding

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"os"
	"reflect"
	"testing"

	"github.com/grokify/mogo/image/colors"
)

var paddingTests = []struct {
	filename                string
	topPadding              int
	rightPadding            int
	bottomPadding           int
	leftPadding             int
	isPaddingFunc           IsPaddingFunc
	nonPaddngColorHistogram map[string]uint
}{
	{"testdata/padding_example.png", 5, 10, 15, 20, CreateIsPaddingFuncSimple(color.White), map[string]uint{"00000.00000.00000.65535": 2500}},
}

func TestPadding(t *testing.T) {
	for _, tt := range paddingTests {
		infile, err := os.Open(tt.filename)
		if err != nil {
			panic(err)
		}
		defer infile.Close()

		img, _, err := image.Decode(infile)
		if err != nil {
			panic(err)
		}
		tryTop, tryRight, tryBottom, tryLeft := PaddingWidths(img, tt.isPaddingFunc)
		if tryTop != tt.topPadding || tryRight != tt.rightPadding || tryBottom != tt.bottomPadding || tryLeft != tt.leftPadding {
			t.Errorf("PaddingWidths({img}, {func}): want [%d,%d,%d,%d] got [%d,%d,%d,%d]",
				tt.topPadding, tt.rightPadding, tt.bottomPadding, tt.leftPadding, tryTop, tryRight, tryBottom, tryLeft)
		}

		tryRect := NonPaddingRectangle(img, tt.isPaddingFunc)
		new := image.NewRGBA(tryRect)
		draw.Draw(new, new.Bounds(), img, tryRect.Min, draw.Over)
		cs := imageColors(new)
		if !reflect.DeepEqual(tt.nonPaddngColorHistogram, cs) {
			t.Errorf("PaddingWidths({img}, {func}): want remaining color map mismatch [%v] [%v]",
				tt.nonPaddngColorHistogram, cs)
		}
	}
}

func imageColors(im image.Image) map[string]uint {
	cm := map[string]uint{}
	for xi := im.Bounds().Min.X; xi < im.Bounds().Max.X; xi++ {
		for yi := im.Bounds().Min.Y; yi < im.Bounds().Max.Y; yi++ {
			cm[colors.ColorString(im.At(xi, yi))]++
		}
	}
	return cm
}
