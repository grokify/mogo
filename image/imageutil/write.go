package imageutil

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

const (
	JPEGQualityDefault int = 80
	JPEGQualityMax     int = 100
)

func WriteFileJPEG(filename string, img image.Image, quality int) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}

	var opt jpeg.Options
	if quality <= 0 {
		quality = JPEGQualityDefault
	}
	if quality > JPEGQualityMax {
		quality = JPEGQualityMax
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
