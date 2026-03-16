// Copyright 2026 John Wang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package lintfix

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed remediations.json
var remediationsJSON []byte

// RemediationDB is the top-level structure for the remediation database.
type RemediationDB struct {
	Version     string                         `json:"version"`
	Description string                         `json:"description"`
	Linters     map[string]map[string]*RuleFix `json:"linters"`
}

// RuleFix contains remediation information for a specific lint rule.
type RuleFix struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Severity    string       `json:"severity,omitempty"`
	Category    string       `json:"category,omitempty"`
	Remediation *Remediation `json:"remediation"`
	References  []string     `json:"references,omitempty"`
}

// Remediation contains the actual fix information.
type Remediation struct {
	Type        string   `json:"type"` // "code", "nolint", "refactor"
	Summary     string   `json:"summary"`
	Pattern     string   `json:"pattern,omitempty"`
	Package     string   `json:"package,omitempty"`
	Function    string   `json:"function,omitempty"`
	Example     string   `json:"example,omitempty"`
	Explanation string   `json:"explanation,omitempty"`
	When        string   `json:"when,omitempty"`
	Avoid       []string `json:"avoid,omitempty"`
	Caveats     []string `json:"caveats,omitempty"`
}

// LoadRemediations loads and parses the embedded remediation database.
func LoadRemediations() (*RemediationDB, error) {
	var db RemediationDB
	if err := json.Unmarshal(remediationsJSON, &db); err != nil {
		return nil, fmt.Errorf("failed to parse remediations.json: %w", err)
	}
	return &db, nil
}

// MustLoadRemediations loads the remediation database or panics.
func MustLoadRemediations() *RemediationDB {
	db, err := LoadRemediations()
	if err != nil {
		panic(err)
	}
	return db
}

// Get retrieves a remediation by linter and rule code.
// Returns nil if not found.
func (db *RemediationDB) Get(linter, code string) *RuleFix {
	if db.Linters == nil {
		return nil
	}
	rules, ok := db.Linters[linter]
	if !ok {
		return nil
	}
	return rules[code]
}

// GetGosec is a convenience method for getting gosec remediations.
func (db *RemediationDB) GetGosec(code string) *RuleFix {
	return db.Get("gosec", code)
}

// GetStaticcheck is a convenience method for getting staticcheck remediations.
func (db *RemediationDB) GetStaticcheck(code string) *RuleFix {
	return db.Get("staticcheck", code)
}

// ListLinters returns all linters in the database.
func (db *RemediationDB) ListLinters() []string {
	linters := make([]string, 0, len(db.Linters))
	for linter := range db.Linters {
		linters = append(linters, linter)
	}
	return linters
}

// ListRules returns all rule codes for a given linter.
func (db *RemediationDB) ListRules(linter string) []string {
	rules, ok := db.Linters[linter]
	if !ok {
		return nil
	}
	codes := make([]string, 0, len(rules))
	for code := range rules {
		codes = append(codes, code)
	}
	return codes
}

// String returns a formatted description of the rule fix.
func (rf *RuleFix) String() string {
	if rf == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s: %s", rf.Name, rf.Description)
}

// HasHelper returns true if this remediation has a helper function.
func (rf *RuleFix) HasHelper() bool {
	if rf == nil || rf.Remediation == nil {
		return false
	}
	return rf.Remediation.Package != "" && rf.Remediation.Function != ""
}
