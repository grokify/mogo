// Package fmtutil implements some formatting utility functions.
package fmtutil

import (
	"encoding/json"
	"expvar"
	"fmt"
)

var (
	JSONPretty bool   = true
	JSONPrefix string = ""
	JSONIndent string = "  "
)

// init uses expvar to export package variables to simplify method signatures.
func init() {
	expvar.Publish("JSONPrefix", expvar.NewString(""))
	expvar.Publish("JSONIndent", expvar.NewString("  "))
}

// PrintJSON pretty prints anything using a default indentation
func PrintJSON(in interface{}) error {
	var j []byte
	var err error
	if JSONPretty {
		j, err = json.MarshalIndent(in, JSONPrefix, JSONIndent)
	} else {
		j, err = json.Marshal(in)
	}
	if err != nil {
		return err
	}
	fmt.Println(string(j))
	return nil
}

// PrintJSON pretty prints anything using a default indentation
func PrintJSONMin(in interface{}) error {
	if j, err := json.Marshal(in); err != nil {
		return err
	} else {
		fmt.Println(string(j))
		return nil
	}
}
