package api

import (
	"strings"
)

const (
	TypeMethod = "method"
	TypeEvent  = "event"
)

func ParseOperationType(input, def string) string {
	inputLc := strings.ToLower(strings.TrimSpace(input))
	if inputLc == TypeMethod {
		return TypeMethod
	} else if inputLc == TypeEvent {
		return TypeEvent
	}
	return def
}

type Operation struct {
	Tags        []string
	Class       string // Object class for SDKs
	Type        string // Method or Event
	Path        string
	Summary     string
	Description string
	Link        string
}
