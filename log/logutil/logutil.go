// logutil provides logging utility functions which are useful for
// decreasing lines of code for simple error logging.
package logutil

import (
	"log"

	"github.com/grokify/mogo/errors/errorsutil"
)

func FatalErr(err error, wrap ...string) {
	if err != nil {
		for _, w := range wrap {
			err = errorsutil.Wrap(err, w)
		}
		log.Fatal(err.Error())
	}
}

func PrintErr(err error) {
	if err != nil {
		log.Print(err.Error())
	}
}

func PrintlnErr(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}
