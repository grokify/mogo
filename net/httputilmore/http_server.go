package httputilmore

import (
	"log"
	"net/http"
)

// Log is a custom Http handler that will log all requests.
// It can be called using
// http.ListenAndServe(":8080", Log(http.DefaultServeMux))
// From: https://groups.google.com/forum/#!topic/golang-nuts/s7Xk1q0LSU0
func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
