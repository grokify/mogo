package reflectutil

import (
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/grokify/gotilla/time/timeutil"
	"gopkg.in/oleiade/reflections.v1"
)

/*
 * http://golang-examples.tumblr.com/post/44089080167/get-set-a-field-value-of-a-struct-using-reflection
 */

func GetString(i interface{}, key string) string {
	immutable := reflect.ValueOf(i)
	return immutable.FieldByName(key).String()
}

func Set(i interface{}, key string, value interface{}) {
	field := reflect.ValueOf(i).Elem().FieldByName(key)
	field.Set(reflect.ValueOf(value))
}

type FieldInfo struct {
	Name    string // . separated name path
	Type    string // . separated type path
	Options []string
}

func SplitFieldInfoString(fieldsRawStr string) []FieldInfo {
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

func GetFieldRecurse(item interface{}, namePath []string) (interface{}, error) {
	if len(namePath) == 0 {
		return item, nil
	}
	nextItem, err := reflections.GetField(item, strings.TrimSpace(namePath[0]))
	if err != nil || len(namePath) == 1 {
		return nextItem, err
	}
	return GetFieldRecurse(nextItem, namePath[1:])
}

func FormatString(s string, options []string) string {
	for _, opt := range options {
		switch strings.TrimSpace(opt) {
		case "StringToLower":
			s = strings.ToLower(s)
		case "SpaceToHyphen":
			s = regexp.MustCompile(`[\s-]+`).ReplaceAllString(s, "-")
		case "SpaceToUnderscore":
			s = regexp.MustCompile(`[\s_]+`).ReplaceAllString(s, "_")
		}
	}
	return s
}

func GetFieldFormatted(item interface{}, f FieldInfo) (interface{}, error) {
	var wip interface{}
	var err error

	val, err := GetFieldRecurse(
		item,
		strings.Split(strings.TrimSpace(f.Name), "."))

	if err != nil {
		return wip, err
	}

	typePath := strings.Split(strings.TrimSpace(f.Type), ".")
	if len(typePath) > 0 {
		switch strings.TrimSpace(typePath[0]) {
		case "string":
			return FormatString(val.(string), f.Options), nil
		case "time":
			dt := val.(time.Time)
			timeString := ""

			format := time.RFC3339
			if len(f.Options) > 0 {
				formatTry := f.Options[0]
				formatFound, err := timeutil.GetFormat(formatTry)
				if err != nil {
					return wip, err
				}
				format = formatFound
			}
			timeString = dt.Format(format)
			return timeString, nil
		}
	}
	return wip, nil
}

// GetFieldsFormatted returns an interface{} slice for the struct and fields
// requested. An error is returned if any fields are not found or parsing options
// fail.
func GetFieldsFormatted(item interface{}, fs []FieldInfo) ([]interface{}, error) {
	vals := []interface{}{}
	for _, f := range fs {
		val, err := GetFieldFormatted(item, f)
		if err != nil {
			return vals, err
		}
		vals = append(vals, val)
	}
	return vals, nil
}

// GetFieldsFormattedForce returns an interface{} slice for the struct and
// fields requested. An empty string value is returned for anything that
// encounters an error.
func GetFieldsFormattedForce(item interface{}, fs []FieldInfo) []interface{} {
	vals := []interface{}{}
	for _, f := range fs {
		val, err := GetFieldFormatted(item, f)
		if err != nil {
			vals = append(vals, "")
		}
		vals = append(vals, val)
	}
	return vals
}
