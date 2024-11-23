package markdown

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/grokify/mogo/type/slicesutil"
	"github.com/grokify/mogo/type/stringsutil"
)

// BoldText bodifies the identified text. It looks for start of words
// using a word boundary and will arbirarily end to match words with
// different suffixes.
func BoldText(haystack, needle string) string {
	output := haystack
	if len(needle) == 0 {
		return "**" + haystack + "**"
	}
	return regexp.MustCompile(`(?i)\b(`+regexp.QuoteMeta(needle)+`)`).ReplaceAllString(output, "**$1**")
}

func URLToMarkdownLinkHostname(url string) string {
	rx := regexp.MustCompile(`(?i)^https?://([^/]+)(/[^/])`)
	m := rx.FindStringSubmatch(url)
	if len(m) > 1 {
		suffix := ""
		if len(m) > 2 {
			suffix = "..."
		}
		return fmt.Sprintf("[%v%v](%v)", m[1], suffix, url)
	}
	return url
}

type Link struct {
	Text string
	URL  string
}

func (lnk Link) Markdown() string {
	return Linkify(lnk.URL, lnk.Text)
}

type Links []Link

func (lnks Links) Texts(condense, dedupe, sortAsc bool) []string {
	var out []string
	for _, lnk := range lnks {
		out = append(out, lnk.Text)
	}
	if condense {
		out = stringsutil.SliceCondenseSpace(out, dedupe, sortAsc)
	} else {
		if dedupe {
			slicesutil.Dedupe(out)
		}
		if sortAsc {
			sort.Strings(out)
		}
	}
	return out
}

// Linkify constructs a link from url and text inputs.
// This function does not handle escaping so it should
// be done before hand, e.g. using `\[` instead of `[`
// by itself in the `text`.
func Linkify(url, text string) string {
	url = strings.TrimSpace(url)
	if len(url) == 0 {
		return text
	}
	if len(text) == 0 {
		return "[" + url + "](" + url + ")"
	}
	return "[" + text + "](" + url + ")"
}

var rxParseLink = regexp.MustCompile(`\[(.*)\]\((.*)\)`)

// ParseLink returns text and url.
func ParseLink(s string) (string, string) {
	m := rxParseLink.FindStringSubmatch(s)
	if len(m) > 0 {
		return m[1], m[2]
	}
	return s, ""
}
