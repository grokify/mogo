package jsonraw

import (
	"encoding/json"
	"io"
	"os"
	"reflect"
)

func Equal(x, y io.Reader) (bool, error) {
	var ax, ay any
	d := json.NewDecoder(x)
	if err := d.Decode(&ax); err != nil {
		return false, err
	}
	d = json.NewDecoder(y)
	if err := d.Decode(&ay); err != nil {
		return false, err
	}
	return reflect.DeepEqual(ax, ay), nil
}

func EqualBytes(x, y []byte) (bool, error) {
	var ax, ay any
	if err := json.Unmarshal(x, &ax); err != nil {
		return false, err
	} else if err := json.Unmarshal(y, &ay); err != nil {
		return false, err
	} else {
		return reflect.DeepEqual(ax, ay), nil
	}
}

func EqualFiles(x, y string) (bool, error) {
	fx, err := os.Open(x)
	if err != nil {
		return false, err
	}
	defer fx.Close()
	fy, err := os.Open(y)
	if err != nil {
		return false, err
	}
	defer fy.Close()
	return Equal(fx, fy)
}
