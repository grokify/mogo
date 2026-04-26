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

- **gosec** - Security-focused rules (G101, G115, G117, G118, G120, G401, G501, G601, G703, G704)
- **staticcheck** - Static analysis (SA1019, SA4006)
- **errcheck** - Error handling

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

## Nolint Generators

The `gosec` subpackage provides type-safe nolint comment generators:

```go
gosec.NolintG101(reason)  // Hardcoded credentials (false positive)
gosec.NolintG115(reason)  // Integer overflow (bounded value)
gosec.NolintG117(reason)  // Secret in JSON response
gosec.NolintG118(reason)  // context.Background in goroutine
gosec.NolintG703(reason)  // Path traversal (CLI entry point only)
gosec.NolintG704(reason)  // SSRF (trusted URL)
```

### Common Reasons

Pre-written reason strings for common scenarios:

```go
gosec.CommonReasons.OAuthTokenResponse        // G117
gosec.CommonReasons.ShutdownHandler           // G118
gosec.CommonReasons.InputValidatedNoPathSep   // G703
gosec.CommonReasons.HttptestServer            // G704
gosec.CommonReasons.BoundedByValidation       // G115
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
