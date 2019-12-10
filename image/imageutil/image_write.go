package imageutil

import (
	"image"
	"image/jpeg"
	"os"
)

const DefaultQualityJPEG int = 80

func WriteFileJPEG(filename string, img image.Image, quality int) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}

	var opt jpeg.Options
	if quality <= 0 {
		quality = DefaultQualityJPEG
	}
	if quality > 100 {
		quality = DefaultQualityJPEG
	}
	opt.Quality = quality

	return jpeg.Encode(out, img, &opt)
}
