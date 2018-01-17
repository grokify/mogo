package nethttputil

import (
	"net/http"
	"strings"

	"github.com/grokify/gotilla/strings/stringsutil"
)

func GetReqHeader(req *http.Request, headerName string) string {
	return strings.TrimSpace(req.Header.Get(headerName))
}

func GetSplitReqHeader(req *http.Request, headerName, sep string) []string {
	return stringsutil.SliceTrimSpace(strings.Split(
		string(req.Header.Get(headerName)), sep))
}
