package jsonutil

import (
	"github.com/valyala/fastjson"
)

func GetSubobjectBytes(data []byte, key string) ([]byte, error) {
	val, err := fastjson.ParseBytes(data)
	if err != nil {
		return []byte{}, err
	}
	return val.GetObject(key).MarshalTo([]byte{}), nil
}

func MustGetSubobjectBytes(data []byte, key string) []byte {
	return fastjson.MustParseBytes(data).GetObject(key).MarshalTo([]byte{})
}
