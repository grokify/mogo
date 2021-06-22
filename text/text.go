package text

import (
	"fmt"
)

type Text struct {
	Display string
	Slug    string
}

type Texts []Text

func (texts Texts) DisplayForSlug(slug string) (string, error) {
	for _, try := range texts {
		if slug == try.Slug {
			return try.Display, nil
		}
	}
	return "", fmt.Errorf("slug not found [%s]", slug)
}

func (texts Texts) DisplayTexts() []string {
	displays := []string{}
	for _, txt := range texts {
		displays = append(displays, txt.Display)
	}
	return displays
}

/*
type TextSet struct {
	Texts []Text
}

func (ts *TextSet) DisplayTexts() []string {
	displays := []string{}
	for _, txt := range ts.Texts {
		displays = append(displays, txt.Display)
	}
	return displays
}
*/
