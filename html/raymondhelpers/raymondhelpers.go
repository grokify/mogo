package raymondhelpers

import (
	"regexp"
	"strings"
	"time"

	"github.com/aymerick/raymond"
	"github.com/grokify/gotilla/time/timeutil"
)

// RegisterAll registers helpers for the Raymond Handlebars template
// engine.
func RegisterAll() {
	raymond.RegisterHelper("timeRfc3339", func(t time.Time) string {
		return t.Format(time.RFC3339)
	})
	raymond.RegisterHelper("timeRfc3339ymd", func(t time.Time) string {
		return t.Format(timeutil.RFC3339YMD)
	})
	raymond.RegisterHelper("spaceToHyphen", func(s string) string {
		return regexp.MustCompile(`[\s-]+`).ReplaceAllString(s, "-")
	})
	raymond.RegisterHelper("spaceToUnderscore", func(s string) string {
		return regexp.MustCompile(`[\s_]+`).ReplaceAllString(s, "_")
	})
	raymond.RegisterHelper("toLower", func(s string) string {
		return strings.ToLower(s)
	})
}
