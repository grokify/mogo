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
	RegisterStringSafe()
	RegisterTimeSafe()
}

func RegisterTimeSafe() {
	raymond.RegisterHelper("timeRfc3339", func(t time.Time) raymond.SafeString {
		return raymond.SafeString(t.Format(time.RFC3339))
	})
	raymond.RegisterHelper("timeRfc3339ymd", func(t time.Time) raymond.SafeString {
		return raymond.SafeString(t.Format(timeutil.RFC3339YMD))
	})
}

func RegisterStringSafe() {
	raymond.RegisterHelper("spaceToHyphen", func(s string) raymond.SafeString {
		return raymond.SafeString(regexp.MustCompile(`[\s-]+`).ReplaceAllString(s, "-"))
	})
	raymond.RegisterHelper("spaceToUnderscore", func(s string) raymond.SafeString {
		return raymond.SafeString(regexp.MustCompile(`[\s_]+`).ReplaceAllString(s, "_"))
	})
	raymond.RegisterHelper("toLower", func(s string) raymond.SafeString {
		return raymond.SafeString(strings.ToLower(s))
	})
}
