// Copyright 2026 John Wang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package lintfix provides a structured database of lint rule remediations
// for Go projects using golangci-lint.
//
// This package serves as a "data overlay" that maps lint errors to:
//   - Remediation strategies (code fix, nolint annotation, refactor)
//   - Helper packages that provide actual fixes (within mogo)
//   - Pre-written nolint comments with proper documentation
//   - Example code and explanations
//
// # Usage
//
// Load the remediation database and query for specific rules:
//
//	db := lintfix.MustLoadRemediations()
//	fix := db.GetGosec("G120")
//	fmt.Println(fix.Remediation.Summary)
//	// "Use http.MaxBytesReader before parsing form data"
//
// # Remediation Types
//
// The database categorizes remediations into three types:
//
//   - "code": Fix by adding/changing code (e.g., LimitRequestBody for G120)
//   - "nolint": Fix by adding a nolint annotation with proper documentation
//   - "refactor": Fix requires broader code changes (e.g., removing hardcoded secrets)
//
// # Nolint Generators
//
// For rules that require nolint annotations, use the gosec subpackage:
//
//	comment := gosec.NolintG117(gosec.CommonReasons.OAuthTokenResponse)
//	// Returns: "//nolint:gosec // G117: OAuth token response per RFC 6749"
//
// # Helper Package References
//
// Code-based remediations reference helper packages within mogo:
//
//	fix := db.GetGosec("G120")
//	fmt.Println(fix.Remediation.Package)
//	// "github.com/grokify/mogo/net/http/httputilmore"
//	fmt.Println(fix.Remediation.Function)
//	// "LimitRequestBody"
//
// # Supported Linters
//
// Currently supported:
//   - gosec: Security-focused linter
//   - staticcheck: Go static analysis
//   - errcheck: Error handling checks
//
// # Documentation
//
// For detailed guides including version-specific caveats, see:
// https://github.com/grokify/mogo/tree/main/docs/lintfix
package lintfix
