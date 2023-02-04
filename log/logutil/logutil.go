// logutil provides logging utility functions which are useful for
// decreasing lines of code for simple error logging.
package logutil

import (
	"bufio"
	"bytes"
	"log"

	"github.com/go-logfmt/logfmt"
	"github.com/grokify/mogo/errors/errorsutil"
)

func FatalErr(err error, wrap ...string) {
	err = errorsutil.Wrap(err, wrap...)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func PrintErr(err error, wrap ...string) {
	err = errorsutil.Wrap(err, wrap...)
	if err != nil {
		log.Print(err.Error())
	}
}

func PrintlnErr(err error, wrap ...string) {
	err = errorsutil.Wrap(err, wrap...)
	if err != nil {
		log.Println(err.Error())
	}
}

func LogfmtString(m map[string][]string) (string, error) {
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	e := logfmt.NewEncoder(w)
	for k, vs := range m {
		for _, v := range vs {
			err := e.EncodeKeyval(k, v)
			if err != nil {
				return "", err
			}
		}
	}
	err := e.EndRecord()
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
