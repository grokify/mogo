package reflectutil

import (
	"reflect"
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
