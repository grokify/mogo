package fasthttputil

import (
	"strings"

	"github.com/grokify/gotilla/strings/stringsutil"
	"github.com/valyala/fasthttp"
)

func GetReqHeader(ctx *fasthttp.RequestCtx, headerName string) string {
	return strings.TrimSpace(string(ctx.QueryArgs().Peek(headerName)))
}

func GetSplitReqHeader(ctx *fasthttp.RequestCtx, headerName, sep string) []string {
	return stringsutil.SliceTrimSpace(strings.Split(
		string(ctx.QueryArgs().Peek(headerName)), sep))
}
