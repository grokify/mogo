# lintfix

Structured lint remediation database for Go projects using golangci-lint.

## Overview

The `lintfix` package provides:

- 📋 **Remediation database** - Embedded JSON database mapping lint rules to fixes
- 🔧 **Helper references** - Links to mogo helper functions for code fixes
- 📝 **Nolint generators** - Properly formatted nolint comments with documented reasons
- 📚 **Documentation** - Version-specific caveats and best practices

## Quick Start

```go
import (
    "github.com/grokify/mogo/lintfix"
    "github.com/grokify/mogo/lintfix/gosec"
)

// Query the remediation database
db := lintfix.MustLoadRemediations()
fix := db.GetGosec("G120")
fmt.Println(fix.Remediation.Summary)
// "Use http.MaxBytesReader inline before parsing form data"

// Generate nolint comments
comment := gosec.NolintG117(gosec.CommonReasons.OAuthTokenResponse)
// "//nolint:gosec // G117: OAuth token response per RFC 6749"
```

## Remediation Types

| Type | Description | Example |
|------|-------------|---------|
| `code` | Add/modify code with helper functions | G120: Use `http.MaxBytesReader` |
| `nolint` | Add nolint annotation with reason | G117: OAuth token response |
| `refactor` | Broader code changes needed | G101: Move secrets to env vars |

## Supported Linters

- **gosec** - Security-focused rules (G101, G112, G115, G117, G118, G120, G122, G124, G401, G501, G601, G703, G704, G705, G706, G710)
- **staticcheck** - Static analysis (SA1019, SA4006, QF1012)
- **errcheck** - Error handling
- **govet** - Inline remediation notes
- **dupl** - Duplicate code detection

## G703: Path Traversal

G703 warns about file paths constructed from user input. The fix depends on where your code lives:

**In `cmd/` (CLI entry points)** - User explicitly provides the path, use nolint:

```go
// User provides path via CLI flag - they own the risk
cleanPath := filepath.Clean(userPath)
if err := os.WriteFile(cleanPath, data, 0600); err != nil { //nolint:gosec // G703: Path from CLI flag
    return err
}
```

**In library code** - Use secure functions that reject `..` sequences:

```go
import "github.com/grokify/mogo/os/osutil"

// Library code - reject paths with traversal sequences
data, err := osutil.ReadFileSecure(path)
if err != nil {
    // Returns: "path contains '..' traversal sequence: ../etc/passwd"
    return err
}

if err := osutil.WriteFileSecure(path, data, 0600); err != nil {
    return err
}
```

**Error returned:** `osutil.ErrPathTraversal` is returned when a path contains `..`:

```go
// errors.Is check
if errors.Is(err, osutil.ErrPathTraversal) {
    log.Println("Invalid path:", err)
}
```

## G101: Config Struct Fields Set From Parameters

G101 also fires on struct literals with credential-named fields (`ClientSecret`, `APIKey`,
`Password`, `Token`, ...) even when the values come from caller-supplied parameters, not
literals - a common shape for any OAuth/API-client config constructor:

```go
func (s *OAuthService) ConfigureGoogle(clientID, clientSecret, redirectURL string) {
	s.RegisterProvider(&OAuthProvider{ //nolint:gosec // G101: ClientID/ClientSecret are set from caller-supplied parameters, not hardcoded literals
		Name:         "google",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
	})
}
```

There is no code fix here - the struct shape is the point, and gosec cannot see that the
values are parameters rather than literals. `nolint` is the correct remediation.

## G706: Log Injection

G706 warns when a value derived from client input (request Host, headers, path, etc.) is
written directly to a log call, since an unescaped newline or control character lets an
attacker forge fake log lines (CWE-117).

**Verified fix - wrap with `strconv.Quote`, not just the `%q` verb:**

```go
import "strconv"

// Correct: strconv.Quote is a recognized sanitizer, clears the finding
log.Printf("Proxy error for %s: %v", strconv.Quote(r.Host), err)
```

```go
// Does NOT clear the finding: gosec inspects the argument expression, not the
// format verb, so the raw tainted value is still flagged even with %q
log.Printf("Proxy error for %q: %v", r.Host, err) // still G706
```

Prefer this code fix over `nolint` in library code - it's a real fix (escapes injected
control characters), not just linter appeasement, and it's what `gosec.NolintG706` is
documented to defer to.

## G710: Open Redirect

G710 warns when an `http.Redirect` target is built by concatenating request-derived data
(e.g. `"https://" + r.Host + r.RequestURI`), since an attacker who controls the Host header
could make the server redirect anywhere (CWE-601).

**Verified fix - build the target with `net/url.URL`, not string concatenation:**

```go
import "net/url"

// Correct: url.URL{}.String() is the recognized safe code shape, clears the finding
target := url.URL{Scheme: "https", Host: r.Host, Path: r.URL.Path, RawQuery: r.URL.RawQuery}
http.Redirect(w, r, target.String(), http.StatusMovedPermanently)
```

**Important - this clears the linter, not the actual vulnerability.** Verified empirically:
gosec accepts the `url.URL{}` construction on its own, with no host validation at all. The
real security fix is a separate step - validate the host against a known allowlist (e.g.
the backends your proxy actually serves) before redirecting:

```go
if !isKnownHost(r.Host) { // e.g. rp.findProxy(r.Host) != nil in a reverse proxy
    http.NotFound(w, r)
    return
}
target := url.URL{Scheme: "https", Host: r.Host, Path: r.URL.Path, RawQuery: r.URL.RawQuery}
http.Redirect(w, r, target.String(), http.StatusMovedPermanently)
```

Do both. Do not treat "gosec is clean" as evidence that a request-derived redirect target
is actually safe.

## Nolint Generators

The `gosec` subpackage provides type-safe nolint comment generators:

```go
gosec.NolintG101(reason)  // Hardcoded credentials (false positive)
gosec.NolintG115(reason)  // Integer overflow (bounded value)
gosec.NolintG117(reason)  // Secret in JSON response
gosec.NolintG118(reason)  // context.Background in goroutine
gosec.NolintG122(reason)  // filepath.Walk TOCTOU race (cmd/ entry point only)
gosec.NolintG124(reason)  // Insecure cookie attributes (set dynamically/from config)
gosec.NolintG703(reason)  // Path traversal (CLI entry point only)
gosec.NolintG704(reason)  // SSRF (trusted URL)
gosec.NolintG705(reason)  // XSS (trusted content)
gosec.NolintG706(reason)  // Log injection (prefer the strconv.Quote code fix instead)
gosec.NolintG710(reason)  // Open redirect (prefer the url.URL{} code fix instead)
```

### Common Reasons

Pre-written reason strings for common scenarios:

```go
gosec.CommonReasons.OAuthTokenResponse        // G117
gosec.CommonReasons.ShutdownHandler           // G118
gosec.CommonReasons.PathFromCLIFlag           // G703
gosec.CommonReasons.HttptestServer            // G704
gosec.CommonReasons.BoundedByValidation       // G115
gosec.CommonReasons.ParameterNotLiteral       // G101 - config struct field set from a parameter
gosec.CommonReasons.TestControlledInputNoUntrustedSource // G706 - nolint fallback only; prefer strconv.Quote
```

## Documentation

- [Gosec Version Caveats](../docs/lintfix/gosec-caveats.md) - Version-specific behaviors
- [GoDoc](https://pkg.go.dev/github.com/grokify/mogo/lintfix) - API reference

## Adding New Rules

Edit `remediations.json` to add new rules:

```json
{
  "linters": {
    "gosec": {
      "G999": {
        "name": "Rule name",
        "description": "What the rule detects",
        "severity": "high|medium|low",
        "category": "security|correctness|maintenance",
        "remediation": {
          "type": "code|nolint|refactor",
          "summary": "Brief fix description",
          "example": "Code example"
        }
      }
    }
  }
}
```
