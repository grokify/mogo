// Copyright 2026 John Wang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package lintfix

import (
	"testing"
)

func TestLoadRemediations(t *testing.T) {
	db, err := LoadRemediations()
	if err != nil {
		t.Fatalf("LoadRemediations() error = %v", err)
	}
	if db == nil {
		t.Fatal("LoadRemediations() returned nil")
	}
	if db.Version == "" {
		t.Error("LoadRemediations() returned empty version")
	}
}

func TestMustLoadRemediations(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MustLoadRemediations() panicked: %v", r)
		}
	}()

	db := MustLoadRemediations()
	if db == nil {
		t.Fatal("MustLoadRemediations() returned nil")
	}
}

func TestRemediationDB_Get(t *testing.T) {
	db := MustLoadRemediations()

	tests := []struct {
		name     string
		linter   string
		code     string
		wantName string
		wantNil  bool
	}{
		{
			name:     "gosec G120",
			linter:   "gosec",
			code:     "G120",
			wantName: "Unbounded request body",
		},
		{
			name:     "gosec G117",
			linter:   "gosec",
			code:     "G117",
			wantName: "Secret in JSON response",
		},
		{
			name:     "gosec G118",
			linter:   "gosec",
			code:     "G118",
			wantName: "context.Background in goroutine",
		},
		{
			name:     "gosec G704",
			linter:   "gosec",
			code:     "G704",
			wantName: "SSRF via taint analysis",
		},
		{
			name:     "staticcheck SA1019",
			linter:   "staticcheck",
			code:     "SA1019",
			wantName: "Deprecated API usage",
		},
		{
			name:    "unknown linter",
			linter:  "unknown",
			code:    "X999",
			wantNil: true,
		},
		{
			name:    "unknown rule",
			linter:  "gosec",
			code:    "G999",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := db.Get(tt.linter, tt.code)
			if tt.wantNil {
				if got != nil {
					t.Errorf("Get(%q, %q) = %v, want nil", tt.linter, tt.code, got)
				}
				return
			}
			if got == nil {
				t.Fatalf("Get(%q, %q) = nil, want %q", tt.linter, tt.code, tt.wantName)
			}
			if got.Name != tt.wantName {
				t.Errorf("Get(%q, %q).Name = %q, want %q", tt.linter, tt.code, got.Name, tt.wantName)
			}
		})
	}
}

func TestRemediationDB_GetGosec(t *testing.T) {
	db := MustLoadRemediations()

	got := db.GetGosec("G120")
	if got == nil {
		t.Fatal("GetGosec(G120) = nil, want non-nil")
	}
	if got.Name != "Unbounded request body" {
		t.Errorf("GetGosec(G120).Name = %q, want %q", got.Name, "Unbounded request body")
	}
}

func TestRemediationDB_GetStaticcheck(t *testing.T) {
	db := MustLoadRemediations()

	got := db.GetStaticcheck("SA1019")
	if got == nil {
		t.Fatal("GetStaticcheck(SA1019) = nil, want non-nil")
	}
	if got.Name != "Deprecated API usage" {
		t.Errorf("GetStaticcheck(SA1019).Name = %q, want %q", got.Name, "Deprecated API usage")
	}
}

func TestRemediationDB_ListLinters(t *testing.T) {
	db := MustLoadRemediations()

	linters := db.ListLinters()
	if len(linters) == 0 {
		t.Fatal("ListLinters() returned empty list")
	}

	// Check that gosec is in the list
	found := false
	for _, l := range linters {
		if l == "gosec" {
			found = true
			break
		}
	}
	if !found {
		t.Error("ListLinters() does not include gosec")
	}
}

func TestRemediationDB_ListRules(t *testing.T) {
	db := MustLoadRemediations()

	rules := db.ListRules("gosec")
	if len(rules) == 0 {
		t.Fatal("ListRules(gosec) returned empty list")
	}

	// Check for nil linter
	nilRules := db.ListRules("nonexistent")
	if nilRules != nil {
		t.Errorf("ListRules(nonexistent) = %v, want nil", nilRules)
	}
}

func TestRuleFix_String(t *testing.T) {
	db := MustLoadRemediations()

	rf := db.GetGosec("G120")
	if rf == nil {
		t.Fatal("GetGosec(G120) = nil")
	}

	s := rf.String()
	if s == "" || s == "<nil>" {
		t.Errorf("RuleFix.String() = %q, want non-empty", s)
	}

	// Test nil case
	var nilRF *RuleFix
	if nilRF.String() != "<nil>" {
		t.Errorf("nil.String() = %q, want %q", nilRF.String(), "<nil>")
	}
}

func TestRuleFix_HasHelper(t *testing.T) {
	db := MustLoadRemediations()

	// G120 has a helper function
	rf := db.GetGosec("G120")
	if rf == nil {
		t.Fatal("GetGosec(G120) = nil")
	}
	if !rf.HasHelper() {
		t.Error("G120.HasHelper() = false, want true")
	}

	// G117 is nolint-only, no helper
	rf117 := db.GetGosec("G117")
	if rf117 == nil {
		t.Fatal("GetGosec(G117) = nil")
	}
	if rf117.HasHelper() {
		t.Error("G117.HasHelper() = true, want false")
	}

	// Test nil case
	var nilRF *RuleFix
	if nilRF.HasHelper() {
		t.Error("nil.HasHelper() = true, want false")
	}
}

func TestRemediationDB_NilLinters(t *testing.T) {
	db := &RemediationDB{Linters: nil}
	got := db.Get("gosec", "G120")
	if got != nil {
		t.Errorf("Get with nil Linters = %v, want nil", got)
	}
}
