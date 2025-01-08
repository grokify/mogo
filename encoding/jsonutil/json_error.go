package jsonutil

import (
	"encoding/json"
	"strings"

	"github.com/grokify/mogo/errors/errorsutil"
)

func UnmarshalWithLoc(data []byte, v any) error {
	err := json.Unmarshal(data, v)
	if err == nil {
		return nil
	}
	if syntaxErr, ok := err.(*json.SyntaxError); ok {
		line, col := getErrorLine(data, syntaxErr.Offset)
		if line >= 0 && col > 0 {
			return errorsutil.Wrapf(err, "json syntax error at line (%d), column (%d)", line, col)
		}
	}
	return err
}

// getErrorLine calculates the line and column number from a byte offset
func getErrorLine(data []byte, offset int64) (line int, col int) {
	lines := strings.Split(string(data), "\n")
	offsetErr := int(offset)
	line = -1
	col = -1
	cur := 0

	for i, l := range lines {
		line = i + 1
		offsetLineEnd := cur + len(l) + 1
		if offsetErr <= offsetLineEnd {
			col = offsetErr - cur
			break
		} else {
			cur += len(l) + 1
		}
	}

	return
}
