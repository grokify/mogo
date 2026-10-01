package main

import (
	"fmt"
	"log"
	"net/mail"

	"github.com/grokify/mogo/mime/multipartutil"
	"github.com/grokify/mogo/net/http/httputilmore"
	"github.com/grokify/mogo/net/mailutil"
)

func main() {
	msg := mailutil.MessageWriter{
		From: &mail.Address{
			Name:    "Alice",
			Address: "alice@example.com",
		},
		To: []mail.Address{
			{Address: "bob@example.com"},
		},
		Subject: "Hello World!",
		BodyPartsSet: multipartutil.PartsSet{
			ContentType: httputilmore.ContentTypeMultipartMixed,
			Parts: multipartutil.Parts{
				{
					Type:        multipartutil.PartTypeRaw,
					ContentType: httputilmore.ContentTypeTextPlainUtf8,
					BodyDataRaw: []byte("Hola!"),
				},
			},
		},
	}

	str, err := msg.String()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(str)

	fmt.Println("DONE")
}
