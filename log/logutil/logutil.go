// logutil provides logging utility functions which are useful for
// decreasing lines of code for simple error logging.
package logutil

import (
	"log"

	"github.com/grokify/mogo/errors/errorsutil"
)

func FatalErr(err error, wrap ...string) {
	if err != nil {
		for i := len(wrap) - 1; i >= 0; i-- {
			err = errorsutil.Wrap(err, wrap[i])
		}
		log.Fatal(err.Error())
	}
}

func PrintErr(err error, wrap ...string) {
	if err != nil {
		for i := len(wrap) - 1; i >= 0; i-- {
			err = errorsutil.Wrap(err, wrap[i])
		}
		log.Print(err.Error())
	}
}

func PrintlnErr(err error, wrap ...string) {
	if err != nil {
		for i := len(wrap) - 1; i >= 0; i-- {
			err = errorsutil.Wrap(err, wrap[i])
		}
		log.Println(err.Error())
	}
}
