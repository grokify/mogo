package fmtutil

import (
	"encoding/json"
	"fmt"
)

func PrintJson(in interface{}) {
	j, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		jerr, _ := json.MarshalIndent(err, "", "  ")
		fmt.Printf("%+v\n", string(jerr))
	} else {
		fmt.Printf("%+v\n", string(j))
	}
}
