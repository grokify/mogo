package main

import (
	"fmt"
	"html"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/grokify/mogo/net/http/httputilmore"
)

func LogEveryLevels(lgr *slog.Logger, inputLevel, inputURL string) {
	lgr.Debug("Example Dynamic Log Level Entry", "contextLevel", inputLevel, "inputURL", inputURL)
	lgr.Info("Example Dynamic Log Level Entry", "contextLevel", inputLevel, "inputURL", inputURL)
	lgr.Warn("Example Dynamic Log Level Entry", "contextLevel", inputLevel, "inputURL", inputURL)
	lgr.Error("Example Dynamic Log Level Entry", "contextLevel", inputLevel, "inputURL", inputURL)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var levelVar slog.LevelVar
		levelVar.Set(slog.LevelInfo) // Default levels
		lgr := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: &levelVar}))

		inputLevel := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("level")))
		switch inputLevel {
		case "DEBUG":
			levelVar.Set(slog.LevelDebug)
		case "INFO":
			levelVar.Set(slog.LevelInfo)
		case "WARN":
			levelVar.Set(slog.LevelWarn)
		case "ERROR":
			levelVar.Set(slog.LevelError)
		}

		LogEveryLevels(lgr, inputLevel, r.URL.RequestURI())

		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path)) //nolint:gosec // G705: already sanitized with html.EscapeString
	})

	svr := httputilmore.NewServerTimeouts(":8080", mux, 1000*time.Second)

	log.Fatal(svr.ListenAndServe())
}
