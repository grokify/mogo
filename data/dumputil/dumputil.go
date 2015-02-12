package dumputil

import (
	"encoding/json"
	"fmt"
)

func DumpJson(in interface{}) {
	j, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		jerr, _ := json.MarshalIndent(err, "", "  ")
		fmt.Printf("%+v\n", string(jerr))
	} else {
		fmt.Printf("%+v\n", string(j))
	}
}
