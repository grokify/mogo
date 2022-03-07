package jsonutil

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/buger/jsonparser"
)

func JSONParserGetArrayIntOneOnly(data []byte, key string) (int, error) {
	values, err := JSONParserGetArrayString(data, key)
	if err != nil {
		return -1, err
	}
	if len(values) != 1 {
		return -1, fmt.Errorf("expecting 1 value, got [%v] values", len(values))
	}
	valueString := values[0]
	valueInt, err := strconv.Atoi(valueString)
	return valueInt, err
}

func JSONParserGetArrayStringOneOnly(data []byte, key string) (string, error) {
	values, err := JSONParserGetArrayString(data, key)
	if err != nil {
		return "", err
	}
	if len(values) != 1 {
		return "", fmt.Errorf("expecting 1 value, got [%v] values", len(values))
	}
	return values[0], nil
}

func JSONParserGetArrayString(data []byte, key string) ([]string, error) {
	strs := []string{}
	value, _, _, err := jsonparser.Get(data, key)
	if err != nil {
		return strs, err
	}
	err = json.Unmarshal(value, &strs)
	return strs, err
}
