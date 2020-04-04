package imageutil

import (
	"image"
	"image/jpeg"
	"image/png"
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
		quality = 100
	}
	opt.Quality = quality

	return jpeg.Encode(out, img, &opt)
}

func WriteFilePNG(filename string, img image.Image) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = png.Encode(out, img)
	if err != nil {
		return err
	}
	return out.Close()
}
