package htmlutil

// import "github.com/grokify/mogo/html/htmlutil"

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
	Term        Tokens
	Description Tokens
}

func (d *Description) Empty() bool {
	return len(d.Term) == 0 && len(d.Description) == 0
}

func (d *Description) TermString() string {
	return d.Term.String()
}

func (d *Description) DescriptionString() string {
	return d.Description.String()
}

func (d *Description) Strings() []string {
	return []string{d.TermString(), d.DescriptionString()}
}
