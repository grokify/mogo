package tokenizer

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TokenizerDescriptionLists(z *html.Tokenizer) (DescriptionLists, error) {
	dls := DescriptionLists{}
	for {
		dl, err := TokenizerDescriptionListNext(z)
		if err != nil {
			return dls, err
		}
		if len(dl) == 0 {
			break
		}
		dls = append(dls, dl)
	}
	return dls, nil
}

func TokenizerDescriptionListNext(z *html.Tokenizer) (DescriptionList, error) {
	descriptionList := DescriptionList{}
	skipErrs := false
	dlToks, err := TokensBetweenAtom(z, skipErrs, true, atom.Dl)
	if err != nil {
		return descriptionList, err
	}
	descriptionList = ParseDescriptionListTokens(dlToks...)
	return descriptionList, nil
}

const (
	matchingTerm = "term"
	matchingDesc = "desc"
)

func ParseDescriptionListTokens(toks ...html.Token) DescriptionList {
	dl := DescriptionList{}
	var curDesc Description
	startDt := NewTokenFilter(html.StartTagToken, atom.Dt)
	endDt := NewTokenFilter(html.EndTagToken, atom.Dt)
	startDd := NewTokenFilter(html.StartTagToken, atom.Dd)
	endDd := NewTokenFilter(html.EndTagToken, atom.Dd)
	matching := ""
	for _, tok := range toks {
		if startDt.Match(tok) {
			curDesc.Term = append(curDesc.Term, tok)
			matching = matchingTerm
		} else if endDt.Match(tok) {
			curDesc.Term = append(curDesc.Term, tok)
			matching = ""
		} else if matching == matchingTerm {
			curDesc.Term = append(curDesc.Term, tok)
		} else if startDd.Match(tok) {
			curDesc.Description = append(curDesc.Description, tok)
			matching = matchingDesc
		} else if endDd.Match(tok) {
			curDesc.Description = append(curDesc.Description, tok)
			dl = append(dl, curDesc)
			matching = ""
			curDesc = Description{}
		} else if matching == matchingDesc {
			curDesc.Description = append(curDesc.Description, tok)
		}
	}
	return dl
}
