package imageutil

import (
	"image/gif"
	"image/jpeg"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/grokify/mogo/io/ioutil"
	"github.com/grokify/mogo/os/osutil"
)

const (
	JPEGExt            = ".jpeg"
	JPEGExtJPG         = ".jpg"
	JPEGQualityDefault = jpeg.DefaultQuality // 75
	JPEGQualityMax     = 100
	JPEGQualityMin     = 1
	JPEGMarkerPrefix   = 0xff
	JPEGMarkerExif     = 0xe1
	JPEGMarkerSOI      = 0xd8
)

var rxJPEGExtension = regexp.MustCompile(`(?i)\.(jpg|jpeg)$`)

func JPEGMarker(b byte) []byte {
	return []byte{JPEGMarkerPrefix, b}
}

func ReadDirJPEGFiles(dir string, rx *regexp.Regexp) (osutil.DirEntries, error) {
	if rx == nil {
		rx = rxJPEGExtension
	}
	return osutil.ReadDirMore(dir, rx, false, true, false)
}

func WriteGIFFile(filename string, img *gif.GIF) error {
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

func ResizeFileJPEG(inputFile, outputFile string, outputWidth, outputHeight uint, opt *JPEGEncodeOptions) error {
	img, _, err := ReadImageFile(inputFile)
	if err != nil {
		return err
	}
	img2 := Resize(outputWidth, outputHeight, img, ScalerBest())
	return writeJPEGFile(outputFile, img2, opt)
}

type JPEGEncodeOptions struct {
	Options            *jpeg.Options
	Exif               []byte
	ReadFilenameRegexp *regexp.Regexp
	WriteExtension     string
}

func (opts JPEGEncodeOptions) ReadFilenameRegexpOrDefault() *regexp.Regexp {
	if opts.ReadFilenameRegexp != nil {
		return opts.ReadFilenameRegexp
	}
	return rxJPEGExtension
}

func (opts JPEGEncodeOptions) WriteExtensionOrDefault() string {
	if strings.TrimSpace(opts.WriteExtension) == "" {
		return JPEGExtJPG // default to 3 letter extension to support Microsoft.
	}
	return opts.WriteExtension
}

var JPEGEncodeOptionsQualityMax = &JPEGEncodeOptions{
	Options: &jpeg.Options{
		Quality: JPEGQualityMax}}

// newWriterExif is used to write Exif to an `io.Writer` before calling `jpeg.Encode()`.
// It is used with `jpeg.Encode()` to remove the Start of Image (SOI) marker after adding
// SOI and Exif.
func newWriterExif(w io.Writer, exif []byte) (io.Writer, error) {
	// Adapted from the following under MIT license: https://github.com/jdeng/goheif/blob/a0d6a8b3e68f9d613abd9ae1db63c72ba33abd14/heic2jpg/main.go
	// See more here: https://en.wikipedia.org/wiki/JPEG_File_Interchange_Format
	// https://www.codeproject.com/Articles/47486/Understanding-and-Reading-Exif-Data
	wExif := ioutil.NewSkipWriter(w, 2)

	if _, err := w.Write(JPEGMarker(JPEGMarkerSOI)); err != nil {
		return nil, err
	}

	if exif == nil {
		return wExif, nil
	}

	markerLen := 2 + len(exif)
	marker := []byte{
		JPEGMarkerPrefix,
		JPEGMarkerExif,
		uint8(markerLen >> 8),
		uint8(markerLen & JPEGMarkerPrefix)}
	if _, err := w.Write(marker); err != nil {
		return nil, err
	}

	if _, err := w.Write(exif); err != nil {
		return nil, err
	}
	return wExif, nil
}
