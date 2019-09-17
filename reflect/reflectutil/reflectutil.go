package reflectutil

import (
	"reflect"
	"strings"

	"gopkg.in/oleiade/reflections.v1"
)

func GetString(i interface{}, key string) string {
	immutable := reflect.ValueOf(i)
	return immutable.FieldByName(key).String()
}

func Set(i interface{}, key string, value interface{}) {
	field := reflect.ValueOf(i).Elem().FieldByName(key)
	field.Set(reflect.ValueOf(value))
}

func GetField(item interface{}, fieldPath ...string) (interface{}, error) {
	if fieldPath == nil || len(fieldPath) == 0 {
		return item, nil
	}
	nextItem, err := reflections.GetField(item, strings.TrimSpace(fieldPath[0]))
	if err != nil || len(fieldPath) == 1 {
		return nextItem, err
	}
	return GetField(nextItem, fieldPath[1:]...)
}

// TypeName returns the name of a struct.
// stackoverflow-answerId:1908967
func TypeName(myvar interface{}) (res string) {
	t := reflect.TypeOf(myvar)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		res += "*"
	}
	return res + t.Name()
}
