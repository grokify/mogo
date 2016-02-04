// Package fmtutil implements some formatting utility functions.
package fmtutil

import (
	"encoding/json"
	"expvar"
	"fmt"
)

var (
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
	j, err := json.MarshalIndent(in, JSONPrefix, JSONIndent)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", string(j))
	return nil
}
