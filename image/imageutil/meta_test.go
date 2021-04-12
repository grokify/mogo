package imageutil

import (
	"image"
	"testing"
)

var imageMetaTests = []struct {
	x0     int
	y0     int
	x1     int
	y1     int
	width  int
	height int
}{
	{0, 0, 1, 1, 1, 1},
	{0, 0, 2, 2, 2, 2},
	{0, 0, 1000, 500, 1000, 500},
}

func TestImageMeta(t *testing.T) {
	for _, tt := range imageMetaTests {
		rect := image.Rect(tt.x0, tt.y0, tt.x1, tt.y1)
		gotStats := ImageStatsRect(rect)
		if gotStats.Width != tt.width || gotStats.Height != tt.height {
			t.Errorf("imageutil.ImageStatsRect(%d,%d,%d,%d) Mismatch: got[%d, %d] want [%d, %d]",
				tt.x0, tt.y0, tt.x1, tt.y1,
				gotStats.Width, gotStats.Height,
				tt.width, tt.height)
		}
	}
}
