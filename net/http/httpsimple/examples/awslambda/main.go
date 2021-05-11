package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/apex/gateway"
)

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/hello", hello)
	log.Fatal(gateway.ListenAndServe(":3000", nil))
}

func hello(w http.ResponseWriter, r *http.Request) {
	// example retrieving values from the api gateway proxy request context.
	requestContext, ok := gateway.RequestContext(r.Context())
	if !ok || requestContext.Authorizer["sub"] == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Hello World from Go")
		return
	}

	userID := requestContext.Authorizer["sub"].(string)
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "Hello %s from Go", userID)
}
