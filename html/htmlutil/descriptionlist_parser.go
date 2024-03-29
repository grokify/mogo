package htmlutil

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
	opts := NextTokensOpts{
		SkipErrors:     false,
		IncludeChain:   true,
		InclusiveMatch: true,
		StartFilter:    []html.Token{{DataAtom: atom.Dl, Type: html.StartTagToken}},
		EndFilter:      []html.Token{{DataAtom: atom.Dl, Type: html.EndTagToken}},
	}
	dlToks, err := NextTokens(z, opts)
	// dlToks, err := TokensBetweenAtom(z, skipErrs, true, atom.Dl)
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
	staDt := Tokens([]html.Token{{Type: html.StartTagToken, DataAtom: atom.Dt}})
	endDt := Tokens([]html.Token{{Type: html.EndTagToken, DataAtom: atom.Dt}})
	staDd := Tokens([]html.Token{{Type: html.StartTagToken, DataAtom: atom.Dd}})
	endDd := Tokens([]html.Token{{Type: html.EndTagToken, DataAtom: atom.Dd}})
	// staDt := NewTokenFilter(html.StartTagToken, atom.Dt)
	// endDt := NewTokenFilter(html.EndTagToken, atom.Dt)
	// staDd := NewTokenFilter(html.StartTagToken, atom.Dd)
	// endDd := NewTokenFilter(html.EndTagToken, atom.Dd)
	matching := ""
	for _, tok := range toks {
		if staDt.MatchLeft(tok, nil) {
			curDesc.Term = append(curDesc.Term, tok)
			matching = matchingTerm
		} else if endDt.MatchLeft(tok, nil) {
			curDesc.Term = append(curDesc.Term, tok)
			matching = ""
		} else if matching == matchingTerm {
			curDesc.Term = append(curDesc.Term, tok)
		} else if staDd.MatchLeft(tok, nil) {
			curDesc.Description = append(curDesc.Description, tok)
			matching = matchingDesc
		} else if endDd.MatchLeft(tok, nil) {
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
