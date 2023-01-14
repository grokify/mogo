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

func (attrs Attributes) GetVal(key string) []string {
	vals := []string{}
	for _, a := range attrs {
		if key == a.Key {
			vals = append(vals, a.Val)
		}
	}
	return vals
}

func (attrs Attributes) Has(key, val string) bool {
	for _, a := range attrs {
		if a.Key == key && a.Val == val {
			return true
		}
	}
	return false
}

func TokenAttribute(token html.Token, attrName string) (string, error) {
	attrNameWant := strings.TrimSpace(strings.ToLower(attrName))
	for _, attr := range token.Attr {
		if strings.TrimSpace(strings.ToLower(attr.Key)) == attrNameWant {
			return attr.Val, nil
		}
	}
	return "", fmt.Errorf("attribute not found [%s]", attrName)
}
