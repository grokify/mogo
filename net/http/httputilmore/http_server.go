package httputilmore

import (
	"log/slog"
	"net/http"
	"time"
)

// Log is a custom Http handler that will log all requests.
// It can be called using
// http.ListenAndServe(":8080", Log(http.DefaultServeMux))
// From: https://groups.google.com/forum/#!topic/golang-nuts/s7Xk1q0LSU0
func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("log handler func", "remote-addr", r.RemoteAddr, "req-method", r.Method, "req-url", r.URL)
		handler.ServeHTTP(w, r)
	})
}

// NewServerTimeouts returns a `*http.Server` with all timeouts set to a single value provided.
func NewServerTimeouts(addr string, handler http.Handler, timeout time.Duration) *http.Server {
	if timeout < 0 {
		timeout = 0
	}
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		IdleTimeout:       timeout,
		ReadHeaderTimeout: timeout,
		ReadTimeout:       timeout,
		WriteTimeout:      timeout,
		MaxHeaderBytes:    1 << 20,
	}
}

func ListenAndServeTimeouts(addr string, handler http.Handler, timeout time.Duration) error {
	return NewServerTimeouts(addr, handler, timeout).ListenAndServe()
}
