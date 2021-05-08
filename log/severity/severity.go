package logutil

import (
	"fmt"
	"strings"
)

const (
	SeverityDisabled      = "disabled"
	SeverityEmergency     = "emergency"
	SeverityAlert         = "alert"
	SeverityCritical      = "critical"
	SeverityError         = "error"
	SeverityWarning       = "warning"
	SeverityNotice        = "notice"
	SeverityInformational = "informational"
	SeverityDebug         = "debug"
)

var mapStringSeverity = map[string]string{
	"disabled":      SeverityDisabled,
	"emergency":     SeverityEmergency,
	"alert":         SeverityAlert,
	"critical":      SeverityCritical,
	"crit":          SeverityCritical,
	"error":         SeverityError,
	"err":           SeverityError,
	"warning":       SeverityWarning,
	"warn":          SeverityWarning,
	"notice":        SeverityNotice,
	"informational": SeverityInformational,
	"info":          SeverityInformational,
	"debug":         SeverityDebug,
}

var severities = map[string]int{
	SeverityDisabled:      -1,
	SeverityEmergency:     0,
	SeverityAlert:         2,
	SeverityCritical:      3,
	SeverityError:         4,
	SeverityWarning:       5,
	SeverityNotice:        6,
	SeverityInformational: 7,
	SeverityDebug:         8,
}

// ParseSeverity takes a string and returns a constant
// `Severity` value. Common aliases such as `warn` and
// `info` are included.
func ParseSeverity(sev string) (string, error) {
	sev = strings.ToLower(strings.TrimSpace(sev))
	sev2, ok := mapStringSeverity[sev]
	if ok {
		return sev2, nil
	}
	return SeverityDisabled, fmt.Errorf("severity not found [%s]", sev)
}

// SeverityInclude checks to see if a severity level
// is included against a certain severity filter.
func SeverityInclude(filterLevel, itemLevel string) (bool, error) {
	filterLevelInt, ok := severities[filterLevel]
	if !ok {
		return false, fmt.Errorf("filterLevel [%s] not supported", filterLevel)
	}
	if filterLevelInt < 1 {
		return false, nil
	}
	itemLevelInt, ok := severities[itemLevel]
	if !ok {
		return false, fmt.Errorf("itemLevel [%s] not supported", itemLevel)
	}
	if filterLevel == SeverityDisabled || itemLevel == SeverityDisabled {
		return false, nil
	}
	if itemLevelInt <= filterLevelInt {
		return true, nil
	}
	return false, nil
}
