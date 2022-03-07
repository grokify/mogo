package grep

import (
	"strings"

	"github.com/grokify/mogo/type/stringsutil"
)

type GrepResults []GrepResult

func (gr GrepResults) Files() []string {
	files := []string{}
	for _, g := range gr {
		g.File = strings.TrimSpace(g.File)
		if len(g.File) > 0 {
			files = append(files, g.File)
		}
	}
	return stringsutil.SliceCondenseSpace(files, true, true)
}

type GrepResult struct {
	File string
	Line string
}

// ParseGrep will parse the results from `os/executil.ExecSimple`.
func ParseGrep(data []byte) GrepResults {
	res := GrepResults{}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		lineParts := strings.SplitN(line, ":", 2)
		if len(lineParts) == 2 {
			res = append(res, GrepResult{
				File: lineParts[0],
				Line: lineParts[1]})
		}
	}
	return res
}
