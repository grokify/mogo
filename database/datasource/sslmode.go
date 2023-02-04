package datasource

import (
	"errors"
	"strings"
)

const (
	OracleParamConnectTimeout = "connect_timeout"

	PgSSLModeParam    = "sslmode"
	SSLModeAllow      = "allow"
	SSLModeDisable    = "disable"
	SSLModePrefer     = "prefer"
	SSLModeRequire    = "require"
	SSLModeVerifyCA   = "verify-ca"
	SSLModeVerifyFull = "verify-full"
	SSLModeDefault    = SSLModeDisable
)

var ErrSSLModeNotSUpported = errors.New("sslmode not supported")

// SSLModeParseOrDefault manages Postgres sslmode query param
func SSLModeParseOrDefault(s, d string) string {
	m, err := SSLModeParse(s)
	if err != nil {
		return d
	}
	return m
}

func sslModes() map[string]int {
	return map[string]int{
		SSLModeAllow:      1,
		SSLModeDisable:    1,
		SSLModePrefer:     1,
		SSLModeRequire:    1,
		SSLModeVerifyCA:   1,
		SSLModeVerifyFull: 1}
}

// SSLModeParse parses Postgres sslmode query param
func SSLModeParse(s string) (string, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	modes := sslModes()
	if _, ok := modes[s]; ok {
		return s, nil
	}
	return "", ErrSSLModeNotSUpported
}
