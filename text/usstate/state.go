package usstate

import "strings"

func Abbreviate(stateName string) string {
	stateNameLc := strings.ToLower(strings.TrimSpace(stateName))
	for abbr, try := range usc {
		if strings.ToLower(try) == stateNameLc {
			return abbr
		}
	}
	return stateName
}
