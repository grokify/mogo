package tokenizer

import (
	"golang.org/x/net/html"
)

type DescriptionLists []DescriptionList

func (dls DescriptionLists) Strings() [][][]string {
	strs := [][][]string{}
	for _, dl := range dls {
		strs = append(strs, dl.Strings())
	}
	return strs
}

type DescriptionList []Description

func (dl DescriptionList) Strings() [][]string {
	strs := [][]string{}
	for _, d := range dl {
		strs = append(strs, d.Strings())
	}
	return strs
}

type Description struct {
	Term        []html.Token
	Description []html.Token
}

func (d *Description) Empty() bool {
	return len(d.Term) == 0 && len(d.Description) == 0
}

func (d *Description) TermString() string {
	return Tokens(d.Term).String()
}

func (d *Description) DescriptionString() string {
	return Tokens(d.Description).String()
}

func (d *Description) Strings() []string {
	return []string{d.TermString(), d.DescriptionString()}
}