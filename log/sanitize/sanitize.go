// Package sanitize provides functions for sanitizing user input before logging.
//
// Log injection (CWE-117) occurs when untrusted data is written to logs without
// proper neutralization, allowing attackers to forge log entries, corrupt log
// integrity, or inject malicious content.
//
// Reference: https://cwe.mitre.org/data/definitions/117.html
//
// Example attack vectors:
//   - Newline injection: "legitimate\nERROR: fake entry" creates fake log lines
//   - Carriage return injection: overwrites log lines on some terminals
//   - Control character injection: can corrupt log files or exploit log viewers
//
// Usage:
//
//	log.Printf("user=%s action=%s", sanitize.String(username), sanitize.String(action))
//	slog.Info("request", "session_id", sanitize.String(sessionID))
package sanitize

import (
	"strings"
	"unicode"
)

// String removes or replaces characters that could be used for log injection.
// This includes newlines, carriage returns, and other ASCII control characters.
//
// Use this function when logging any user-controlled input to prevent CWE-117.
//
// Example:
//
//	log.Printf("session=%s", sanitize.String(req.SessionID))
func String(s string) string {
	return strings.Map(func(r rune) rune {
		// Remove ASCII control characters (0x00-0x1F and 0x7F)
		if r < 0x20 || r == 0x7F {
			return -1 // -1 means delete the rune
		}
		return r
	}, s)
}

// StringReplace is like String but replaces control characters with a
// replacement string instead of removing them. This preserves the visual
// indication that something was sanitized.
//
// Example:
//
//	log.Printf("input=%s", sanitize.StringReplace(userInput, "?"))
func StringReplace(s, replacement string) string {
	var b strings.Builder
	b.Grow(len(s))

	for _, r := range s {
		if r < 0x20 || r == 0x7F {
			b.WriteString(replacement)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// Strings sanitizes multiple strings, returning a new slice.
//
// Example:
//
//	safe := sanitize.Strings(userID, sessionID, action)
func Strings(values ...string) []string {
	result := make([]string, len(values))
	for i, v := range values {
		result[i] = String(v)
	}
	return result
}

// StringOrTruncate sanitizes a string and truncates it to maxLen if longer.
// This is useful for logging potentially large user inputs.
//
// Example:
//
//	log.Printf("body=%s", sanitize.StringOrTruncate(requestBody, 1000))
func StringOrTruncate(s string, maxLen int) string {
	s = String(s)
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// IsClean returns true if the string contains no control characters
// that would require sanitization.
//
// Example:
//
//	if !sanitize.IsClean(sessionID) {
//	    log.Warn("suspicious session ID detected")
//	}
func IsClean(s string) bool {
	for _, r := range s {
		if r < 0x20 || r == 0x7F {
			return false
		}
	}
	return true
}

// HasControlChars returns true if the string contains any Unicode control
// characters (category Cc), which is broader than just ASCII control chars.
func HasControlChars(s string) bool {
	for _, r := range s {
		if unicode.IsControl(r) {
			return true
		}
	}
	return false
}

// StripAllControl removes all Unicode control characters (category Cc),
// not just ASCII control characters. Use this for stricter sanitization.
func StripAllControl(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, s)
}
