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

func GetField(item interface{}, namePath ...string) (interface{}, error) {
	if namePath == nil || len(namePath) == 0 {
		return item, nil
	}
	nextItem, err := reflections.GetField(item, strings.TrimSpace(namePath[0]))
	if err != nil || len(namePath) == 1 {
		return nextItem, err
	}
	return GetField(nextItem, namePath[1:]...)
}
