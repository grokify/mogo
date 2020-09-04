// qrutil provides wrapper functions for https://github.com/skip2/go-qrcode s.
package qrutil

import (
	"image"

	"github.com/skip2/go-qrcode"
)

type Opts struct {
	QRCode    qrcode.QRCode
	Filename  string
	Filesize  int
	FileWrite bool
}

// New creates a new `*qrcode.QRCode` using `qrcode.QRCode`
// as the request parameter.
func New(opts qrcode.QRCode) (*qrcode.QRCode, error) {
	q, err := qrcode.New(opts.Content, opts.Level)
	if err != nil {
		return nil, err
	}

	if opts.BackgroundColor != opts.ForegroundColor {
		q.BackgroundColor = opts.BackgroundColor
		q.ForegroundColor = opts.ForegroundColor
	}

	if opts.DisableBorder {
		q.DisableBorder = true
	}
	return q, nil
}

func NewImage(opts qrcode.QRCode, size int) (image.Image, error) {
	qr, err := New(opts)
	if err != nil {
		return nil, err
	}
	return qr.Image(size), nil
}

func Create(opts Opts) (*qrcode.QRCode, error) {
	q, err := New(opts.QRCode)
	if err != nil {
		return nil, err
	}

	if opts.FileWrite {
		err = q.WriteFile(opts.Filesize, opts.Filename)
	}

	return q, err
}
