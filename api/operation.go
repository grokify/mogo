package api

import (
	"strings"

	"github.com/grokify/mogo/type/maputil"
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

func (ops Operations) TrimSpace(sortTags bool) {
	for i, op := range ops {
		op.TrimSpace(sortTags)
		ops[i] = op
	}
}

func (ops Operations) RemovePaths(paths []string) Operations {
	paths = stringsutil.SliceDedupe(paths)
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

func (ops Operations) CountsByTag(concatenate bool, sep string) map[string]int {
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

func (ops Operations) DuplicatePaths() map[string]int {
	mapPaths := map[string]int{}
	for _, op := range ops {
		mapPaths[op.Path] += 1
	}
	dupPaths := map[string]int{}
	for opPath, pathCount := range mapPaths {
		if pathCount > 1 {
			dupPaths[opPath] = pathCount
		}
	}
	return dupPaths
}

func (ops Operations) CountsByTagCount() map[int]int {
	mii := map[int]int{}
	for _, op := range ops {
		mii[len(op.Tags)] += 1
	}
	return mii
}

func (ops Operations) CountsByClass() map[string]int {
	msi := map[string]int{}
	for _, op := range ops {
		msi[op.Class] += 1
	}
	return msi
}

func (ops Operations) CountsByType() (operationCounts map[string]int, methodCounts map[string]int) {
	operationCounts = map[string]int{}
	methodCounts = map[string]int{}
	for _, op := range ops {
		operationCounts[op.OperationType] += 1
		if op.OperationType == TypeMethod {
			methodCounts[op.MethodType] += 1
		}
	}
	return
}

func (ops Operations) Table() ([]string, [][]string) {
	cols := []string{
		"Tags", "Class", "Op Type", "Method Type", "Path", "Summary", "Description", "Link"}
	rows := [][]string{}
	for _, op := range ops {
		rows = append(rows, []string{
			strings.Join(op.Tags, ", "),
			op.Class,
			op.OperationType,
			op.MethodType,
			op.Path,
			op.Summary,
			op.Description,
			op.Link})
	}
	return cols, rows
}

func (ops Operations) PathsByClass() maputil.MapStringSlice {
	mss := maputil.MapStringSlice{}
	for _, op := range ops {
		mss.Add(op.Class, op.Path)
	}
	mss.Sort(true)
	return mss
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

func (op *Operation) TrimSpace(sortTags bool) {
	op.Tags = stringsutil.SliceCondenseSpace(op.Tags, true, sortTags)
	op.Class = strings.TrimSpace(op.Class)
	op.OperationType = strings.TrimSpace(op.OperationType)
	op.MethodType = strings.TrimSpace(op.MethodType)
	op.Path = strings.TrimSpace(op.Path)
	op.Summary = strings.TrimSpace(op.Summary)
	op.Description = strings.TrimSpace(op.Description)
	op.Link = strings.TrimSpace(op.Link)
}
