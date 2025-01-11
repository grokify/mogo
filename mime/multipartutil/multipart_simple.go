package multipartutil

type MultipartSimple struct {
	ContentType string
	Parts       []Part
}

func NewMultipartSimple() MultipartSimple {
	return MultipartSimple{Parts: []Part{}}
}

func (ms MultipartSimple) Builder(close bool) (MultipartBuilder, error) {
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
func (ms MultipartSimple) Part() (Part, error) {
	ct, body, err := ms.Strings()
	if err != nil {
		return Part{}, err
	} else {
		return Part{
			Type:         PartTypeRaw,
			ContentType:  ct,
			Base64Encode: false,
			RawBody:      []byte(body),
		}, nil
	}
}

func (ms MultipartSimple) Strings() (ctHeader, body string, err error) {
	mb, err := ms.Builder(true)
	if err != nil {
		return
	}
	ctHeader = mb.ContentType(ms.ContentType)
	body = mb.String()
	return
}
