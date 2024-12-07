package multipartutil

import (
	"fmt"
	"strings"
)

type MultipartSimple struct {
	Parts []Part
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
		typ := strings.ToLower(strings.TrimSpace(p.Type))
		switch typ {
		case PartTypeFilepath:
			if err := mb.WriteFilePathPlus(p.Name, p.Filepath, p.Base64Encode); err != nil {
				return mb, nil
			}
		case PartTypeJSON:
			if err := mb.WriteFieldAsJSON(p.Name, p.Data, p.Base64Encode); err != nil {
				return mb, nil
			}
		case PartTypeString:
			if err := mb.WriteFieldString(p.Name, p.String); err != nil {
				return mb, nil
			}
		default:
			return mb, fmt.Errorf("type not supported (%s)", p.Type)
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

func (ms MultipartSimple) Strings() (ctHeader, body string, err error) {
	mb, err := ms.Builder(true)
	if err != nil {
		return
	}
	ctHeader = mb.ContentType()
	body = mb.String()
	return
}

const (
	PartTypeJSON     = "json"
	PartTypeFilepath = "filepath"
	PartTypeString   = "string"
)

type Part struct {
	Type         string
	Name         string
	Base64Encode bool
	Data         any
	Filepath     string
	String       string
}
