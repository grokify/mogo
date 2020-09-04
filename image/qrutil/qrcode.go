// qrutil provides wrapper functions for https://github.com/skip2/go-qrcode
package qrutil

import (
	"image"

	"github.com/skip2/go-qrcode"
)

// New creates a new `*qrcode.QRCode` using `qrcode.QRCode`
// as the request parameter. `Content` is qrcode value.
// `ForegroundColor` and `BackgroundColor` are `color.Color`.
// `Level` is in [`qrcode.Low`,`Medium`,`High`,`Highest`].
func New(opts qrcode.QRCode) (*qrcode.QRCode, error) {
	qrc, err := qrcode.New(opts.Content, opts.Level)
	if err != nil {
		return nil, err
	}

	if opts.BackgroundColor != opts.ForegroundColor {
		qrc.BackgroundColor = opts.BackgroundColor
		qrc.ForegroundColor = opts.ForegroundColor
	}

	if opts.DisableBorder {
		qrc.DisableBorder = true
	}
	return qrc, nil
}

func NewImage(opts qrcode.QRCode, size int) (image.Image, error) {
	qrc, err := New(opts)
	if err != nil {
		return nil, err
	}
	return qrc.Image(size), nil
}

func WritePNG(filename string, opts qrcode.QRCode, pixels int) error {
	qrc, err := New(opts)
	if err != nil {
		return err
	}

	return qrc.WriteFile(pixels, filename)
}
