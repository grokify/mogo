package htmlutil

import (
	"errors"
	"fmt"
	"html"
	"strings"

	"github.com/grokify/mogo/type/maputil"
	"github.com/grokify/mogo/type/stringsutil"
)

const (
	TagDiv           = "div"
	AttributeClass   = "class"
	AttributeOnclick = "onclick"
	DelimitSemicolon = ";"
	DelimitSpace     = " "
)

var (
	ErrAttributeNameIsRequired = errors.New("attribute name is required")
	ErrTagNameIsRequired       = errors.New("tag name is required")
)

type Element struct {
	TagName   string
	Attrs     map[string][]string
	SelfClose bool
	InnerHTML []stringsutil.Stringable
}

func NewElement() Element {
	return Element{
		Attrs:     map[string][]string{},
		InnerHTML: []stringsutil.Stringable{}}
}

func (el Element) AddAttribute(key string, values ...string) error {
	key = strings.TrimSpace(key)
	if len(key) == 0 {
		return ErrAttributeNameIsRequired
	}
	if _, ok := el.Attrs[key]; !ok {
		el.Attrs[key] = []string{}
	}
	if len(values) > 0 {
		el.Attrs[key] = append(el.Attrs[key], values...)
	}
	return nil
}

func BuildAttributeHTML(key string, values []string, delimiter string, htmlEscape bool) string {
	vals2 := []string{}
	if htmlEscape {
		for _, val := range values {
			vals2 = append(vals2, html.EscapeString(val))
		}
	} else {
		vals2 = values
	}
	valStr := strings.Join(vals2, delimiter)
	return fmt.Sprintf(`%s="%s"`, key, valStr)
}

func (el Element) String() (string, error) {
	el.TagName = strings.TrimSpace(el.TagName)
	if len(el.TagName) == 0 {
		return "", ErrTagNameIsRequired
	}
	attrs := []string{}
	keysSorted := maputil.StringKeys(el.Attrs, nil, true)
	for _, key := range keysSorted {
		vals, ok := el.Attrs[key]
		if !ok {
			panic("key not found")
		}
		vals = stringsutil.SliceTrimSpace(vals, true)
		if len(vals) == 0 {
			attrs = append(attrs, key)
		} else if key == AttributeClass {
			attrs = append(attrs, BuildAttributeHTML(
				key,
				stringsutil.SliceCondenseSpace(vals, true, false),
				DelimitSpace, true))
			//escaped := []string{}
			//for _, val := range vals {
			//	escaped = append(escaped, html.EscapeString(val))
			//}
			//attrs = append(attrs, key+"=\""+strings.Join(escaped, DelimitSpace)+"\"")
		} else if key == AttributeOnclick {
			attrs = append(attrs, BuildAttributeHTML(key, vals, DelimitSemicolon, false))
			//attrs = append(attrs, key+"=\""+strings.Join(vals, DelimitSemicolon)+"\"")
		} else {
			attrs = append(attrs, BuildAttributeHTML(key, vals, DelimitSemicolon, true))
			//escaped := []string{}
			//for _, val := range vals {
			//	escaped = append(escaped, html.EscapeString(val))
			//}
			//attrs = append(attrs, key+"=\""+strings.Join(escaped, DelimitSemicolon)+"\"")
		}
	}
	attrsStr := strings.Join(attrs, " ")
	if len(attrsStr) > 0 {
		attrsStr = " " + attrsStr
	}
	openingTagClose := ">"
	if el.SelfClose && len(el.InnerHTML) == 0 {
		openingTagClose = " />"
	}
	openingTag := "<" + el.TagName + attrsStr + openingTagClose

	var innerHTML string
	for _, child := range el.InnerHTML {
		innerHTML += child.String()
	}

	closingTag := ""
	if openingTagClose == ">" {
		closingTag = "</" + el.TagName + ">"
	}
	return openingTag + innerHTML + closingTag, nil

	/*
		elString := "<" + el.TagName
		if len(attrs) > 0 {
			elString += " " + strings.Join(attrs, " ")
		}
		if len(el.InnerHTML) == 0 {
			if el.SelfClose {
				elString += " />"
			} else {
				elString += "></" + el.TagName + ">"
			}
			return elString
		}

		for _, child := range el.InnerHTML {
			elString += child.String()
		}
		elString += "</" + el.TagName + ">"

		return elString
	*/
}
