package reflectutil

import (
	"reflect"
	"strings"

	reflections "gopkg.in/oleiade/reflections.v1"
)

func GetString(i any, key string) string {
	immutable := reflect.ValueOf(i)
	return immutable.FieldByName(key).String()
}

func Set(i any, key string, value any) {
	field := reflect.ValueOf(i).Elem().FieldByName(key)
	field.Set(reflect.ValueOf(value))
}

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

// TypeName returns the name of a struct. stackoverflow-answerId:1908967
func TypeName(i any) string {
	res := ""
	t := reflect.TypeOf(i)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		res += "*"
	}
	return res + t.Name()
}

// SliceInterfaceToString converts an `interface{}` to a `[]string`.
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
