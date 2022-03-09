// logutil provides logging utility functions.
package logutil

import "log"

func FatalOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func PrintOnError(err error) {
	if err != nil {
		log.Print(err.Error())
	}
}

func PrintlnOnError(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}
