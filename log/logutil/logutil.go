// logutil provides logging utility functions which are useful for
// decreasing lines of code for simple error logging.
package logutil

import (
	"log"

	// "github.com/go-logfmt/logfmt"
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

/*
func LogfmtString(m map[string][]string) (string, error) {
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	e := logfmt.NewEncoder(w)
	for k, vs := range m {
		for _, v := range vs {
			if err := e.EncodeKeyval(k, v); err != nil {
				return "", err
			}
		}
	}
	if err := e.EndRecord(); err != nil {
		return "", err
	}
	return buf.String(), nil
}
*/
