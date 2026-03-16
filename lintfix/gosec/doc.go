// Copyright 2026 John Wang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package gosec provides helpers for generating nolint comments for gosec rules.
//
// This package generates properly formatted nolint:gosec comments with
// rule codes and documented reasons. Use these when a gosec warning is a
// false positive or when the flagged pattern is intentional.
//
// # Usage
//
// Generate a nolint comment:
//
//	comment := gosec.NolintG117(gosec.CommonReasons.OAuthTokenResponse)
//	// Returns: "//nolint:gosec // G117: OAuth token response per RFC 6749"
//
// Or use the generic Nolint function:
//
//	comment := gosec.Nolint("G117", "custom reason here")
//
// # Common Reasons
//
// The CommonReasons variable provides pre-written reason strings for
// common scenarios:
//
//	gosec.CommonReasons.OAuthTokenResponse        // G117
//	gosec.CommonReasons.ShutdownHandler           // G118
//	gosec.CommonReasons.HttptestServer            // G704
//	gosec.CommonReasons.URLPathNotCredential      // G101
//
// # Supported Rules
//
//   - G101: Hardcoded credentials (false positive for URL paths, test fixtures)
//   - G117: Secret in JSON response (OAuth tokens, registration responses)
//   - G118: context.Background in goroutine (shutdown handlers, background jobs)
//   - G703: Path traversal (validated input, trusted paths)
//   - G704: SSRF (trusted URLs, httptest servers, validated allowlists)
//
// # For Code Fixes
//
// For rules that require code changes (not nolint), see the helper packages:
//
//   - G120: github.com/grokify/mogo/net/http/httputilmore.LimitRequestBody
package gosec
