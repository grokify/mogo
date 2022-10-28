package htmlutil

import (
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

func (el Element) Add(key string, values ...string) {
	key = strings.ToLower(strings.TrimSpace(key))
	if len(key) == 0 {
		return
	}
	if len(values) == 0 {
		if _, ok := el.Attrs[key]; !ok {
			el.Attrs[key] = []string{}
		}
		return
	}
	el.Attrs[key] = append(el.Attrs[key], values...)
}

func (el Element) String() string {
	el.TagName = strings.ToLower(strings.TrimSpace(el.TagName))
	if len(el.TagName) == 0 {
		el.TagName = TagDiv
	}
	attrs := []string{}
	keysSorted := maputil.StringKeys(el.Attrs, nil, true)
	for _, key := range keysSorted {
		vals, ok := el.Attrs[key]
		if !ok {
			panic("key not found")
		}
		//}
		//for key, vals := range el.Attrs {
		if len(vals) == 0 {
			attrs = append(attrs, key)
		} else if key == AttributeClass {
			escaped := []string{}
			for _, val := range vals {
				escaped = append(escaped, html.EscapeString(val))
			}
			attrs = append(attrs, key+"=\""+strings.Join(escaped, DelimitSpace)+"\"")
		} else if key == AttributeOnclick {
			attrs = append(attrs, key+"=\""+strings.Join(vals, DelimitSemicolon)+"\"")
		} else {
			escaped := []string{}
			for _, val := range vals {
				escaped = append(escaped, html.EscapeString(val))
			}
			attrs = append(attrs, key+"=\""+strings.Join(escaped, DelimitSemicolon)+"\"")
		}
	}
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
}
