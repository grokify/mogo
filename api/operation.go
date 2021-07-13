package api

import (
	"strings"

	"github.com/grokify/simplego/type/stringsutil"
)

const (
	TypeAPI               = "api"
	TypeMethod            = "method"
	TypeEvent             = "event"
	TypeEventAliasWebhook = "webhook"
)

type Operations []Operation

func (ops Operations) Table() ([]string, [][]string) {
	cols := []string{
		"Tags", "Class", "Type", "Path", "Summary", "Description", "Link"}
	rows := [][]string{}
	for _, op := range ops {
		row := []string{
			strings.Join(op.Tags, ", "),
			op.Class,
			op.Type,
			op.Path,
			op.Summary,
			op.Description,
			op.Link}
		rows = append(rows, row)
	}
	return cols, rows
}

func ParseOperationType(input, def string) string {
	inputLc := strings.ToLower(strings.TrimSpace(input))
	switch inputLc {
	case TypeMethod:
		return TypeMethod
	case TypeAPI:
		return TypeAPI
	case TypeEvent:
		return TypeEvent
	case TypeEventAliasWebhook:
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

func (op *Operation) TrimSpace() {
	op.Tags = stringsutil.SliceCondenseSpace(op.Tags, true, false)
	op.Class = strings.TrimSpace(op.Class)
	op.Type = strings.TrimSpace(op.Type)
	op.Path = strings.TrimSpace(op.Path)
	op.Summary = strings.TrimSpace(op.Summary)
	op.Description = strings.TrimSpace(op.Description)
	op.Link = strings.TrimSpace(op.Link)
}
