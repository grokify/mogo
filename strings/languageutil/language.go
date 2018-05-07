package languageutil

import (
	"fmt"
	"strings"

	"golang.org/x/text/language"
)

var languageConjunctionMap = map[language.Tag]string{
	language.English: "and",
}

func JoinLanguage(slice []string, sep string, joinLang language.Tag) (string, error) {
	switch len(slice) {
	case 0:
		return "", nil
	case 1:
		return slice[0], nil
	case 2:
		if joinWord, ok := languageConjunctionMap[joinLang]; ok {
			return slice[0] + " " + joinWord + " " + slice[1], nil
		}
		return strings.Join(slice, sep), fmt.Errorf("Join word not found for language [%v]", joinLang)
	default:
		last, rest := slice[len(slice)-1], slice[:len(slice)-1]
		if joinWord, ok := languageConjunctionMap[joinLang]; ok {
			rest = append(rest, joinWord+" "+last)
			return strings.Join(rest, sep+" "), nil
		}
		return strings.Join(slice, sep), fmt.Errorf("Join word not found for language [%v]", joinLang)
	}
}
