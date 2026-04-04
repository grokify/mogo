// Copyright 2026 John Wang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gosec

import (
	"fmt"
	"net/http"
)

// Nolint formats a nolint:gosec comment with the given rule and reason.
//
// Example:
//
//	comment := gosec.Nolint("G117", "OAuth token response per RFC 6749")
//	// Returns: "//nolint:gosec // G117: OAuth token response per RFC 6749"
func Nolint(rule, reason string) string {
	return fmt.Sprintf("//nolint:gosec // %s: %s", rule, reason)
}

// NolintG101 returns a nolint comment for G101 (hardcoded credentials).
//
// Use this when a string matches credential patterns but is not actually
// a credential (e.g., URL paths, test fixtures, documentation).
//
// Example reasons:
//   - "URL path, not a credential"
//   - "Test fixture with fake credentials"
//   - "Documentation example"
func NolintG101(reason string) string {
	return Nolint("G101", reason)
}

// NolintG115 returns a nolint comment for G115 (integer overflow conversion).
//
// Use this when an integer conversion is known to be safe due to
// validated input ranges or domain constraints.
//
// Example reasons:
//   - "Value bounded by prior validation"
//   - "Domain constraint ensures safe range"
//   - "Year value always fits in int32"
func NolintG115(reason string) string {
	return Nolint("G115", reason)
}

// NolintG117 returns a nolint comment for G117 (secret in JSON response).
//
// Use this when marshaling structs with intentional secret fields like
// OAuth access_token, client_secret, etc.
//
// Example reasons:
//   - "OAuth token response per RFC 6749"
//   - "OAuth registration response per RFC 7591"
//   - "API key response for authenticated user"
func NolintG117(reason string) string {
	return Nolint("G117", reason)
}

// NolintG118 returns a nolint comment for G118 (context.Background in goroutine).
//
// Use this when a goroutine intentionally uses context.Background because
// the request context is not appropriate (e.g., shutdown handlers).
//
// Example reasons:
//   - "Shutdown handler runs after request context is cancelled"
//   - "Background job outlives request lifecycle"
//   - "Cleanup routine needs independent timeout"
func NolintG118(reason string) string {
	return Nolint("G118", reason)
}

// NolintG703 returns a nolint comment for G703 (path traversal via taint analysis).
//
// IMPORTANT: Only use this in cmd/ directories (CLI entry points) where users
// explicitly provide paths. For library code in pkg/, use osutil.WriteFileSecure,
// osutil.CopyFileSecure, or osutil.CleanPathSecure instead.
//
// CLI entry points should:
//  1. Use filepath.Clean() on user-provided paths before calling libraries
//  2. Use this nolint comment to suppress the warning
//
// Example reasons for cmd/:
//   - "Path from CLI flag"
//   - "Path from user config"
//   - "Output path specified by user"
func NolintG703(reason string) string {
	return Nolint("G703", reason)
}

// NolintG704 returns a nolint comment for G704 (SSRF via taint analysis).
//
// Use this when making HTTP requests to URLs from trusted sources.
//
// Example reasons:
//   - "Test uses httptest server URL"
//   - "URL from validated allowlist"
//   - "Internal service URL from config"
//   - "URL constructed from trusted constants"
func NolintG704(reason string) string {
	return Nolint("G704", reason)
}

// NolintG122 returns a nolint comment for G122 (filepath.Walk TOCTOU race).
//
// IMPORTANT: Only use this in cmd/ directories (CLI entry points) where users
// explicitly provide directories. For library code in pkg/, use os.Root (Go 1.24+)
// or osutil.ReadDirFilesSecure to perform symlink-safe filesystem operations.
//
// G122 detects TOCTOU (time-of-check-time-of-use) race conditions in
// filepath.Walk/WalkDir callbacks where filesystem operations use the
// potentially race-prone path provided by Walk.
//
// For pkg/ code, use osutil.ReadDirFilesSecure or os.Root directly:
//
//	// Option 1: Use osutil helper
//	files, err := osutil.ReadDirFilesSecure(dir)
//
//	// Option 2: Use os.Root directly
//	root, err := os.OpenRoot(dir)
//	if err != nil { return err }
//	defer root.Close()
//	content, err := root.ReadFile(relativePath) // Symlink-safe
//
// Example reasons for cmd/:
//   - "Directory from CLI flag"
//   - "Walking config directory from trusted source"
func NolintG122(reason string) string {
	return Nolint("G122", reason)
}

// CommonReasons provides pre-written reason strings for common scenarios.
//
//nolint:gosec // G101: These are reason strings, not credentials
var CommonReasons = struct {
	// G101 reasons
	URLPathNotCredential string
	TestFixture          string
	DocumentationExample string

	// G115 reasons
	BoundedByValidation string
	DomainConstraint    string
	YearValueSafeRange  string
	SmallEnumValue      string

	// G117 reasons
	OAuthTokenResponse        string
	OAuthRegistrationResponse string

	// G118 reasons
	ShutdownHandler   string
	BackgroundJob     string
	CleanupRoutine    string
	IndependentCancel string

	// G703 reasons (use only in cmd/, not library code)
	PathFromCLIFlag  string
	PathFromConfig   string
	OutputPathByUser string

	// G704 reasons
	HttptestServer      string
	ValidatedAllowlist  string
	InternalServiceURL  string
	TrustedConstantsURL string

	// G122 reasons (cmd/ only - use os.Root in pkg/)
	DirectoryFromCLIFlag string
	TrustedConfigDir     string
}{
	// G101
	URLPathNotCredential: "URL path, not a credential",
	TestFixture:          "Test fixture with fake credentials",
	DocumentationExample: "Documentation example",

	// G115
	BoundedByValidation: "Value bounded by prior validation",
	DomainConstraint:    "Domain constraint ensures safe range",
	YearValueSafeRange:  "Year value always fits in target type",
	SmallEnumValue:      "Small enum value, no overflow possible",

	// G117
	OAuthTokenResponse:        "OAuth token response per RFC 6749",
	OAuthRegistrationResponse: "OAuth registration response per RFC 7591",

	// G118
	ShutdownHandler:   "Shutdown handler runs after request context is cancelled",
	BackgroundJob:     "Background job outlives request lifecycle",
	CleanupRoutine:    "Cleanup routine needs independent timeout",
	IndependentCancel: "Requires independent cancellation from request",

	// G703 (use only in cmd/, not library code)
	PathFromCLIFlag:  "Path from CLI flag",
	PathFromConfig:   "Path from user config",
	OutputPathByUser: "Output path specified by user",

	// G704
	HttptestServer:      "Test uses httptest server URL",
	ValidatedAllowlist:  "URL from validated allowlist",
	InternalServiceURL:  "Internal service URL from config",
	TrustedConstantsURL: "URL constructed from trusted constants",

	// G122 (cmd/ only - use os.Root in pkg/)
	DirectoryFromCLIFlag: "Directory from CLI flag",
	TrustedConfigDir:     "Walking config directory from trusted source",
}

// G120 Fix Helpers
//
// G120 warns about parsing form data without limiting request body size.
// Unlike other gosec rules, G120 requires code changes rather than nolint.
//
// The fix requires:
//  1. Call http.MaxBytesReader to limit body size (MUST be inline, not a helper)
//  2. Call r.ParseForm() or r.ParseMultipartForm()
//  3. Use r.Form.Get() instead of r.FormValue()
//
// Caveats (gosec 2.11+):
//   - Only inline http.MaxBytesReader is recognized; helper functions are not
//   - r.FormValue() is still flagged even after ParseForm; use r.Form.Get()

// G120MaxBytes provides common max body size limits for G120 fixes.
var G120MaxBytes = struct {
	// Form is the default limit for simple form submissions (1MB).
	Form int64

	// Multipart is the limit for file uploads (32MB).
	Multipart int64

	// Webhook is the limit for webhook payloads (64KB).
	Webhook int64

	// Twilio is the limit for Twilio webhook callbacks (64KB).
	// Twilio webhook bodies are typically small (under 10KB).
	Twilio int64
}{
	Form:      1 << 20,  // 1MB
	Multipart: 32 << 20, // 32MB
	Webhook:   64 << 10, // 64KB
	Twilio:    64 << 10, // 64KB
}

// LimitAndParseForm limits the request body and parses form data.
// This is the recommended pattern to fix G120, but note that gosec 2.11+
// may not recognize helper functions - copy the inline pattern if needed.
//
// After calling this, use r.Form.Get() instead of r.FormValue().
//
// Example:
//
//	if err := gosec.LimitAndParseForm(w, r, gosec.G120MaxBytes.Webhook); err != nil {
//	    http.Error(w, "Bad Request", http.StatusBadRequest)
//	    return
//	}
//	value := r.Form.Get("key")
func LimitAndParseForm(w http.ResponseWriter, r *http.Request, maxBytes int64) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
	return r.ParseForm()
}

// G120InlinePattern returns a code comment with the inline fix pattern.
// Use this as documentation when gosec doesn't recognize helper functions.
func G120InlinePattern() string {
	return `// G120 fix pattern (inline for gosec 2.11+ compatibility):
//
//   r.Body = http.MaxBytesReader(w, r.Body, 64<<10) // 64KB
//   if err := r.ParseForm(); err != nil {
//       http.Error(w, "Bad Request", http.StatusBadRequest)
//       return
//   }
//   value := r.Form.Get("key") // NOT r.FormValue("key")`
}
