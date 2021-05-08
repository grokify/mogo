package severity

import "fmt"

type Severity string

const (
	SeverityDisabled      Severity = "disabled"
	SeverityEmergency     Severity = "emergency"
	SeverityAlert         Severity = "alert"
	SeverityCritical      Severity = "critical"
	SeverityError         Severity = "error"
	SeverityWarning       Severity = "warning"
	SeverityNotice        Severity = "notice"
	SeverityInformational Severity = "informational"
	SeverityDebug         Severity = "debug"
)

var severities = map[Severity]int{
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

func Parse(sev string) Severity {
	switch sev {
	case string(SeverityDisabled):
		return SeverityDisabled
	case string(SeverityEmergency):
		return SeverityEmergency
	case string(SeverityAlert):
		return SeverityAlert
	case string(SeverityCritical):
		return SeverityCritical
	case string(SeverityError):
		return SeverityError
	case string(SeverityWarning):
		return SeverityWarning
	case string(SeverityNotice):
		return SeverityNotice
	case string(SeverityInformational):
		return SeverityInformational
	case string(SeverityDebug):
		return SeverityDebug
	}
	return SeverityDisabled
}

func SeverityInclude(reportLevel, itemLevel Severity) (bool, error) {
	reportLevelInt, ok := severities[reportLevel]
	if !ok {
		return false, fmt.Errorf("reportLevel [%s] not supported", reportLevel)
	}
	if reportLevelInt < 1 {
		return false, nil
	}
	itemLevelInt, ok := severities[itemLevel]
	if !ok {
		return false, fmt.Errorf("itemLevel [%s] not supported", itemLevel)
	}
	if reportLevel == SeverityDisabled || itemLevel == SeverityDisabled {
		return false, nil
	}
	if itemLevelInt <= reportLevelInt {
		return true, nil
	}
	return false, nil
}
