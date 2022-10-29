package htmlutil

import "html"

// Text represents a text string that fulfills the `Stringable` interface.
type Text struct {
	Text    string
	Escaped bool
}

func (s Text) String() (string, error) {
	if s.Escaped {
		return s.Text, nil
	}
	return html.EscapeString(s.Text), nil
}
