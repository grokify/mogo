package fmtutil

import (
	"encoding/json"
	"fmt"
)

func PrintJson(in interface{}) error {
	j, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", string(j))
	return nil
}
