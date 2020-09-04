package text

type Text struct {
	Display string
	Slug    string
}

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
