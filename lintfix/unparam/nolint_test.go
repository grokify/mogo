// Copyright 2026 John Wang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package unparam

import "testing"

func TestNolint(t *testing.T) {
	got := Nolint("custom reason")
	want := "//nolint:unparam // custom reason"
	if got != want {
		t.Errorf("Nolint() = %q, want %q", got, want)
	}
}

func TestNolintWithCommonReasons(t *testing.T) {
	tests := []struct {
		name   string
		reason string
	}{
		{"InterfaceSignature", CommonReasons.InterfaceSignature},
		{"CallbackSignature", CommonReasons.CallbackSignature},
		{"ExportedAPICompat", CommonReasons.ExportedAPICompat},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.reason == "" {
				t.Fatalf("CommonReasons.%s is empty", tt.name)
			}
			got := Nolint(tt.reason)
			want := "//nolint:unparam // " + tt.reason
			if got != want {
				t.Errorf("Nolint(%s) = %q, want %q", tt.name, got, want)
			}
		})
	}
}
