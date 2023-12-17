package stringsutil

import "strings"

type Quoter struct {
	Beg         string
	End         string
	SkipNesting bool
}

func (qtr Quoter) Quote(input string) string {
	if qtr.Beg == "" && qtr.End == "" {
		return input
	} else if qtr.SkipNesting && strings.Index(input, qtr.Beg) == 0 && ReverseIndex(input, qtr.End) == 0 {
		return input
	} else {
		return qtr.Beg + input + qtr.End
	}
}

func Quote(str, beg, end string) string {
	return beg + str + end
}
