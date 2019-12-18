package markdown

import "strings"

type PresentationData struct {
	Slides []RemarkSlideData
}

type RemarkSlideData struct {
	Layout   string
	Class    string
	Markdown string
}

func (data *RemarkSlideData) ToRemarkString() string {
	parts := []string{}
	data.Layout = strings.TrimSpace(data.Layout)
	data.Class = strings.TrimSpace(data.Class)
	data.Markdown = strings.TrimSpace(data.Markdown)
	if len(data.Layout) > 0 {
		parts = append(parts, "layout: "+data.Layout)
	}
	if len(data.Class) > 0 {
		parts = append(parts, "class: "+data.Class)
	}
	if len(data.Markdown) > 0 {
		parts = append(parts, data.Markdown)
	}
	return strings.Join(parts, "\n\n")
}
