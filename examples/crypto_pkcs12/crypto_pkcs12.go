package main

import (
	"fmt"
	"log"
	"os"

	p12 "github.com/grokify/simplego/crypto/pkcs12util"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(os.Getenv("ENV_PATH"))
	if err != nil {
		log.Fatal(err)
	}

	opts := p12.Options{
		InKey: os.Getenv("X509_KEY_FILE"),
		In:    os.Getenv("X509_CERT_FILE"),
		Out:   os.Getenv("X509_P12_FILE"),
	}
	fmt.Println(opts.CreateCommand())

	opts2 := p12.Options{In: os.Getenv("X509_P12_FILE")}
	fmt.Println(opts2.InfoCommand())

	fmt.Println("DONE")
}
