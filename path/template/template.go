package template

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/grokify/mogo/reflect/reflectutil"
	"github.com/grokify/mogo/time/timeutil"
	"github.com/grokify/mogo/type/stringsutil"
)

// Execute takes a template string and an `any`
// struct, substituting struct values for the variables.
// Field names can be nested.
func Execute(pattern string, item any) string {
	r := regexp.MustCompile(`{{(.*?)}}`)
	m := r.FindAllStringSubmatch(pattern, -1)

	if len(m) == 0 {
		return pattern
	}

	fields := []FieldInfo{}
	for _, mi := range m {
		fields = append(fields, ParseFieldInfoString(mi[1])...)
	}

	return fmt.Sprintf(
		r.ReplaceAllString(pattern, "%v"),
		GetFieldsFormattedForce(item, fields)...)
}

type FieldInfo struct {
	Name    string // . separated name path
	Type    string // . separated type path
	Options []string
}

func ParseFieldInfoString(fieldsRawStr string) []FieldInfo {
	fields := []FieldInfo{}
	fieldsRaw := strings.Split(fieldsRawStr, ";")
	for _, fieldStr := range fieldsRaw {
		f := FieldInfo{}
		fieldSections := strings.Split(fieldStr, "|")
		if len(fieldSections) > 2 {
			opts := strings.Split(fieldSections[2], ",")
			if len(opts) > 0 {
				f.Options = opts
			}
		}
		if len(fieldSections) > 1 {
			f.Type = strings.ToLower(strings.TrimSpace(fieldSections[1]))
		}
		if len(fieldSections) > 0 {
			f.Name = fieldSections[0]
			fields = append(fields, f)
		}
	}
	return fields
}

func GetFieldFormatted(item any, f FieldInfo) (any, error) {
	val, err := reflectutil.GetField(item, strings.Split(strings.TrimSpace(f.Name), ".")...)
	if err != nil {
		return val, err
	}

	typePath := strings.Split(strings.TrimSpace(f.Type), ".")
	if len(typePath) > 0 {
		switch strings.TrimSpace(typePath[0]) {
		case "string":
			return stringsutil.FormatString(val.(string), f.Options), nil
		case "time":
			dt := val.(time.Time)
			format := time.RFC3339
			if len(f.Options) > 0 {
				formatTry := f.Options[0]
				formatFound, err := timeutil.GetFormat(formatTry)
				if err != nil {
					return val, err
				}
				format = formatFound
			}
			return dt.Format(format), nil
		}
	}
	return val, nil
}

// GetFieldsFormatted returns an `any` slice for the struct and fields requested.
// An error is returned if any fields are not found or parsing options fail.
func GetFieldsFormatted(item any, fs []FieldInfo) ([]any, error) {
	vals := []any{}
	for _, f := range fs {
		val, err := GetFieldFormatted(item, f)
		if err != nil {
			return vals, err
		}
		vals = append(vals, val)
	}
	return vals, nil
}

// GetFieldsFormattedForce returns an `any` slice for the struct and
// fields requested. An empty string value is returned for anything that
// encounters an error.
func GetFieldsFormattedForce(item any, fs []FieldInfo) []any {
	vals := []any{}
	for _, f := range fs {
		val, err := GetFieldFormatted(item, f)
		if err != nil {
			vals = append(vals, "")
		}
		vals = append(vals, val)
	}
	return vals
}
