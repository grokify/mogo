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
