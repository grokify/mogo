# Lintfix Documentation

The `lintfix` package provides structured lint remediation guidance for Go projects.

## Contents

- [Gosec Version Caveats](gosec-caveats.md) - Version-specific behaviors that affect remediation strategies

## Overview

When golangci-lint flags an issue, `lintfix` helps you determine the appropriate fix:

1. **Query the database** to understand the rule and get remediation guidance
2. **Choose the fix type** based on your situation (code fix, nolint, or refactor)
3. **Apply the fix** using helper functions or generated nolint comments

## Remediation Decision Tree

```
Lint issue flagged
    │
    ├─ Is it a false positive?
    │   └─ Yes → Use nolint with documented reason
    │
    ├─ Is it intentional behavior?
    │   └─ Yes → Use nolint with documented reason
    │
    ├─ Can it be fixed with a helper function?
    │   └─ Yes → Use the referenced mogo helper
    │
    └─ Does it require broader changes?
        └─ Yes → Plan refactoring (move secrets, update crypto, etc.)
```

## Package Structure

```
mogo/lintfix/
├── remediations.go      # Database loading and query API
├── remediations.json    # Embedded remediation database
└── gosec/
    └── nolint.go        # Nolint comment generators
```

## Related Packages

Helper functions referenced by code remediations:

| Rule | Package | Function |
|------|---------|----------|
| G120 | `mogo/net/http/httputilmore` | `LimitRequestBody` |
| G703 | `mogo/os/osutil` | `ReadFileSecure`, `WriteFileSecure`, `CopyFileSecure` |
| G122 | `mogo/os/osutil` | `ReadDirFilesSecure` |

## Contributing

To add support for new lint rules:

1. Add the rule to `remediations.json`
2. If it needs a nolint generator, add to `gosec/nolint.go`
3. If it needs common reasons, add to `gosec.CommonReasons`
4. Document any caveats in `docs/lintfix/`
