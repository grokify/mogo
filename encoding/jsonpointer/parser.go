package jsonpointer

import (
	"errors"
	"strings"
)

var (
	ErrJSONPointerInvalidSyntaxNoAnchorSlash     = errors.New("invalid JSON Pointer format - no `#/`")
	ErrJSONPointerInvalidSyntaxNonOneAnchorSlash = errors.New("invalid JSON Pointer format - non-1 `#/`")
)

type JSONPointer struct {
	Document   string
	String     string
	PathString string
	Path       []string
}

func ParseJSONPointer(s string) (JSONPointer, error) {
	anchorSlash := "#/"
	ptr := JSONPointer{String: s}
	if strings.Index(s, anchorSlash) == 0 {
		ptr.PathString = s
		pathTrimmed := strings.TrimLeft(s, anchorSlash)
		ptr.Path = strings.Split(pathTrimmed, "/")
		return ptr, nil
	}
	if !strings.Contains(s, anchorSlash) {
		return ptr, ErrJSONPointerInvalidSyntaxNoAnchorSlash
	}
	parts := []string{}
	if strings.Contains(s, anchorSlash) {
		parts = strings.Split(s, anchorSlash)
	}
	if len(parts) != 2 {
		return ptr, ErrJSONPointerInvalidSyntaxNonOneAnchorSlash
	}
	ptr.Document = parts[0]
	ptr.PathString = parts[1]
	ptr.Path = strings.Split(ptr.PathString, "/")
	return ptr, nil
}
