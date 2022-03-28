package datasource

import (
	"errors"
	"strings"
)

const (
	SchemePostgreSQL  = "postgres"
	SSLModeAllow      = "allow"
	SSLModeDisable    = "disable"
	SSLModePrefer     = "prefer"
	SSLModeRequire    = "require"
	SSLModeVerifyCA   = "verify-ca"
	SSLModeVerifyFull = "verify-full"
)

func SSLModeParseOrDefault(s, d string) string {
	m, err := SSLModeParse(s)
	if err != nil {
		m2, err := SSLModeParse(d)
		if err != nil {
			return ""
		}
		return m2
	}
	return m
}

func SSLModeParse(s string) (string, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case SSLModeAllow:
		return s, nil
	case SSLModeDisable:
		return s, nil
	case SSLModePrefer:
		return s, nil
	case SSLModeRequire:
		return s, nil
	case SSLModeVerifyCA:
		return s, nil
	case SSLModeVerifyFull:
		return s, nil
	}
	return "", errors.New("sslMode not supported")
}
