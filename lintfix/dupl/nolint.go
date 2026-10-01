// Copyright 2026 John Wang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package dupl

import "fmt"

// Nolint formats a nolint:dupl comment with the given reason.
//
// Example:
//
//	comment := dupl.Nolint("standalone test for a distinct endpoint")
//	// Returns: "//nolint:dupl // standalone test for a distinct endpoint"
func Nolint(reason string) string {
	return fmt.Sprintf("//nolint:dupl // %s", reason)
}

// CommonReasons provides pre-written reason strings for common scenarios.
var CommonReasons = struct {
	// ParallelResourceWrapper is for source code: sibling methods that wrap
	// distinct generated types for parallel resources (e.g. a Create/Get/
	// Update/Delete quartet repeated per resource kind in an ogen/
	// openapi-generator client). The bodies look identical but switch over
	// unrelated per-operation response and error types, so there is no
	// common interface to abstract over without reflection or a per-type
	// adapter layer.
	ParallelResourceWrapper string

	// StandaloneTestClarity is for test code: each test is a self-contained,
	// independently readable case for a distinct endpoint or input. Table-
	// driving or extracting a shared subtest helper would save a few lines
	// at the cost of making each test harder to read in isolation and
	// breaking from the rest of the suite's per-endpoint test convention.
	StandaloneTestClarity string
}{
	ParallelResourceWrapper: "Structurally parallel to sibling wrapper methods over distinct generated types; not meaningfully extractable without reflection or per-type adapters",
	StandaloneTestClarity:   "Standalone test for a distinct endpoint; consolidating would reduce per-test clarity for a negligible reduction in duplication",
}
