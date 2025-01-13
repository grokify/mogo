package multipartutil

import (
	"slices"

	"github.com/grokify/mogo/net/http/httputilmore"
)

type PartsSet struct {
	ContentType string
	Parts       Parts
}

func NewPartsSet(ct string) PartsSet {
	return PartsSet{
		ContentType: ct,
		Parts:       []Part{}}
}

// NewPartsSetMail returns a single part that represents an email message.
func NewPartsSetMail(textBody, htmlBody []byte, additionalParts Parts) (PartsSet, error) {
	if len(additionalParts) > 0 {
		ps := NewPartsSet(httputilmore.ContentTypeMultipartMixed)
		ps.Parts = slices.Clone(additionalParts)
		err := ps.AddMailBody(textBody, htmlBody)
		return ps, err
	} else {
		return NewPartsSetAlternative(textBody, htmlBody), nil
	}
}

func (ps *PartsSet) AddMailBody(textBody, htmlBody []byte) error {
	if p, err := NewPartAlternativeOrNot(textBody, htmlBody); err != nil {
		return err
	} else {
		ps.Parts = append([]Part{p}, ps.Parts...)
		return nil
	}
}

func (ps *PartsSet) Clone() PartsSet {
	new := NewPartsSet(ps.ContentType)
	new.Parts = slices.Clone(ps.Parts)
	return new
}

func (ms *PartsSet) Builder(close bool) (MultipartBuilder, error) {
	mb := NewMultipartBuilder()
	if len(ms.Parts) == 0 {
		err := mb.Close()
		return mb, err
	}
	for _, p := range ms.Parts {
		if err := p.Write(mb.Writer); err != nil {
			return mb, err
		}
	}
	if close {
		if err := mb.Close(); err != nil {
			return mb, err
		} else {
			return mb, nil
		}
	} else {
		return mb, nil
	}
}

// Part returns the MultipartSimple as a Part. This can be used for
// creating parts such as `multipart/alternative`.
func (ms *PartsSet) Part() (Part, error) {
	ct, body, err := ms.Strings()
	if err != nil {
		return Part{}, err
	} else {
		return Part{
			Type:             PartTypeRaw,
			ContentType:      ct,
			BodyEncodeBase64: false,
			BodyDataRaw:      []byte(body),
		}, nil
	}
}

func (ps *PartsSet) Strings() (ctHeader, body string, err error) {
	if len(ps.Parts) == 0 {
		return ps.ContentType, "", nil
	}

	mb, err := ps.Builder(true)
	if err != nil {
		return
	}
	ctHeader = mb.ContentType(ps.ContentType)
	body = mb.String()
	return
}
