package httputilmore

// Constants ensuring that header names are correctly spelled and consistently cased.
const (
	HeaderAccept                       = "Accept"
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
	ContentTypeAppJsonUtf8             = "application/json; charset=utf-8"
	ContentTypeAppFormUrlEncoded       = "application/x-www-form-urlencoded"
	ContentTypeAppXml                  = "application/xml"
	ContentTypeAppXmlUtf8              = "application/xml; charset=utf-8"
	ContentTypeTextCalendarUtf8Request = `text/calendar; charset="utf-8"; method=REQUEST`
	ContentTypeTextHtmlUtf8            = "text/html; charset=utf-8"
	ContentTypeTextPlainUsAscii        = "text/plain; charset=us-ascii"
	ContentTypeTextPlainUtf8           = "text/plain; charset=utf-8"
	ContentTypeTextXmlUtf8             = "text/xml; charset=utf-8"
	SchemeHTTPS                        = "https"
)
