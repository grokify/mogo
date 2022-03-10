// logutil provides logging utility functions which are useful for
// decreasing lines of code for simple error logging.
package logutil

import "log"

func FatalErr(err error) {
	if err != nil {
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
