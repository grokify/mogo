package api

import (
	"strings"

	"github.com/grokify/mogo/type/stringsutil"
)

const (
	TypeAPI               = "api"
	TypeMethod            = "method"
	TypeCallbackMethod    = "callback"
	TypeEvent             = "event"
	TypeEventAliasWebhook = "webhook"
	TypeCategory          = "category" // mobile
	TypeClass             = "class"    // mobile
	TypeConstant          = "constant" // mobile
	TypeProtocol          = "protocol" // mobile

	MethodClass    = "class"
	MethodInstance = "instance"
)

type Operations []Operation

func (ops Operations) TrimSpace() {
	for i, op := range ops {
		op.TrimSpace()
		ops[i] = op
	}
}

func (ops Operations) RemovePaths(paths []string) Operations {
	paths = stringsutil.SliceUnique(paths)
	if len(paths) == 0 {
		return ops
	}
	newOps := Operations{}
OP:
	for _, op := range ops {
		for _, rem := range paths {
			if op.Path == rem {
				continue OP
			}
		}
		newOps = append(newOps, op)
	}
	return newOps
}

func (ops Operations) TagCounts(concatenate bool, sep string) map[string]int {
	msi := map[string]int{}
	for _, op := range ops {
		tags := stringsutil.SliceCondenseSpace(op.Tags, true, true)
		if concatenate {
			tagc := strings.Join(tags, sep)
			msi[tagc] += 1
			continue
		}
		for _, tag := range tags {
			msi[tag] += 1
		}
	}
	return msi
}

func (ops Operations) PathsUnique() bool {
	mapPaths := map[string]int{}
	for _, op := range ops {
		mapPaths[op.Path] += 1
	}
	for _, v := range mapPaths {
		if v > 1 {
			return false
		}
	}
	return true
}

func (ops Operations) Table() ([]string, [][]string) {
	cols := []string{
		"Tags", "Class", "Op Type", "Method Type", "Path", "Summary", "Description", "Link"}
	rows := [][]string{}
	for _, op := range ops {
		row := []string{
			strings.Join(op.Tags, ", "),
			op.Class,
			op.OperationType,
			op.MethodType,
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
	Tags          []string
	Class         string // Object class for SDKs
	OperationType string // Method or Event
	MethodType    string
	Path          string
	Summary       string
	Description   string
	Link          string
}

func (op *Operation) TrimSpace() {
	op.Tags = stringsutil.SliceCondenseSpace(op.Tags, true, false)
	op.Class = strings.TrimSpace(op.Class)
	op.OperationType = strings.TrimSpace(op.OperationType)
	op.MethodType = strings.TrimSpace(op.MethodType)
	op.Path = strings.TrimSpace(op.Path)
	op.Summary = strings.TrimSpace(op.Summary)
	op.Description = strings.TrimSpace(op.Description)
	op.Link = strings.TrimSpace(op.Link)
}
