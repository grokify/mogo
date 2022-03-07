package httputilmore

import (
	"bytes"
	"net/http"
)

// rw.Header().Set("WWW-Authenticate", "Basic realm=Restricted")

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
	HeaderWWWAuthenticate              = "WWW-Authenticate"
	HeaderXContentTypeOptions          = "X-Content-Type-Options"
	ContentTypeAppJSON                 = "application/json"
	ContentTypeAppJSONUtf8             = "application/json; charset=utf-8"
	ContentTypeAppOctetStream          = "application/octet-stream"
	ContentTypeAppFormURLEncoded       = "application/x-www-form-urlencoded"
	ContentTypeAppFormURLEncodedUtf8   = "application/x-www-form-urlencoded; charset=utf-8"
	ContentTypeAppXML                  = "application/xml"
	ContentTypeAppXMLUtf8              = "application/xml; charset=utf-8"
	ContentTypeTextCalendarUtf8Request = "text/calendar; charset=utf-8; method=REQUEST"
	ContentTypeTextHTML                = "text/html"
	ContentTypeTextHTMLUtf8            = "text/html; charset=utf-8"
	ContentTypeTextMarkdown            = "text/markdown"
	ContentTypeTextPlain               = "text/plain"
	ContentTypeTextPlainUsASCII        = "text/plain; charset=us-ascii"
	ContentTypeTextPlainUtf8           = "text/plain; charset=utf-8"
	ContentTypeTextXMLUtf8             = "text/xml; charset=utf-8"
	SchemeHTTPS                        = "https"
	WWWAuthenticateBasicRestricted     = "Basic realm=Restricted"
)

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

// HeaderMerge combines data from multiple `http.Header` structs.
func HeaderMerge(headers ...http.Header) http.Header {
	merged := http.Header{}
	for _, h := range headers {
		for k, vals := range h {
			for _, v := range vals {
				merged.Add(k, v)
			}
		}
	}
	return merged
}

// HeaderString converts a `http.Header` to a string.
func HeaderString(h http.Header) (string, error) {
	b := bytes.NewBuffer([]byte{})
	err := h.Write(b)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
