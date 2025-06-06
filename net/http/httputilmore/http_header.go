package httputilmore

import (
	"bytes"
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
	HeaderContentID                    = "Content-ID"
	HeaderContentLanguage              = "Content-Language"
	HeaderContentLength                = "Content-Length"
	HeaderContentMD5                   = "Content-MD5"
	HeaderContentTransferEncoding      = "Content-Transfer-Encoding"
	HeaderContentType                  = "Content-Type"
	ContentTypeAppJSON                 = "application/json"
	ContentTypeAppJSONUtf8             = "application/json; charset=utf-8"
	ContentTypeAppOctetStream          = "application/octet-stream"
	ContentTypeAppFormURLEncoded       = "application/x-www-form-urlencoded"
	ContentTypeAppFormURLEncodedUtf8   = "application/x-www-form-urlencoded; charset=utf-8"
	ContentTypeAppXML                  = "application/xml"
	ContentTypeAppXMLUtf8              = "application/xml; charset=utf-8"
	ContentTypeImageGIF                = "image/gif"
	ContentTypeImageJPEG               = "image/jpeg"
	ContentTypeImagePNG                = "image/png"
	ContentTypeImageSVG                = "image/svg+xml"
	ContentTypeImageWebP               = "image/webp"
	ContentTypeMultipartAlternative    = "multipart/alternative"
	ContentTypeMultipartFormData       = "multipart/form-data"
	ContentTypeMultipartMixed          = "multipart/mixed"
	ContentTypeTextCalendarUtf8Request = "text/calendar; charset=utf-8; method=REQUEST"
	ContentTypeTextHTML                = "text/html"
	ContentTypeTextHTMLUtf8            = "text/html; charset=utf-8"
	ContentTypeTextMarkdown            = "text/markdown"
	ContentTypeTextPlain               = "text/plain"
	ContentTypeTextPlainUsASCII        = "text/plain; charset=us-ascii"
	ContentTypeTextPlainUtf8           = "text/plain; charset=utf-8"
	ContentTypeTextXMLUtf8             = "text/xml; charset=utf-8"
	DispositionTypeAttachment          = "attachment"
	DispositionTypeFormData            = "form-data"
	DispositionTypeInline              = "inline"
	HeaderDate                         = "Date"
	HeaderExpect                       = "Expect"
	Expect100Continue                  = "100-continue"
	HeaderIfMatch                      = "If-Match"
	HeaderIfModifiedSince              = "If-Modified-Since"
	HeaderIfNoneMatch                  = "If-None-Match"
	HeaderIfUnmodifiedSince            = "If-Unmodified-Since"
	HeaderLocation                     = "Location"
	HeaderPrefer                       = "Prefer" // IETF RFC-7240
	PreferHandlingStrict               = "handling=strict"
	PreferRespondAsync                 = "respond-async"
	PreferReturnMinimal                = "return=minimal"
	PreferReturnRepresentation         = "return=representation"
	HeaderRange                        = "Range"
	HeaderUserAgent                    = "User-Agent"
	HeaderWWWAuthenticate              = "WWW-Authenticate"
	WWWAuthenticateBasicRestricted     = "Basic realm=Restricted"
	HeaderXContentTypeOptions          = "X-Content-Type-Options"

	SchemeHTTPS = "https"

	HeaderNgrokSkipBrowserWarning      = "ngrok-skip-browser-warning" // header needs to be present. value can be anything. See more at: https://stackoverflow.com/questions/73017353/how-to-bypass-ngrok-browser-warning
	HeaderNgrokSkipBrowserWarningValue = "skip-browser-warning"
)

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

func HeaderContentTypeContains(h http.Header, ct string) bool {
	return strings.Contains(
		strings.ToLower(strings.TrimSpace(h.Get(HeaderContentType))),
		ct)
}
