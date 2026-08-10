// Copyright 2026 John Wang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package unparam provides helpers for generating nolint comments for the
// unparam linter (mvdan.cc/unparam), which flags a function parameter (or
// result) that never varies across its call sites.
//
// unparam's default remediation - delete the parameter and inline the
// constant value at its one effective use - is usually right, and is a real
// simplification rather than linter appeasement. But unparam cannot see that
// a signature is fixed by something other than its own call sites: an
// interface method set, a function-type variable (http.HandlerFunc,
// sort.Interface, a callback field), or an exported API whose signature is a
// compatibility contract. In those cases the parameter is genuinely unused
// today, but removing it isn't possible (or isn't safe) without breaking the
// thing the signature exists to satisfy - nolint is the correct call.
//
// # Usage
//
//	comment := unparam.Nolint(unparam.CommonReasons.InterfaceSignature)
//	// Returns: "//nolint:unparam // Signature fixed by an interface method
//	// set this type implements"
//
// # When to fix instead of nolint
//
// Reach for the real fix first - delete the parameter and hardcode the
// constant inside the function body - whenever the signature isn't
// externally constrained. This is especially common in test helpers, where
// a parameter added for generality often ends up receiving the same value
// at every call site as the test suite grows. See CommonReasons for the
// constrained-signature cases nolint is actually for.
package unparam
