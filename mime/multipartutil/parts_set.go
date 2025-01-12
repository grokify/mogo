package multipartutil

type PartsSet struct {
	ContentType string
	Parts       []Part
}

func NewPartsSet(ct string) PartsSet {
	return PartsSet{
		ContentType: ct,
		Parts:       []Part{}}
}

func (ps *PartsSet) AddMailBody(textBody, htmlBody []byte) error {
	if p, err := NewPartAlternativeOrNot(textBody, htmlBody); err != nil {
		return err
	} else {
		ps.Parts = append([]Part{p}, ps.Parts...)
		return nil
	}
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
