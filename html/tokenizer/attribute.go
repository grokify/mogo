package tokenizer

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

const (
	AttrHref = "href"
)

type Attributes []html.Attribute

func (attrs Attributes) GetOne(attributeKey string) (html.Attribute, error) {
	matches := []html.Attribute{}
	for _, attr := range attrs {
		if attributeKey == strings.TrimSpace(attr.Key) {
			matches = append(matches, attr)
		}
	}
	if len(matches) == 0 {
		return html.Attribute{}, fmt.Errorf("attribute key not found [%s]", attributeKey)
	} else if len(matches) > 1 {
		return html.Attribute{}, fmt.Errorf("attribute key found multiple times [%d]", len(matches))
	}
	return matches[0], nil
}

func TokenAttribute(token html.Token, attrName string, attrLower bool) (string, error) {
	attrNameMatch := attrName
	if attrLower {
		attrNameMatch = strings.ToLower(attrNameMatch)
	}
	for _, attr := range token.Attr {
		if attrLower && strings.ToLower(attr.Key) == attrNameMatch {
			return attr.Val, nil
		} else if !attrLower && attr.Key == attrNameMatch {
			return attr.Val, nil
		}
	}
	return "", fmt.Errorf("attribute not found [%s] matchLower [%v]", attrName, attrLower)
}
