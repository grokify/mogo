package htmlutil

import (
	"html"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	xhtml "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// ChartColor1 is the color palette for Google Charts as collected by
// Craig Davis here: https://gist.github.com/there4/2579834
var ChartColor1 = [...]string{
	"#3366CC",
	"#DC3912",
	"#FF9900",
	"#109618",
	"#990099",
	"#3B3EAC",
	"#0099C6",
	"#DD4477",
	"#66AA00",
	"#B82E2E",
	"#316395",
	"#994499",
	"#22AA99",
	"#AAAA11",
	"#6633CC",
	"#E67300",
	"#8B0707",
	"#329262",
	"#5574A6",
	"#3B3EAC",
}

// Link is a struct to hold information for an HTML link.
type Link struct {
	Href      string
	InnerHTML string
}

const (
	Color2GreenHex       = "#00FF2A"
	Color2YellowHex      = "#DFDD13"
	Color2RedHex         = "#FF0000"
	RingCentralOrangeHex = "#FF8800"
	RingCentralBlueHex   = "#0073AE"
	RingCentralGreyHex   = "#585858"
)

var (
	bluemondayStrictPolicy                  = bluemonday.StrictPolicy()
	rxHTMLToTextNewLines     *regexp.Regexp = regexp.MustCompile(`(?i:</?p>)`)
	rxCarriageReturn         *regexp.Regexp = regexp.MustCompile(`\r`)
	rxDiv                    *regexp.Regexp = regexp.MustCompile(`(?i)<div>`)
	rxLineFeed               *regexp.Regexp = regexp.MustCompile(`\n`)
	rxLineFeedMore           *regexp.Regexp = regexp.MustCompile(`\n+`)
	rxCarriageReturnLineFeed *regexp.Regexp = regexp.MustCompile(`\r\n`)
	rxLineFeedMore2          *regexp.Regexp = regexp.MustCompile(`\n\n+`)
	doubleLinefeed                          = "\n\n"
	// rxCarriageReturnLineFeedMore *regexp.Regexp = regexp.MustCompile(`[\r\n]+`)
	// rxEndingSpacesLineFeed       *regexp.Regexp = regexp.MustCompile(`\s+\n`)
)

func EscapeStrings(s []string) []string {
	var n []string
	for _, si := range s {
		n = append(n, html.EscapeString(si))
	}
	return n
}

func StreamlineCRLFs(s string) string {
	newLines := []string{}
	extLines := strings.Split(s, "\n")
	for _, line := range extLines {
		newLines = append(newLines, strings.TrimSpace(line))
	}
	s2 := strings.Join(newLines, "\n")
	s2 = rxLineFeedMore2.ReplaceAllString(
		rxCarriageReturn.ReplaceAllString(
			rxCarriageReturnLineFeed.ReplaceAllString(s2, "\n"),
			"\n"),
		"\n",
	)
	return s2
}

// HTMLToTextCondensed removes HTML tags, unescapes HTML entities,
// and removes extra whitespace including non-breaking spaces.
func HTMLToTextCondensed(s string) string {
	return strings.Join(
		strings.Fields(
			html.UnescapeString(
				bluemondayStrictPolicy.Sanitize(s),
			),
		),
		" ",
	)
}

// HTMLToText converts HTML to multi-line text.
func HTMLToText(s string) string {
	return rxLineFeedMore2.ReplaceAllString(
		strings.TrimSpace(
			html.UnescapeString(
				bluemondayStrictPolicy.Sanitize(
					rxHTMLToTextNewLines.ReplaceAllString(
						rxDiv.ReplaceAllString(s, doubleLinefeed),
						"$1"+doubleLinefeed),
				),
			),
		),
		doubleLinefeed,
	)
}

func HTMLToTextH1(b []byte, policy *bluemonday.Policy) (string, error) {
	return HTMLToTextAtom(b, policy, atom.H1)
}

func HTMLToTextAtom(b []byte, policy *bluemonday.Policy, a atom.Atom) (string, error) {
	z := NewTokenizerBytes(b)
	// filter := []golanghtml.Token{{DataAtom: a}}
	opts := NextTokensOpts{
		SkipErrors:     false,
		IncludeChain:   true,
		InclusiveMatch: true,
		StartFilter:    []xhtml.Token{{DataAtom: a, Type: xhtml.StartTagToken}},
		EndFilter:      []xhtml.Token{{DataAtom: a, Type: xhtml.EndTagToken}},
	}
	toks, err := NextTokens(z, opts)
	// toks, err := TokensBetweenAtom(t, false, true, a)
	if err != nil {
		return "", err
	}
	if policy == nil {
		policy = bluemonday.StrictPolicy()
	}
	return strings.TrimSpace(policy.Sanitize(toks.String())), nil
}

func SimplifyHTMLText(s string) string {
	text := HTMLToText(s)
	lines := strings.Split(text, "\n")
	newlines := []string{}
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if len(line) > 0 {
			newlines = append(newlines, "<p>"+line+"</p>")
		}
	}
	return strings.Join(newlines, "")
}

func TextToHTML(s string) string {
	return rxLineFeed.ReplaceAllString(StreamlineCRLFs(s), "<br/>")
}

func TextToHTMLBr2(s string) string {
	return rxLineFeed.ReplaceAllString(
		rxLineFeedMore.ReplaceAllString(StreamlineCRLFs(s), "\n"),
		"<br/><br/>",
	)
}
