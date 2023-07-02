package imageutil

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"

	"github.com/grokify/mogo/net/http/httputilmore"
)

var ErrImageNotSet = errors.New("image not set")

type ImageMore struct {
	Image image.Image
}

func (im ImageMore) BytesJPEG(opt *JPEGEncodeOptions) ([]byte, error) {
	return bytesJPEG(im.Image, opt)
}

func bytesJPEG(img image.Image, opt *JPEGEncodeOptions) ([]byte, error) {
	if img == nil {
		return []byte{}, ErrImageNotSet
	}
	buf := new(bytes.Buffer)
	err := writeJPEG(buf, img, opt)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func (im ImageMore) BytesPNG() ([]byte, error) {
	return bytesPNG(im.Image)
}

func bytesPNG(img image.Image) ([]byte, error) {
	if img == nil {
		return []byte{}, ErrImageNotSet
	}
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func (im ImageMore) WriteJPEG(w io.Writer, opt *JPEGEncodeOptions) error {
	return writeJPEG(w, im.Image, opt)
}

func writeJPEG(w io.Writer, img image.Image, opt *JPEGEncodeOptions) error {
	if img == nil {
		return ErrImageNotSet
	}
	if opt != nil && len(opt.Exif) > 0 {
		if wexif, err := newWriterExif(w, opt.Exif); err != nil {
			return err
		} else {
			return jpeg.Encode(wexif, img, opt.Options)
		}
	}
	return jpeg.Encode(w, img, opt.Options)
}

func (im ImageMore) WriteJPEGFile(filename string, opt *JPEGEncodeOptions) error {
	return writeJPEGFile(filename, im.Image, opt)
}

func writeJPEGFile(filename string, img image.Image, opt *JPEGEncodeOptions) error {
	if img == nil {
		return ErrImageNotSet
	} else if w, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600); err != nil {
		return err
	} else {
		defer w.Close()
		return writeJPEG(w, img, opt)
	}
}

func (im ImageMore) WriteJPEGResponseWriter(w http.ResponseWriter, addContentTypeHeader bool, opt *JPEGEncodeOptions) error {
	return writeJPEGResponseWriter(w, addContentTypeHeader, im.Image, opt)
}

func writeJPEGResponseWriter(w http.ResponseWriter, addContentTypeHeader bool, img image.Image, opt *JPEGEncodeOptions) error {
	if img == nil {
		return ErrImageNotSet
	} else if b, err := bytesJPEG(img, opt); err != nil {
		return err
	} else {
		if addContentTypeHeader {
			w.Header().Set(httputilmore.HeaderContentType, httputilmore.ContentTypeImageJPEG)
		}
		_, err = w.Write(b)
		return err
	}
}

func (im ImageMore) WritePNG(w io.Writer) error {
	return writePNG(w, im.Image)
}

func writePNG(w io.Writer, img image.Image) error {
	if img == nil {
		return ErrImageNotSet
	}
	return png.Encode(w, img)
}

func (im ImageMore) WritePNGFile(filename string) error {
	return writePNGFile(filename, im.Image)
}

func writePNGFile(filename string, img image.Image) error {
	if img == nil {
		return ErrImageNotSet
	} else if w, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600); err != nil {
		return err
	} else if err = png.Encode(w, img); err != nil {
		return err
	} else {
		return w.Close()
	}
}

func (im ImageMore) WritePNGResponseWriter(w http.ResponseWriter, addContentTypeHeader bool) error {
	return writePNGResponseWriter(w, addContentTypeHeader, im.Image)
}

func writePNGResponseWriter(w http.ResponseWriter, addContentTypeHeader bool, img image.Image) error {
	if img == nil {
		return ErrImageNotSet
	} else if b, err := bytesPNG(img); err != nil {
		return err
	} else {
		if addContentTypeHeader {
			w.Header().Set(httputilmore.HeaderContentType, httputilmore.ContentTypeImagePNG)
		}
		_, err = w.Write(b)
		return err
	}
}
