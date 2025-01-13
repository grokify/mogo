package mailutil

import (
	"bytes"
	"io"
	"mime"
	"net/mail"
	"net/textproto"
	"strings"
	"time"

	"github.com/grokify/mogo/mime/multipartutil"
	"github.com/grokify/mogo/net/http/httputilmore"
)

type MessageWriter struct {
	To           Addresses
	Cc           Addresses
	Bcc          Addresses
	From         *mail.Address
	Sender       *mail.Address
	Date         *time.Time
	Subject      string
	Header       textproto.MIMEHeader
	BodyPartsSet multipartutil.PartsSet // Attachments or Inline
}

func NewMessageWriter() *MessageWriter {
	return &MessageWriter{
		Header:       textproto.MIMEHeader{},
		BodyPartsSet: multipartutil.NewPartsSet(httputilmore.ContentTypeMultipartMixed),
	}
}

func (mw *MessageWriter) RecipientCount() int {
	return len(mw.To.FilterInclWithAddress()) +
		len(mw.Cc.FilterInclWithAddress()) +
		len(mw.Bcc.FilterInclWithAddress())
}

func (mw *MessageWriter) HeaderLines() []string {
	var lines []string
	if mw.From != nil {
		if from := mw.From.String(); from != "" {
			lines = append(lines, MustEncodeHeaderLine(HeaderFrom, from))
		}
	}
	if to := mw.To.String(false, true, false); to != "" {
		lines = append(lines, MustEncodeHeaderLine(HeaderTo, to))
	}
	if cc := mw.Cc.String(false, true, false); cc != "" {
		lines = append(lines, MustEncodeHeaderLine(HeaderCc, cc))
	}
	if bcc := mw.Bcc.String(false, true, false); bcc != "" {
		lines = append(lines, MustEncodeHeaderLine(HeaderBcc, bcc))
	}
	if sub := strings.TrimSpace(mw.Subject); sub != "" {
		lines = append(lines, MustEncodeHeaderLine(HeaderSubject, sub))
	}
	for k, vals := range mw.Header {
		for _, v := range vals {
			lines = append(lines, MustEncodeHeaderLine(k, v))
		}
	}
	return lines
}

func MustEncodeHeaderLine(key, val string) string {
	if s, err := EncodeHeaderLine(key, val); err != nil {
		panic(err)
	} else {
		return s
	}
}

func EncodeHeaderLine(key, val string) (string, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		panic("no header key")
	}
	return key + HeaderSep + mime.QEncoding.Encode("utf-8", val), nil
}

func (mw *MessageWriter) Bytes() ([]byte, error) {
	var b bytes.Buffer
	if err := mw.Write(&b); err != nil {
		return []byte{}, err
	} else {
		return b.Bytes(), nil
	}
}

func (mw *MessageWriter) String() (string, error) {
	var b strings.Builder
	if err := mw.Write(&b); err != nil {
		return "", err
	} else {
		return b.String(), nil
	}
}

func (mw *MessageWriter) Write(w io.Writer) error {
	ctHeader, body, err := mw.BodyPartsSet.Strings()
	if err != nil {
		return err
	}

	if ctHeader = strings.TrimSpace(ctHeader); ctHeader != "" {
		if mw.Header == nil {
			mw.Header = textproto.MIMEHeader{}
		}
		mw.Header[httputilmore.HeaderContentType] = []string{ctHeader}
	}

	if count, err := w.Write([]byte(strings.Join(mw.HeaderLines(), "\n"))); err != nil {
		return err
	} else if count > 0 {
		if _, err := w.Write([]byte("\n\n")); err != nil {
			return err
		}
	}
	if len(body) > 0 {
		if _, err := w.Write([]byte(body)); err != nil {
			return err
		}
	}
	return nil
}
