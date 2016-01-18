package fmtutil

import (
	"encoding/json"
	"fmt"
)

// PrintJSON pretty prints anything using a default indentation
func PrintJSON(in interface{}) error {
	j, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", string(j))
	return nil
}
