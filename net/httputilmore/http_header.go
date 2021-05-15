package httputilmore

import (
	"fmt"
	"net/http"
	"strings"
)

// Constants ensuring that header names are correctly spelled and consistently cased.
const (
	HeaderAccept                       = "Accept"
	HeaderAccessControlAllowHeaders    = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowMethods    = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowOrigin     = "Access-Control-Allow-Origin"
	HeaderAccessControlRequestHeaders  = "Access-Control-Request-Headers"
	HeaderAccessControlRequestMethod   = "Access-Control-Request-Method"
	HeaderAuthorization                = "Authorization"
	HeaderCacheControl                 = "Cache-Control"
	HeaderContentDisposition           = "Content-Disposition"
	HeaderContentEncoding              = "Content-Encoding"
	HeaderContentLanguage              = "Content-Language"
	HeaderContentLength                = "Content-Length"
	HeaderContentMD5                   = "Content-MD5"
	HeaderContentTransferEncoding      = "Content-Transfer-Encoding"
	HeaderContentType                  = "Content-Type"
	HeaderDate                         = "Date"
	HeaderIfMatch                      = "If-Match"
	HeaderIfModifiedSince              = "If-Modified-Since"
	HeaderIfNoneMatch                  = "If-None-Match"
	HeaderIfUnmodifiedSince            = "If-Unmodified-Since"
	HeaderLocation                     = "Location"
	HeaderRange                        = "Range"
	HeaderUserAgent                    = "User-Agent"
	HeaderXContentTypeOptions          = "X-Content-Type-Options"
	ContentTypeAppJson                 = "application/json"
	ContentTypeAppJsonUtf8             = "application/json; charset=utf-8"
	ContentTypeAppFormUrlEncoded       = "application/x-www-form-urlencoded"
	ContentTypeAppXml                  = "application/xml"
	ContentTypeAppXmlUtf8              = "application/xml; charset=utf-8"
	ContentTypeTextCalendarUtf8Request = `text/calendar; charset="utf-8"; method=REQUEST`
	ContentTypeTextHtml                = "text/html"
	ContentTypeTextHtmlUtf8            = "text/html; charset=utf-8"
	ContentTypeTextMarkdown            = "text/markdown"
	ContentTypeTextPlain               = "text/plain"
	ContentTypeTextPlainUsAscii        = "text/plain; charset=us-ascii"
	ContentTypeTextPlainUtf8           = "text/plain; charset=utf-8"
	ContentTypeTextXmlUtf8             = "text/xml; charset=utf-8"
	SchemeHTTPS                        = "https"
)

type HTTPMethod string

const (
	MethodConnect HTTPMethod = http.MethodConnect
	MethodDelete             = http.MethodDelete
	MethodGet                = http.MethodGet
	MethodHead               = http.MethodHead
	MethodOptions            = http.MethodOptions
	MethodPatch              = http.MethodPatch
	MethodPost               = http.MethodPost
	MethodPut                = http.MethodPut
	MethodTrace              = http.MethodTrace
)

// ParseHTTPMethod returns a HTTPMethod type for a string.
func ParseHTTPMethod(method string) (HTTPMethod, error) {
	method = strings.ToUpper(strings.TrimSpace(method))
	switch method {
	case http.MethodConnect:
		return MethodConnect, nil
	case http.MethodDelete:
		return MethodDelete, nil
	case http.MethodGet:
		return MethodGet, nil
	case http.MethodHead:
		return MethodHead, nil
	case http.MethodOptions:
		return MethodOptions, nil
	case http.MethodPatch:
		return MethodPatch, nil
	case http.MethodPost:
		return MethodPost, nil
	case http.MethodPut:
		return MethodPut, nil
	case http.MethodTrace:
		return MethodTrace, nil
	}
	return MethodConnect, fmt.Errorf("E_NO_METHOD_FOR [%v]", method)
}

/*
// HTTPMethod is a type of HTTP Methods. Of note, do not
// rely on the integer value but use for definitions.
type HTTPMethod int

const (
	MethodConnect HTTPMethod = iota
	MethodDelete
	MethodGet
	MethodHead
	MethodOptions
	MethodPatch
	MethodPost
	MethodPut
	MethodTrace
)

var methods = [...]string{
	http.MethodConnect,
	http.MethodDelete,
	http.MethodGet,
	http.MethodHead,
	http.MethodOptions,
	http.MethodPatch,
	http.MethodPost,
	http.MethodPut,
	http.MethodTrace,
}

// String returns the English name of the HTTPMethod.
func (method HTTPMethod) String() string {
	if MethodConnect <= method && method <= MethodTrace {
		return methods[method]
	}
	buf := make([]byte, 20)
	n := fmtInt(buf, uint64(method))
	return "%!HTTPMethod(" + string(buf[n:]) + ")"
}

// ParseHTTPMethod returns a HTTPMethod type for a string.
func ParseHTTPMethod(method string) (HTTPMethod, error) {
	method = strings.ToUpper(strings.TrimSpace(method))
	for i, try := range methods {
		if method == try {
			return HTTPMethod(i), nil
		}
	}
	return MethodConnect, fmt.Errorf("E_NO_METHOD_FOR [%v]", method)
}
*/

// fmtInt formats v into the tail of buf.
// It returns the index where the output begins.
func fmtInt(buf []byte, v uint64) int {
	w := len(buf)
	if v == 0 {
		w--
		buf[w] = '0'
	} else {
		for v > 0 {
			w--
			buf[w] = byte(v%10) + '0'
			v /= 10
		}
	}
	return w
}

// NewHeadersMSS returns a `http.Header` struct give a `map[string]string`
func NewHeadersMSS(headersMap map[string]string) http.Header {
	header := http.Header{}
	for k, v := range headersMap {
		header.Add(k, v)
	}
	return header
}
