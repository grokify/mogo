package main

import (
	"fmt"
	"log"

	"github.com/grokify/mogo/crypto/token"
)

func main() {
	tok, err := token.Base58(32)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("token: %s\n", tok)

	fmt.Println("DONE")
}
