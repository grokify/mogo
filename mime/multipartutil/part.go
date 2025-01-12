package multipartutil

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"

	hum "github.com/grokify/mogo/net/http/httputilmore"
)

const (
	PartTypeFilepath = "filepath"
	PartTypeJSON     = "json"
	PartTypeRaw      = "raw"
)

type Part struct {
	Type               string
	Name               string
	ContentType        string
	ContentDisposition string // short
	HeaderRaw          textproto.MIMEHeader
	BodyDataFilepath   string
	BodyDataJSON       any
	BodyDataRaw        []byte
	BodyEncodeBase64   bool
}

// HeaderBodyFilepath sets Content-Disposition and Content-Type.
func (p Part) HeaderBodyFilepath() (textproto.MIMEHeader, []byte, error) {
	header := textproto.MIMEHeader{}

	if cd := p.ContentDispositionHeader(); cd != "" {
		header.Add(hum.HeaderContentDisposition, cd)
	}

	_, filename := filepath.Split(p.BodyDataFilepath)
	filename = strings.TrimSpace(filename)
	mimeType := mime.TypeByExtension(filepath.Ext(p.BodyDataFilepath))
	if len(mimeType) > 0 {
		header.Add(hum.HeaderContentType, mimeType)
	}

	for k, vals := range p.HeaderRaw {
		for _, v := range vals {
			header.Add(k, v)
		}
	}

	body, err := os.ReadFile(p.BodyDataFilepath)
	if err != nil {
		return header, []byte{}, err
	}

	if p.BodyEncodeBase64 {
		header.Add(hum.HeaderContentTransferEncoding, "base64")
		body = []byte(base64.StdEncoding.EncodeToString(body))
	}

	return header, body, nil
}

func (p Part) FilepathToRaw() (Part, error) {
	if p.Type != PartTypeFilepath {
		return Part{}, errors.New("part is not filepath type")
	}
	header, body, err := p.HeaderBodyFilepath()
	if err != nil {
		return Part{}, err
	}
	return Part{
		Type:        PartTypeRaw,
		HeaderRaw:   header,
		BodyDataRaw: body,
	}, nil
}

func (p *Part) ContentDispositionHeader() string {
	disposition := strings.TrimSpace(p.ContentDisposition)
	params := map[string]string{}
	name := strings.TrimSpace(p.Name)
	if name != "" {
		params["name"] = name
	}
	_, filename := filepath.Split(p.BodyDataFilepath)
	if filename = strings.TrimSpace(filename); filename != "" {
		params["filename"] = filename
	}

	return mime.FormatMediaType(disposition, params)
}

// HeaderBodyJSON adds a JSON part.
func (p *Part) HeaderBodyJSON() (textproto.MIMEHeader, []byte, error) {
	header := textproto.MIMEHeader{}

	if cd := p.ContentDispositionHeader(); cd != "" {
		header.Add(hum.HeaderContentDisposition, cd)
	}

	header.Add(hum.HeaderContentType, hum.ContentTypeAppJSONUtf8)

	body, err := json.Marshal(p.BodyDataJSON)
	if err != nil {
		return header, []byte{}, err
	}

	if p.BodyEncodeBase64 {
		header.Add(hum.HeaderContentTransferEncoding, "base64")
		body = []byte(base64.StdEncoding.EncodeToString(body))
	}

	return header, body, nil
}

func (p Part) HeaderBodyRaw() (textproto.MIMEHeader, []byte, error) {
	header := textproto.MIMEHeader{}
	if strings.TrimSpace(p.ContentType) != "" {
		header.Add(hum.HeaderContentType, p.ContentType)
	}
	body := p.BodyDataRaw
	if p.BodyEncodeBase64 {
		header.Add(hum.HeaderContentTransferEncoding, "base64")
		body = []byte(base64.StdEncoding.EncodeToString(body))
	}
	for k, vals := range p.HeaderRaw {
		for _, v := range vals {
			header.Add(k, v)
		}
	}
	return header, body, nil
}

func (p Part) Write(w *multipart.Writer) error {
	header := textproto.MIMEHeader{}
	var body []byte
	var err error
	switch p.Type {
	case PartTypeFilepath:
		header, body, err = p.HeaderBodyFilepath()
	case PartTypeJSON:
		header, body, err = p.HeaderBodyJSON()
	case PartTypeRaw:
		header, body, err = p.HeaderBodyRaw()
	default:
		return fmt.Errorf("part type not supported (%s)", p.Type)
	}
	if err != nil {
		return err
	} else if partWriter, err := w.CreatePart(header); err != nil {
		return err
	} else {
		_, err = bytes.NewBuffer(body).WriteTo(partWriter)
		return err
	}
}
