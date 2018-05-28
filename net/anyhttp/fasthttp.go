package anyhttp

import (
	"io"
	"mime/multipart"
	"net"
	"strings"

	"github.com/valyala/fasthttp"
)

type RequestFastHttp struct {
	Raw       *fasthttp.RequestCtx
	allArgs   *ArgsFastHttpMulti
	queryArgs *ArgsFastHttp
	postArgs  *ArgsFastHttp
}

func NewRequestFastHttp(ctx *fasthttp.RequestCtx) *RequestFastHttp {
	return &RequestFastHttp{
		Raw:       ctx,
		allArgs:   &ArgsFastHttpMulti{Raw: []*fasthttp.Args{ctx.QueryArgs(), ctx.PostArgs()}},
		queryArgs: &ArgsFastHttp{Raw: ctx.QueryArgs()},
		postArgs:  &ArgsFastHttp{Raw: ctx.PostArgs()},
	}
}

func (r RequestFastHttp) ParseForm() error                        { return nil }
func (r RequestFastHttp) AllArgs() Args                           { return r.allArgs }
func (r RequestFastHttp) QueryArgs() Args                         { return r.queryArgs }
func (r RequestFastHttp) PostArgs() Args                          { return r.postArgs }
func (r RequestFastHttp) Method() []byte                          { return r.Raw.Method() }
func (r RequestFastHttp) MultipartForm() (*multipart.Form, error) { return r.Raw.MultipartForm() }
func (r RequestFastHttp) RemoteAddr() net.Addr                    { return r.Raw.RemoteAddr() }
func (r RequestFastHttp) RemoteAddress() string                   { return r.Raw.RemoteAddr().String() }
func (r RequestFastHttp) RequestURI() []byte                      { return r.Raw.RequestURI() }
func (r RequestFastHttp) UserAgent() []byte                       { return r.Raw.UserAgent() }

type ResponseFastHttp struct {
	Raw *fasthttp.RequestCtx
}

func NewResponseFastHttp(ctx *fasthttp.RequestCtx) ResponseFastHttp {
	return ResponseFastHttp{Raw: ctx}
}

func (w ResponseFastHttp) GetHeader(k string) []byte { return w.Raw.Response.Header.Peek(k) }
func (w ResponseFastHttp) SetHeader(k, v string)     { w.Raw.Response.Header.Set(k, v) }
func (w ResponseFastHttp) SetStatusCode(code int)    { w.Raw.SetStatusCode(code) }
func (w ResponseFastHttp) SetContentType(ct string)  { w.Raw.SetContentType(ct) }
func (w ResponseFastHttp) SetBodyBytes(body []byte) (int, error) {
	w.Raw.SetBody(body)
	return -1, nil
}
func (w ResponseFastHttp) SetBodyStream(bodyStream io.Reader, bodySize int) error {
	w.Raw.SetBodyStream(bodyStream, bodySize)
	return nil
}
func (w ResponseFastHttp) SetCookie(cookie *Cookie) {
	w.Raw.Response.Header.SetCookie(cookie.ToFastHttp())
}

type ArgsFastHttp struct{ Raw *fasthttp.Args }

func NewArgsFastHttp(args *fasthttp.Args) ArgsFastHttp {
	return ArgsFastHttp{Raw: args}
}

func (a ArgsFastHttp) GetBytes(key string) []byte        { return a.Raw.Peek(key) }
func (a ArgsFastHttp) GetBytesSlice(key string) [][]byte { return a.Raw.PeekMulti(key) }
func (a ArgsFastHttp) GetString(key string) string       { return string(a.Raw.Peek(key)) }
func (a ArgsFastHttp) GetStringSlice(key string) []string {
	slice := a.Raw.PeekMulti(key)
	newSlice := []string{}
	for _, bytes := range slice {
		newSlice = append(newSlice, string(bytes))
	}
	return newSlice
}

type ArgsFastHttpMulti struct {
	Raw []*fasthttp.Args
}

func NewArgsFastHttpMulti(args []*fasthttp.Args) ArgsFastHttpMulti {
	return ArgsFastHttpMulti{Raw: args}
}

func (args ArgsFastHttpMulti) GetBytes(key string) []byte {
	for _, raw := range args.Raw {
		try := raw.Peek(key)
		if len(try) == 0 {
			return try
		}
	}
	return []byte("")
}

func (args ArgsFastHttpMulti) GetBytesSlice(key string) [][]byte {
	newSlice := [][]byte{}
	for _, raw := range args.Raw {
		slice := raw.PeekMulti(key)
		for _, bytes := range slice {
			if len(string(bytes)) > 0 {
				newSlice = append(newSlice, bytes)
			}
		}
	}
	return newSlice
}

func (args ArgsFastHttpMulti) GetString(key string) string {
	for _, raw := range args.Raw {
		try := strings.TrimSpace(string(raw.Peek(key)))
		if len(try) > 0 {
			return try
		}
	}
	return ""
}

func (args ArgsFastHttpMulti) GetStringSlice(key string) []string {
	newSlice := []string{}
	for _, raw := range args.Raw {
		slice := raw.PeekMulti(key)
		for _, bytes := range slice {
			try := strings.TrimSpace(string(bytes))
			if len(try) > 0 {
				newSlice = append(newSlice, try)
			}
		}
	}
	return newSlice
}

func NewResReqFastHttp(ctx *fasthttp.RequestCtx) (ResponseFastHttp, *RequestFastHttp) {
	return NewResponseFastHttp(ctx), NewRequestFastHttp(ctx)
}
