package imageutil

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/chai2010/webp"
)

const (
	JPEGQualityDefault int = 80
	JPEGQualityMax     int = 100
)

func ResizeFileJPEG(inputFile, outputFile string, outputWidth, outputHeight uint, quality int) error {
	img, _, err := ReadImageFile(inputFile)
	if err != nil {
		return err
	}
	img2 := Resize(outputWidth, outputHeight, img, ScalerBest())
	return WriteFileJPEG(outputFile, img2, quality)
}

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

func WriteFileWEBP(filename string, img image.Image, lossless bool, perm os.FileMode) error {
	var buf bytes.Buffer
	err := webp.Encode(&buf, img, &webp.Options{Lossless: lossless})
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, buf.Bytes(), perm)
}
