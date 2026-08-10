// Copyright 2026 John Wang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package unparam

import "fmt"

// Nolint formats a nolint:unparam comment with the given reason.
//
// Example:
//
//	comment := unparam.Nolint("r required by http.Handler")
//	// Returns: "//nolint:unparam // r required by http.Handler"
func Nolint(reason string) string {
	return fmt.Sprintf("//nolint:unparam // %s", reason)
}

// CommonReasons provides pre-written reason strings for common scenarios
// where a parameter or result unparam flags as unused is actually required
// by something outside the function's own call sites.
var CommonReasons = struct {
	// InterfaceSignature is for a method whose parameter list is fixed by
	// an interface it implements (e.g. sort.Interface, io.Writer,
	// database/sql/driver.Valuer) even though this particular
	// implementation doesn't use every parameter.
	InterfaceSignature string

	// CallbackSignature is for a function passed as a value where the
	// signature is fixed by the receiving API (http.HandlerFunc, a
	// third-party callback/hook type, a function-type struct field) rather
	// than chosen by this function itself.
	CallbackSignature string

	// ExportedAPICompat is for an exported function or method whose
	// signature is a compatibility contract: removing a currently-unused
	// parameter would be a breaking change for callers outside this
	// module, even though no internal call site varies it today.
	ExportedAPICompat string
}{
	InterfaceSignature: "Signature fixed by an interface method set this type implements",
	CallbackSignature:  "Signature fixed by the callback/function type this value is passed as",
	ExportedAPICompat:  "Parameter kept for exported API compatibility; removing it would be a breaking change for external callers",
}
