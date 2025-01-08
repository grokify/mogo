package mailutil

import (
	"mime"
	"net/http"
	"net/mail"
	"strings"
	"time"
)

type Header map[string][]string

type MessageWriter struct {
	To      Addresses
	Cc      Addresses
	Bcc     Addresses
	From    *mail.Address
	Sender  *mail.Address
	Date    *time.Time
	Subject string
	Header  http.Header
	Body    string
}

func (mw MessageWriter) HeaderLines() []string {
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
	return key + HeaderSep + EncodeHeaderValueUTF8(val), nil
}

func EncodeHeaderValueUTF8(s string) string {
	return mime.QEncoding.Encode("utf-8", s)
}

func (mw MessageWriter) String() string {
	lines := mw.HeaderLines()
	header := strings.Join(lines, "\n")
	var parts []string
	if header != "" {
		parts = append(parts, header)
	}
	if mw.Body != "" {
		parts = append(parts, mw.Body)
	}
	return strings.Join(parts, "\n\n")
}
