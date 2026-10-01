// Copyright 2026 John Wang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package dupl provides helpers for generating nolint comments for the dupl
// linter (github.com/mibk/dupl), which flags structurally similar code blocks.
//
// dupl's default remediation is "extract a shared helper function" - and that
// is usually right. But one shape recurs across generated-client wrappers
// (ogen, openapi-generator, protoc, etc.): sibling methods that each wrap a
// distinct, codegen-produced type for a parallel resource (e.g. CreateFolder /
// CreateRequest / CreateResponse, each switching over its own
// *api.CreateFolderRes / *api.CreateRequestRes union with its own per-status
// error types). The method bodies read as duplicates to dupl, but the
// generated types share no common interface, so a real extraction would need
// reflection or a per-type adapter - trading a few duplicated lines for a
// genuinely harder-to-follow abstraction. In that case, nolint is the correct
// call, not a forced refactor.
//
// # Usage
//
//	comment := dupl.Nolint(dupl.CommonReasons.ParallelResourceWrapper)
//	// Returns: "//nolint:dupl // Structurally parallel to sibling wrapper
//	// methods over distinct generated types; not meaningfully extractable
//	// without reflection or per-type adapters"
//
// # When to refactor instead
//
// Reach for a shared helper first when the duplicated blocks operate on the
// same concrete type, or when the difference is a single parameter that's
// trivial to lift into a function argument. Reserve nolint for cases where
// the duplication is inherent to the shape being wrapped - see
// CommonReasons for the two most common ones.
package dupl
