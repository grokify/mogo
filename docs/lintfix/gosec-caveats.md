# Gosec Version-Specific Caveats

This document describes version-specific behaviors in gosec that affect remediation strategies.

## G120: Unbounded Request Body

### gosec 2.11+ Behavior Changes

Starting with gosec 2.11 (included in golangci-lint 2.11+), the G120 rule has stricter detection:

#### 1. Helper Functions Not Recognized

gosec only recognizes **inline** `http.MaxBytesReader` calls. Helper functions that wrap this call are not detected.

**Does NOT work with gosec 2.11+:**

```go
// Helper function - gosec doesn't trace the call
httputilmore.LimitRequestBody(w, r, httputilmore.DefaultMaxBodySize)
if err := r.ParseForm(); err != nil { ... }  // G120 flagged
```

**Works with gosec 2.11+:**

```go
// Inline call - gosec recognizes this pattern
r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
if err := r.ParseForm(); err != nil { ... }  // OK
```

#### 2. FormValue() Flagged After ParseForm()

gosec 2.11+ flags `r.FormValue()` calls even when `ParseForm()` was already called with a limited body. This is because `FormValue()` internally calls `ParseForm()` if not already parsed, and gosec doesn't track that the form is already parsed.

**Does NOT work with gosec 2.11+:**

```go
r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
if err := r.ParseForm(); err != nil { ... }
value := r.FormValue("key")  // G120 flagged!
```

**Works with gosec 2.11+:**

```go
r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
if err := r.ParseForm(); err != nil { ... }
value := r.Form.Get("key")  // OK - directly accesses parsed form
```

### Complete Example

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // 1. Limit body size INLINE (gosec recognizes this)
    r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB

    // 2. Parse the form
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    // 3. Use r.Form.Get() instead of r.FormValue()
    username := r.Form.Get("username")
    password := r.Form.Get("password")

    // ... handle request
}
```

### Why This Matters

The stricter detection in gosec 2.11+ means:

1. **Helper functions** like `httputilmore.LimitRequestBody()` provide good abstractions but won't satisfy the linter
2. **FormValue()** is a common pattern but triggers false positives after proper limiting
3. **Local/CI version mismatch** can cause CI failures that don't reproduce locally

### Recommendations

1. **Keep golangci-lint versions in sync** between local development and CI
2. **Use inline `http.MaxBytesReader`** rather than helper functions for G120
3. **Use `r.Form.Get()`** instead of `r.FormValue()` after calling `ParseForm()`
4. **Document the pattern** in code comments referencing G120

### Version Reference

| golangci-lint | gosec | G120 Behavior |
|---------------|-------|---------------|
| 2.10.x | 2.21.x | Recognizes inline MaxBytesReader only |
| 2.11.x | 2.24.x | Same + flags FormValue() after ParseForm |

## G703: Path Traversal

G703 warns about file paths constructed from user input that may allow directory traversal attacks.

### cmd/ vs Library Code

The fix depends on where your code lives:

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

**Error handling:**

```go
// Check for path traversal error specifically
if errors.Is(err, osutil.ErrPathTraversal) {
    log.Println("Invalid path:", err)
}
```

### Additional Validation Patterns

**Validate input to reject path separators:**

```go
func findExecutable(name string) (string, error) {
    // Validate name doesn't contain path separators
    if strings.ContainsAny(name, `/\`) {
        return "", fmt.Errorf("invalid name: %s", name)
    }

    fullPath := filepath.Join(dir, name)
    if info, err := os.Stat(fullPath); err == nil { // #nosec G703 - name validated above
        return fullPath, nil
    }
    // ...
}
```

**Use filepath.Base for filenames:**

```go
// Sanitize to just the filename
safeName := filepath.Base(userInput)
fullPath := filepath.Join(baseDir, safeName)
```

## G706: Log Injection via Taint Analysis

G706 flags request-derived values (Host, headers, path, etc.) passed to a log call, since
an unescaped newline or control character lets an attacker forge fake log lines (CWE-117).

### Format Verb vs. Sanitizing Call

Changing only the `Printf` verb does **not** clear the finding - gosec's taint analysis
inspects the argument *expression*, not the verb string, so a raw tainted value is still
flagged even when quoted for display purposes only:

**Does NOT clear G706:**

```go
log.Printf("Proxy error for %q: %v", r.Host, err) // still flagged - r.Host is unchanged
```

**Clears G706** - wrap the value in an actual call to `strconv.Quote`, which gosec
recognizes as a sanitizer boundary:

```go
import "strconv"

log.Printf("Proxy error for %s: %v", strconv.Quote(r.Host), err)
```

This is verified empirically (golangci-lint 2.12.2 / gosec bundled therein), not merely
inferred from the rule name. It is also a genuine fix, not just linter appeasement:
`strconv.Quote` escapes newlines and control characters, so injected log lines are
rendered inert (visible as `\n` inside the quoted string) rather than executed as literal
line breaks in the log stream.

### When Nolint Is Appropriate Instead

Use `gosec.NolintG706` only when the value cannot reasonably be quoted (e.g., written to a
structured/non-text log sink) or in test code with fully controlled input - not as a
substitute for the `strconv.Quote` fix in production code that logs request data.

## G710: Open Redirect via Taint Analysis

G710 flags an `http.Redirect` target built by string concatenation from request-derived
data (e.g. `"https://" + r.Host + r.RequestURI`), since an attacker controlling the Host
header could make the server redirect anywhere (CWE-601).

### Code Shape vs. Actual Validation

**Does NOT clear G710** - string concatenation, even downstream of a validation check:

```go
if rp.findProxy(r.Host) == nil { // real validation - but gosec doesn't see it
    http.NotFound(w, r)
    return
}
target := "https://" + r.Host + r.RequestURI // still flagged
http.Redirect(w, r, target, http.StatusMovedPermanently)
```

**Clears G710** - construct the target with `net/url.URL` and call `.String()`:

```go
import "net/url"

target := url.URL{Scheme: "https", Host: r.Host, Path: r.URL.Path, RawQuery: r.URL.RawQuery}
http.Redirect(w, r, target.String(), http.StatusMovedPermanently)
```

**Verified empirically that this is a pure code-shape check**: the `url.URL{}` +
`.String()` construction clears G710 *on its own*, with the host-validation check removed
entirely. gosec has no way to detect or require an allowlist check - it is matching the
construction pattern only.

### This Means the Linter and the Vulnerability Are Two Different Things

Do not treat a clean gosec run as evidence that a request-derived redirect is safe. The
`url.URL{}` shape satisfies the linter; an allowlist/known-host check is what actually
prevents the open redirect. Ship both:

```go
if !isKnownHost(r.Host) {
    http.NotFound(w, r)
    return
}
target := url.URL{Scheme: "https", Host: r.Host, Path: r.URL.Path, RawQuery: r.URL.RawQuery}
http.Redirect(w, r, target.String(), http.StatusMovedPermanently)
```

A `nolint` is appropriate only when the redirect target is fully server-controlled (e.g. a
constant string) - never when any part of it derives from the request.

## G101: Credential-Named Fields Set From Parameters

G101 pattern-matches on struct fields/variables whose *names* look like credentials
(`ClientSecret`, `APIKey`, `Password`, `Token`, ...), independent of whether the assigned
*value* is a literal. This makes it a reliable false positive on any config/options struct
built from caller-supplied parameters:

```go
func (s *OAuthService) ConfigureGoogle(clientID, clientSecret, redirectURL string) {
	s.RegisterProvider(&OAuthProvider{ //nolint:gosec // G101: ClientID/ClientSecret are set from caller-supplied parameters, not hardcoded literals
		ClientID:     clientID,     // parameter, not a literal
		ClientSecret: clientSecret, // parameter, not a literal
		RedirectURL:  redirectURL,
	})
}
```

There is no code-shape fix for this one (unlike G706/G710 above) - the struct literal
*is* the correct code, and gosec cannot distinguish "field set from a parameter" from
"field set from a literal" by field name alone. `gosec.NolintG101` with
`gosec.CommonReasons.ParameterNotLiteral` is the right remediation.

## G115: Integer Overflow Conversion

### When to Use Nolint

G115 flags integer conversions that could overflow. Use nolint when:

1. **Domain constraints guarantee safety**: Year values (2020-2100) fit in any integer type
2. **Prior validation bounds the value**: Input already checked to be in safe range
3. **Small enum values**: Converting small constants that obviously fit

**Example - Year conversion (safe):**

```go
year := time.Now().Year()
// Year is always ~2000-2100, fits in any integer type
prefix := fmt.Sprintf("PRD-%d", year)  // Use fmt instead of manual rune conversion
```

**Example - Validated input:**

```go
if value < 0 || value > 255 {
    return errors.New("value out of range")
}
b := byte(value) // #nosec G115 - value bounded by validation above
```

## Keeping Versions in Sync

To avoid local/CI lint mismatches:

```bash
# Check local version
golangci-lint --version

# Install specific version to match CI
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.1
```

Or pin the version in CI to match local:

```yaml
# .github/workflows/lint.yaml
- uses: golangci/golangci-lint-action@v8
  with:
    version: v2.11.1
```
