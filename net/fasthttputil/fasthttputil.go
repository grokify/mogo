package fasthttputil

import (
	"strings"

	"github.com/grokify/simplego/type/stringsutil"
	"github.com/valyala/fasthttp"
)

func GetReqQueryParam(ctx *fasthttp.RequestCtx, headerName string) string {
	return strings.TrimSpace(string(ctx.QueryArgs().Peek(headerName)))
}

func GetSplitReqQueryParam(ctx *fasthttp.RequestCtx, headerName, sep string) []string {
	return stringsutil.SliceTrimSpace(
		strings.Split(
			string(ctx.QueryArgs().Peek(headerName)), sep),
		true,
	)
}
