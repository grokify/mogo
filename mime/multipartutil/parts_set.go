package multipartutil

type PartsSet struct {
	ContentType string
	Parts       []Part
}

func NewPartsSet() PartsSet {
	return PartsSet{Parts: []Part{}}
}

func (ms PartsSet) Builder(close bool) (MultipartBuilder, error) {
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
func (ms PartsSet) Part() (Part, error) {
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

func (ms PartsSet) Strings() (ctHeader, body string, err error) {
	mb, err := ms.Builder(true)
	if err != nil {
		return
	}
	ctHeader = mb.ContentType(ms.ContentType)
	body = mb.String()
	return
}
