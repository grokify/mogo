package reflectutil

import (
	"errors"
	"reflect"
	// reflections "gopkg.in/oleiade/reflections.v1"
)

func GetString(i any, key string) string {
	immutable := reflect.ValueOf(i)
	return immutable.FieldByName(key).String()
}

func Set(i any, key string, value any) {
	field := reflect.ValueOf(i).Elem().FieldByName(key)
	field.Set(reflect.ValueOf(value))
}

/*
func GetField(i any, fieldPath ...string) (any, error) {
	if len(fieldPath) == 0 {
		return i, nil
	}
	nextItem, err := reflections.GetField(i, strings.TrimSpace(fieldPath[0]))
	if err != nil || len(fieldPath) == 1 {
		return nextItem, err
	}
	return GetField(nextItem, fieldPath[1:]...)
}
*/

var ErrFieldNotFound = errors.New("field not found")

// FieldTagValue returns a tag name. For example, in 'Attribute string: `json:"attribute,omitempty"`', the
// usage would be `FieldTagValue(s, "Attribute", "json")` which would return `attribute`.
func FieldTagValue(a any, fieldName, tagName string) (string, error) {
	val := reflect.ValueOf(a)
	if field, ok := val.Type().FieldByName(fieldName); !ok {
		return "", ErrFieldNotFound
	} else {
		return field.Tag.Get(tagName), nil
	}
}

// NameOf returns the name of a struct. If `inclPkgPath` is set to `true`, a
// fully-qualified name is returned including package path.
func NameOf(i any, inclPkgPath bool) string {
	var ptr, name string
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		ptr = "*"
	}
	name = t.Name()
	if inclPkgPath {
		pkgPath := t.PkgPath()
		if len(pkgPath) > 0 {
			name = pkgPath + "." + name
		}
	}
	return ptr + name
}

// SliceInterfaceToString converts an `any` to a `[]string`.
func SliceInterfaceToString(raws any) []string {
	out := []string{}
	switch reflect.TypeOf(raws).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(raws)

		for i := 0; i < s.Len(); i++ {
			val := s.Index(i)
			out = append(out, val.Interface().(string))
		}
	}
	return out
}

func IsNil(i any) bool {
	// From https://medium.com/@mangatmodi/go-check-nil-interface-the-right-way-d142776edef1
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
