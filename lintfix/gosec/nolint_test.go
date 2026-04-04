// Copyright 2026 John Wang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gosec

import (
	"strings"
	"testing"
)

func TestNolint(t *testing.T) {
	tests := []struct {
		name   string
		rule   string
		reason string
		want   string
	}{
		{
			name:   "G117 OAuth",
			rule:   "G117",
			reason: "OAuth token response per RFC 6749",
			want:   "//nolint:gosec // G117: OAuth token response per RFC 6749",
		},
		{
			name:   "G118 shutdown",
			rule:   "G118",
			reason: "Shutdown handler runs after request context is cancelled",
			want:   "//nolint:gosec // G118: Shutdown handler runs after request context is cancelled",
		},
		{
			name:   "G703 path validation",
			rule:   "G703",
			reason: "Input validated to reject path separators",
			want:   "//nolint:gosec // G703: Input validated to reject path separators",
		},
		{
			name:   "G704 test",
			rule:   "G704",
			reason: "Test uses httptest server URL",
			want:   "//nolint:gosec // G704: Test uses httptest server URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Nolint(tt.rule, tt.reason)
			if got != tt.want {
				t.Errorf("Nolint(%q, %q) = %q, want %q", tt.rule, tt.reason, got, tt.want)
			}
		})
	}
}

func TestNolintG101(t *testing.T) {
	got := NolintG101("URL path, not a credential")
	if !strings.Contains(got, "G101") {
		t.Errorf("NolintG101() = %q, does not contain G101", got)
	}
	if !strings.Contains(got, "//nolint:gosec") {
		t.Errorf("NolintG101() = %q, does not contain //nolint:gosec", got)
	}
}

func TestNolintG115(t *testing.T) {
	got := NolintG115("Value bounded by prior validation")
	if !strings.Contains(got, "G115") {
		t.Errorf("NolintG115() = %q, does not contain G115", got)
	}
	if !strings.Contains(got, "//nolint:gosec") {
		t.Errorf("NolintG115() = %q, does not contain //nolint:gosec", got)
	}
}

func TestNolintG117(t *testing.T) {
	got := NolintG117("OAuth token response per RFC 6749")
	if !strings.Contains(got, "G117") {
		t.Errorf("NolintG117() = %q, does not contain G117", got)
	}
	if !strings.Contains(got, "//nolint:gosec") {
		t.Errorf("NolintG117() = %q, does not contain //nolint:gosec", got)
	}
}

func TestNolintG118(t *testing.T) {
	got := NolintG118("Shutdown handler runs after request context is cancelled")
	if !strings.Contains(got, "G118") {
		t.Errorf("NolintG118() = %q, does not contain G118", got)
	}
	if !strings.Contains(got, "//nolint:gosec") {
		t.Errorf("NolintG118() = %q, does not contain //nolint:gosec", got)
	}
}

func TestNolintG703(t *testing.T) {
	got := NolintG703("Input validated to reject path separators")
	if !strings.Contains(got, "G703") {
		t.Errorf("NolintG703() = %q, does not contain G703", got)
	}
	if !strings.Contains(got, "//nolint:gosec") {
		t.Errorf("NolintG703() = %q, does not contain //nolint:gosec", got)
	}
}

func TestNolintG704(t *testing.T) {
	got := NolintG704("Test uses httptest server URL")
	if !strings.Contains(got, "G704") {
		t.Errorf("NolintG704() = %q, does not contain G704", got)
	}
	if !strings.Contains(got, "//nolint:gosec") {
		t.Errorf("NolintG704() = %q, does not contain //nolint:gosec", got)
	}
}

func TestCommonReasons(t *testing.T) {
	// G101 reasons
	if CommonReasons.URLPathNotCredential == "" {
		t.Error("CommonReasons.URLPathNotCredential is empty")
	}
	if CommonReasons.TestFixture == "" {
		t.Error("CommonReasons.TestFixture is empty")
	}
	if CommonReasons.DocumentationExample == "" {
		t.Error("CommonReasons.DocumentationExample is empty")
	}

	// G115 reasons
	if CommonReasons.BoundedByValidation == "" {
		t.Error("CommonReasons.BoundedByValidation is empty")
	}
	if CommonReasons.DomainConstraint == "" {
		t.Error("CommonReasons.DomainConstraint is empty")
	}
	if CommonReasons.YearValueSafeRange == "" {
		t.Error("CommonReasons.YearValueSafeRange is empty")
	}
	if CommonReasons.SmallEnumValue == "" {
		t.Error("CommonReasons.SmallEnumValue is empty")
	}

	// G117 reasons
	if CommonReasons.OAuthTokenResponse == "" {
		t.Error("CommonReasons.OAuthTokenResponse is empty")
	}
	if CommonReasons.OAuthRegistrationResponse == "" {
		t.Error("CommonReasons.OAuthRegistrationResponse is empty")
	}

	// G118 reasons
	if CommonReasons.ShutdownHandler == "" {
		t.Error("CommonReasons.ShutdownHandler is empty")
	}
	if CommonReasons.BackgroundJob == "" {
		t.Error("CommonReasons.BackgroundJob is empty")
	}
	if CommonReasons.CleanupRoutine == "" {
		t.Error("CommonReasons.CleanupRoutine is empty")
	}
	if CommonReasons.IndependentCancel == "" {
		t.Error("CommonReasons.IndependentCancel is empty")
	}

	// G703 reasons (use only in cmd/, not library code)
	if CommonReasons.PathFromCLIFlag == "" {
		t.Error("CommonReasons.PathFromCLIFlag is empty")
	}
	if CommonReasons.PathFromConfig == "" {
		t.Error("CommonReasons.PathFromConfig is empty")
	}
	if CommonReasons.OutputPathByUser == "" {
		t.Error("CommonReasons.OutputPathByUser is empty")
	}

	// G704 reasons
	if CommonReasons.HttptestServer == "" {
		t.Error("CommonReasons.HttptestServer is empty")
	}
	if CommonReasons.ValidatedAllowlist == "" {
		t.Error("CommonReasons.ValidatedAllowlist is empty")
	}
	if CommonReasons.InternalServiceURL == "" {
		t.Error("CommonReasons.InternalServiceURL is empty")
	}
	if CommonReasons.TrustedConstantsURL == "" {
		t.Error("CommonReasons.TrustedConstantsURL is empty")
	}
}

func TestNolintWithCommonReasons(t *testing.T) {
	// Test that Nolint* functions work with CommonReasons
	tests := []struct {
		name string
		got  string
	}{
		{
			name: "G101 with URLPathNotCredential",
			got:  NolintG101(CommonReasons.URLPathNotCredential),
		},
		{
			name: "G115 with BoundedByValidation",
			got:  NolintG115(CommonReasons.BoundedByValidation),
		},
		{
			name: "G117 with OAuthTokenResponse",
			got:  NolintG117(CommonReasons.OAuthTokenResponse),
		},
		{
			name: "G118 with ShutdownHandler",
			got:  NolintG118(CommonReasons.ShutdownHandler),
		},
		{
			name: "G703 with PathFromCLIFlag",
			got:  NolintG703(CommonReasons.PathFromCLIFlag),
		},
		{
			name: "G704 with HttptestServer",
			got:  NolintG704(CommonReasons.HttptestServer),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.HasPrefix(tt.got, "//nolint:gosec") {
				t.Errorf("got = %q, does not start with //nolint:gosec", tt.got)
			}
		})
	}
}
