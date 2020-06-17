package stringsutil

type Quoter struct {
	Beg string
	End string
}

func (qtr Quoter) Quote(input string) string {
	return qtr.Beg + input + qtr.End
}

func Quote(str, beg, end string) string {
	return beg + str + end
}
