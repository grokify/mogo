// Package fmtutil implements some formatting utility functions.
package fmtutil

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/grokify/mogo/encoding/jsonutil"
)

var (
	JSONPretty bool   = true
	JSONPrefix string = ""
	JSONIndent string = "  "
)

// init uses expvar to export package variables to simplify method signatures.
/*
func init() {
	expvar.Publish("JSONPrefix", expvar.NewString(""))
	expvar.Publish("JSONIndent", expvar.NewString("  "))
}
*/

// PrintJSON pretty prints anything using a default indentation
func PrintJSON(in any) error {
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

func MustPrintJSON(in any) {
	if err := PrintJSON(in); err != nil {
		panic(err)
	}
}

// PrintJSONMore pretty prints anything using supplied indentation.
func PrintJSONMore(in any, jsonPrefix, jsonIndent string) error {
	j, err := jsonutil.MarshalSimple(in, jsonPrefix, jsonIndent)
	if err != nil {
		return err
	}
	fmt.Println(string(j))
	return nil
}

// PrintJSONMin pretty prints anything using a default indentation
func PrintJSONMin(in any) error {
	if j, err := json.Marshal(in); err != nil {
		return err
	} else {
		fmt.Println(string(j))
		return nil
	}
}

func PrintReader(r io.Reader) error {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))
	return nil
}
