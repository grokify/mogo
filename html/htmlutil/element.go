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
	AttributeHref    = "href"
	AttributeOnclick = "onclick"
	AttributeStyle   = "style"
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
	InnerHTML []stringsutil.StringableWithErr
}

func NewElement() *Element {
	return &Element{
		Attrs:     map[string][]string{},
		InnerHTML: []stringsutil.StringableWithErr{}}
}

func (el *Element) AddAttribute(key string, values ...string) error {
	key = strings.TrimSpace(key)
	if len(key) == 0 {
		return ErrAttributeNameIsRequired
	}
	if el.Attrs == nil {
		el.Attrs = map[string][]string{}
	}
	if _, ok := el.Attrs[key]; !ok {
		el.Attrs[key] = []string{}
	}
	if len(values) > 0 {
		el.Attrs[key] = append(el.Attrs[key], values...)
	}
	return nil
}

func (el *Element) AddInnerHTML(innerHTML stringsutil.StringableWithErr) {
	el.InnerHTML = append(el.InnerHTML, innerHTML)
}

func (el *Element) AddInnerHTMLText(text string, escaped bool) {
	el.InnerHTML = append(el.InnerHTML, &Text{
		Text:    text,
		Escaped: escaped})
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

func (el *Element) String() (string, error) {
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
		} else if key == AttributeOnclick {
			attrs = append(attrs, BuildAttributeHTML(key, vals, DelimitSemicolon, false))
		} else {
			attrs = append(attrs, BuildAttributeHTML(key, vals, DelimitSemicolon, true))
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
		childStr, err := child.String()
		if err != nil {
			return "", err
		}
		innerHTML += childStr
	}

	closingTag := ""
	if openingTagClose == ">" {
		closingTag = "</" + el.TagName + ">"
	}
	return openingTag + innerHTML + closingTag, nil
}
