package httputilmore

import (
	"net/http"
	"strings"

	"github.com/grokify/mogo/type/stringsutil"
)

// https://stackoverflow.com/questions/15407719/in-gos-http-package-how-do-i-get-the-query-string-on-a-post-request

func GetReqQueryParam(req *http.Request, paramName string) string {
	return strings.TrimSpace(req.URL.Query().Get(paramName))
}

func GetReqQueryParamSplit(req *http.Request, paramName, sep string) []string {
	return stringsutil.SliceTrimSpace(
		strings.Split(
			GetReqQueryParam(req, paramName),
			sep),
		true,
	)
}

/*
type RequestUtil struct {
	Request *http.Request
}

func (ru *RequestUtil) QueryParamString(paramName string) string {
	return strings.TrimSpace(ru.Request.URL.Query().Get(paramName))
}
*/
