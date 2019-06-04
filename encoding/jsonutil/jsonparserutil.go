package jsonutil

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/buger/jsonparser"
)

func JsonParserGetArrayIntOneOnly(data []byte, key string) (int, error) {
	values, err := JsonParserGetArrayString(data, key)
	if err != nil {
		return -1, err
	}
	if len(values) != 1 {
		return -1, fmt.Errorf("Expecting 1 Value, Got [%v]", len(values))
	}
	valueString := values[0]
	valueInt, err := strconv.Atoi(valueString)
	return valueInt, err
}

func JsonParserGetArrayStringOneOnly(data []byte, key string) (string, error) {
	values, err := JsonParserGetArrayString(data, key)
	if err != nil {
		return "", err
	}
	if len(values) != 1 {
		return "", fmt.Errorf("Expecting 1 Value, Got [%v]", len(values))
	}
	return values[0], nil
}

func JsonParserGetArrayString(data []byte, key string) ([]string, error) {
	strs := []string{}
	value, _, _, err := jsonparser.Get(data, key)
	if err != nil {
		return strs, err
	}
	err = json.Unmarshal(value, &strs)
	return strs, err
}
