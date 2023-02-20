package imageutil

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
)

const (
	JPEGQualityDefault int = 80
	JPEGQualityMax     int = 100
)

func WriteFileGIF(filename string, img *gif.GIF) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	err = gif.EncodeAll(f, img)
	if err != nil {
		return err
	}
	return f.Close()
}

func ResizeFileJPEG(inputFile, outputFile string, outputWidth, outputHeight uint, quality int) error {
	img, _, err := ReadImageFile(inputFile)
	if err != nil {
		return err
	}
	img2 := Resize(outputWidth, outputHeight, img, ScalerBest())
	return WriteFileJPEG(outputFile, img2, quality)
}

func WriteFileJPEG(filename string, img image.Image, quality int) error {
	out, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	return WriteJPEG(out, img, quality)
	/*
		var opt jpeg.Options
		if quality <= 0 {
			quality = JPEGQualityDefault
		}
		if quality > JPEGQualityMax {
			quality = JPEGQualityMax
		}
		opt.Quality = quality

		return jpeg.Encode(out, img, &opt)
	*/
}

func WriteJPEG(w io.Writer, img image.Image, quality int) error {
	/*
		out, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return err
		}
	*/

	var opt jpeg.Options
	if quality <= 0 {
		quality = JPEGQualityDefault
	}
	if quality > JPEGQualityMax {
		quality = JPEGQualityMax
	}
	opt.Quality = quality

	return jpeg.Encode(w, img, &opt)
}

func BytesJPEG(img image.Image, quality int) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := WriteJPEG(buf, img, quality)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func WriteFilePNG(filename string, img image.Image) error {
	out, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	err = png.Encode(out, img)
	if err != nil {
		return err
	}
	return out.Close()
}
