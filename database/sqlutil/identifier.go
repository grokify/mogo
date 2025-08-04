package sqlutil

import "regexp"

var rxUnquotedIdentifier = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

func IsUnquotedIdentifier(s string) bool {
	return rxUnquotedIdentifier.MatchString(s)
}
