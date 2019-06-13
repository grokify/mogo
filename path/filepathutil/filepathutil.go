package filepathutil

// $ go get github.com/grokify/gotilla/git/apps/gitremoteaddupstream

import (
	"os"
	"strconv"
	"strings"
)

func PathSeparatorString() (string, error) {
	return strconv.Unquote(strconv.QuoteRune(os.PathSeparator))
}

func CurLeafDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	sep, err := PathSeparatorString()
	if err != nil {
		return "", err
	}
	dirParts := strings.Split(dir, sep)
	if len(dirParts) > 0 {
		return dirParts[len(dirParts)-1], nil
	}
	return "", nil
}
