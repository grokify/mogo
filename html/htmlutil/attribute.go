package htmlutil

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

const (
	AttrHref = "href"
)

type Attributes []html.Attribute

func (attrs Attributes) Exists(attr html.Attribute) bool {
	for _, attrTry := range attrs {
		if attrTry.Namespace == attr.Namespace &&
			attrTry.Key == attr.Key &&
			attrTry.Val == attr.Val {
			return true
		}
	}
	return false
}

func (attrs Attributes) Find(ns, key, val *string, n int) []html.Attribute {
	finds := []html.Attribute{}
	if n == 0 {
		return finds
	}
	for _, attr := range attrs {
		if ns != nil && *ns != attr.Namespace {
			continue
		}
		if key != nil && *key != attr.Key {
			continue
		}
		if val != nil && *val != attr.Val {
			continue
		}
		finds = append(finds, attr)
		if n == len(finds) {
			return finds
		}
	}
	return finds
}

func (attrs Attributes) FindOne(ns, key, val *string, errOnNone, errOnMulti bool) (html.Attribute, error) {
	finds := attrs.Find(ns, key, val, 2)
	if errOnNone && len(finds) == 0 {
		return html.Attribute{}, fmt.Errorf("attribute key not found for N[%v]K[%v]V[%v]", *ns, *key, *val)
	} else if errOnMulti && len(finds) > 1 {
		return html.Attribute{}, fmt.Errorf("attribute key found multiple times [%d] for N[%v]K[%v]V[%v]", len(finds), *ns, *key, *val)
	}
	return finds[0], nil
}

func (attrs Attributes) FindVal(ns, key *string) string {
	finds := attrs.Find(ns, key, nil, 1)
	for _, a := range finds {
		return a.Val
	}
	return ""
}

func (attrs Attributes) FindVals(ns, key *string, n int) []string {
	finds := attrs.Find(ns, key, nil, n)
	vals := []string{}
	for _, a := range finds {
		vals = append(vals, a.Val)
	}
	return vals
}

func (attrs Attributes) FindExists(ns, key, val *string) bool {
	finds := attrs.Find(ns, key, val, 1)
	return len(finds) > 0
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
