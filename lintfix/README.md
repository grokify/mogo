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

```go
import "github.com/grokify/mogo/lintfix/dupl"

comment := dupl.Nolint(dupl.CommonReasons.ParallelResourceWrapper)
// "//nolint:dupl // Structurally parallel to sibling wrapper methods over
// distinct generated types; not meaningfully extractable without reflection
// or per-type adapters"
```

```go
import "github.com/grokify/mogo/lintfix/unparam"

comment := unparam.Nolint(unparam.CommonReasons.InterfaceSignature)
// "//nolint:unparam // Signature fixed by an interface method set this type implements"
```

## Remediation Types

| Type | Description | Example |
|------|-------------|---------|
| `code` | Add/modify code with helper functions | G120: Use `http.MaxBytesReader` |
| `nolint` | Add nolint annotation with reason | G117: OAuth token response |
| `refactor` | Broader code changes needed | G101: Move secrets to env vars |

## Supported Linters

- **gosec** - Security-focused rules (G101, G112, G115, G117, G118, G120, G122, G124, G401, G404, G501, G601, G703, G704, G705, G706, G710)
- **staticcheck** - Static analysis (SA1019, SA4006, QF1003, QF1012)
- **errcheck** - Error handling
- **govet** - Inline remediation notes
- **dupl** - Duplicate code detection; see the `dupl` subpackage for nolint generators covering the generated-client-wrapper case
- **unparam** - Unused function parameters/results; see the `unparam` subpackage for nolint generators covering interface/callback-constrained signatures
- **unused** - Dead code (unused functions, vars, consts, types); the fix is deletion, including any import the removed code solely required

## G404: Weak Random Number Generator

G404 flags any use of `math/rand` or `math/rand/v2` — it has no way to tell whether
the value is used for something security-sensitive (tokens, keys, nonces, passwords)
or not (shuffling display data, jitter, sampling, non-cryptographic test fixtures).
Golangci-lint version skew commonly surfaces this: an older locally-installed gosec
may not flag a call that a newer one (e.g. CI's `version: latest`) does, since gosec's
G404 detection coverage (which stdlib functions it recognizes, e.g. `rand.Shuffle`)
has expanded across releases.

**If the value IS security-sensitive** - switch to `crypto/rand`, don't nolint:

```go
import "crypto/rand"

n, err := rand.Int(rand.Reader, max)
if err != nil {
    return err
}
```

**If the value is NOT security-sensitive** - nolint with a reason:

```go
import "math/rand"

//nolint:gosec // G404: Shuffling display data, not security-sensitive
rand.Shuffle(len(items), func(i, j int) { items[i], items[j] = items[j], items[i] })
```

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

## G115: Integer Overflow Conversion — Length-Prefixed External Data

G115's default remediation (a `nolint` backed by "domain knowledge guarantees
the value fits") is right for small constants/enums, but wrong for a
length/size prefix read off a file format or wire protocol — there the value
comes from outside the program and nothing guarantees it fits until you check.

**A common bad shape**: converting through a narrower signed type first, then
validating the *signed* result:

```go
n := int(int32(binary.LittleEndian.Uint32(lenBuf[:])))
if n < 5 {
    return fmt.Errorf("invalid length %d", n)
}
```

This only catches lengths that wrapped negative (raw values ≥ 2^31). Any raw
value below that — up to `2^31-1`, over two billion — sails through as a
large *positive* `n` and drives `make([]byte, n)` with an
attacker-or-corruption-controlled size, unless something downstream happens
to bound it separately.

**Verified fix — validate the raw unsigned value's range before converting:**

```go
const maxDocSize = 16 * 1024 * 1024 // real domain ceiling, not an arbitrary guess

raw := binary.LittleEndian.Uint32(lenBuf[:])
if raw < 5 || raw > maxDocSize {
    return fmt.Errorf("invalid length %d", raw)
}
n := int(raw) // safe: raw is now proven in [5, maxDocSize]
```

This clears the G115 finding (the conversion now only ever sees a
pre-validated range) and is a strictly stronger real fix than the nolint
default: it also closes the unbounded-allocation gap the "check after
converting" shape left open. Pick `maxDocSize` from a real domain limit (a
format spec's own max, a protocol's own frame cap) — never an arbitrary
round number.

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

## G101: Environment Variable Names and Enum/Const Identifiers

G101's identifier-name heuristic also fires on constants whose Go *name*
merely contains a credential-flagged substring (`secret`, `cred`, `apikey`,
...) even though the *value* is not a secret at all - an environment variable
name to read at runtime, or a plain enum tag:

```go
const (
	EnvAPIKey = "POSTMAN_API_KEY" //nolint:gosec // G101: This is an environment variable name, not a credential
)

const (
	SecretTypeOriginTeamRegex SecretTypeOrigin = "TEAM_REGEX" //nolint:gosec // G101: Enum/constant identifier matches the credential-name heuristic, but the value is a public tag, not a secret
)
```

Use `gosec.CommonReasons.EnvVarName` and `gosec.CommonReasons.EnumTagNotCredential`
for these. gosec's match is per-identifier, not per-`const` block, so only
annotate the specific line(s) it actually flags - a sibling constant in the
same block often isn't flagged at all.

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

## dupl: Structurally Parallel Generated-Client Wrappers

dupl's default remediation - extract a shared helper - is usually right. But
one shape recurs across generated-client wrappers (ogen, openapi-generator,
protoc): sibling methods per resource kind (`CreateFolder` / `CreateRequest` /
`CreateResponse`, or `GetX` / `DeleteX` repeated per `X`) that each switch over
a *distinct*, codegen-produced response/error union:

```go
func (s *Service) GetFolder(ctx context.Context, collectionID, folderID string, opts *GetOptions) (*FolderResult, error) {
	// ...
	switch r := res.(type) {
	case *api.CollectionFolderInfo:
		// ...
	case *api.GetCollectionFolderNotFound:
		return nil, postmanerr.FromProblemDetails([]byte(*r), http.StatusNotFound)
	// ...
	}
}

//nolint:dupl // Structurally parallel to sibling wrapper methods over distinct generated types; not meaningfully extractable without reflection or per-type adapters
func (s *Service) GetRequest(ctx context.Context, collectionID, requestID string, opts *GetOptions) (*RequestResult, error) {
	// ...
	switch r := res.(type) {
	case *api.CollectionRequestInfo:  // <- unrelated type to CollectionFolderInfo
		// ...
	case *api.GetCollectionRequestNotFound:  // <- unrelated type to GetCollectionFolderNotFound
		return nil, postmanerr.FromProblemDetails([]byte(*r), http.StatusNotFound)
	// ...
	}
}
```

`CollectionFolderInfo` and `CollectionRequestInfo` share no common interface -
a real extraction needs reflection or a per-type adapter layer, which is
harder to follow than the duplication it removes. Use `dupl.Nolint` from the
`dupl` subpackage:

```go
import "github.com/grokify/mogo/lintfix/dupl"

comment := dupl.Nolint(dupl.CommonReasons.ParallelResourceWrapper)
// "//nolint:dupl // Structurally parallel to sibling wrapper methods over
// distinct generated types; not meaningfully extractable without reflection
// or per-type adapters"
```

The same reasoning applies to test files: standalone, one-test-per-endpoint
httptest cases are usually clearer than a table-driven consolidation forced
just to satisfy dupl. Use `dupl.CommonReasons.StandaloneTestClarity` there.

**Reach for the real refactor first** when the duplicated blocks operate on
the *same* concrete type, or the difference is a single value trivial to lift
into a function parameter - see `remediations.json`'s `dupl.duplicate` entry
for the general case.

## unparam: Unused Parameters and Results

unparam (mvdan.cc/unparam) flags a function parameter (or result) that never
actually varies across its call sites - most often leftover generality from
an earlier version of the function, and especially common in test helpers as
call sites accumulate over time:

```go
// unparam: category always receives ClaimStatistical
func verifiedClaim(id string, category ClaimCategory) Claim {
    return Claim{ID: id, Category: category}
}
```

**The default remediation - delete the parameter, hardcode the constant - is
almost always right** in unexported code, and is a real simplification, not
just linter appeasement:

```go
func verifiedClaim(id string) Claim {
    return Claim{ID: id, Category: ClaimStatistical}
}
```

**Exception: the signature is constrained by something other than its own
call sites.** unparam only sees call sites within the analyzed code - it
can't see that a signature is fixed by an interface method set, a
function-type variable (`http.HandlerFunc`, `sort.Interface`, a callback
struct field), or an exported API whose signature is a compatibility
contract. Deleting the parameter there isn't possible (or isn't safe)
without breaking the thing the signature exists to satisfy - use `nolint`:

```go
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) { //nolint:unparam // r required by http.Handler
    w.WriteHeader(http.StatusOK)
}
```

Use the `unparam` subpackage to generate the comment:

```go
import "github.com/grokify/mogo/lintfix/unparam"

comment := unparam.Nolint(unparam.CommonReasons.InterfaceSignature)
// "//nolint:unparam // Signature fixed by an interface method set this type implements"
```

## Nolint Generators

The `gosec` subpackage provides type-safe nolint comment generators:

```go
gosec.NolintG101(reason)  // Hardcoded credentials (false positive)
gosec.NolintG115(reason)  // Integer overflow (bounded value)
gosec.NolintG117(reason)  // Secret in JSON response
gosec.NolintG118(reason)  // context.Background in goroutine
gosec.NolintG122(reason)  // filepath.Walk TOCTOU race (cmd/ entry point only)
gosec.NolintG404(reason)  // Weak random number generator (non-security use only; use crypto/rand otherwise)
gosec.NolintG124(reason)  // Insecure cookie attributes (set dynamically/from config)
gosec.NolintG703(reason)  // Path traversal (CLI entry point only)
gosec.NolintG704(reason)  // SSRF (trusted URL)
gosec.NolintG705(reason)  // XSS (trusted content)
gosec.NolintG706(reason)  // Log injection (prefer the strconv.Quote code fix instead)
gosec.NolintG710(reason)  // Open redirect (prefer the url.URL{} code fix instead)
```

The `dupl` subpackage provides the equivalent for duplicate-code findings:

```go
dupl.Nolint(reason)  // Structurally-required duplication (see "dupl" section above)
```

The `unparam` subpackage provides the equivalent for unused parameter/result findings:

```go
unparam.Nolint(reason)  // Signature constrained by an interface, callback, or exported API (see "unparam" section above)
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
gosec.CommonReasons.EnvVarName                // G101 - environment variable name, not a credential
gosec.CommonReasons.EnumTagNotCredential       // G101 - enum/const identifier matches heuristic, value is a public tag
gosec.CommonReasons.TestControlledInputNoUntrustedSource // G706 - nolint fallback only; prefer strconv.Quote
gosec.CommonReasons.ShufflingDisplayData       // G404 - non-security use only; use crypto/rand otherwise

dupl.CommonReasons.ParallelResourceWrapper    // sibling wrapper methods over distinct generated types
dupl.CommonReasons.StandaloneTestClarity      // standalone per-endpoint test, not worth consolidating

unparam.CommonReasons.InterfaceSignature      // parameter required by an interface method set
unparam.CommonReasons.CallbackSignature       // parameter required by a callback/function-type value
unparam.CommonReasons.ExportedAPICompat       // parameter kept for exported API compatibility
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
