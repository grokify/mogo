// severity provides syslog-type severity level handling.
package severity

import (
	"fmt"
	"strings"
)

const (
	SeverityDisabled      = "disabled"
	SeverityEmergency     = "emerg"
	SeverityAlert         = "alert"
	SeverityCritical      = "crit"
	SeverityError         = "err"
	SeverityWarning       = "warning"
	SeverityNotice        = "notice"
	SeverityInformational = "info"
	SeverityDebug         = "debug"
)

var mapStringSeverity = map[string]string{
	"disabled":      SeverityDisabled,
	"disable":       SeverityDisabled,
	"off":           SeverityDisabled,
	"emergency":     SeverityEmergency,
	"emerg":         SeverityEmergency,
	"panic":         SeverityEmergency, // deprecated by syslog
	"exception":     SeverityEmergency, // used by PostgreSQL
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
	"hint":          SeverityDebug, // used by Spectral
}

var severities = map[string]int{
	SeverityDisabled:      -1,
	SeverityEmergency:     0,
	SeverityAlert:         1,
	SeverityCritical:      2,
	SeverityError:         3,
	SeverityWarning:       4,
	SeverityNotice:        5,
	SeverityInformational: 6,
	SeverityDebug:         7,
}

// Parse takes a string and returns a constant
// `Severity` value. Common aliases such as `warn` and
// `info` are included.
func Parse(sev string) (string, error) {
	sev = strings.ToLower(strings.TrimSpace(sev))
	sev2, ok := mapStringSeverity[sev]
	if ok {
		return sev2, nil
	}
	return SeverityDisabled, fmt.Errorf("severity not found [%s]", sev)
}

// Severities returns a list of severities.
func Severities() []string {
	return []string{
		SeverityDisabled,
		SeverityEmergency,
		SeverityAlert,
		SeverityCritical,
		SeverityError,
		SeverityWarning,
		SeverityNotice,
		SeverityInformational,
		SeverityDebug}
}

// SeverityInclude checks to see if a severity level
// is included against a certain severity filter.
func SeverityInclude(filterLevel, itemLevel string) (bool, error) {
	filterLevelGood, err := Parse(filterLevel)
	if err != nil || filterLevelGood == SeverityDisabled {
		return false, err
	}
	filterLevelInt, ok := severities[filterLevelGood]
	if !ok {
		return false, fmt.Errorf("filterLevel [%s] not supported", filterLevel)
	}
	if filterLevelInt < 1 {
		return false, nil
	}
	itemLevelGood, err := Parse(itemLevel)
	if err != nil || itemLevelGood == SeverityDisabled {
		return false, err
	}
	itemLevelInt, ok := severities[itemLevelGood]
	if !ok {
		return false, fmt.Errorf("itemLevel [%s] not supported", itemLevel)
	}
	if itemLevelInt > filterLevelInt {
		return false, nil
	}
	return true, nil
}
