package multipartutil

import "github.com/grokify/mogo/net/http/httputilmore"

// NewPartAlternativeOrNot can be used for email bodies which come with
// text and HTML alternatives.
func NewPartAlternativeOrNot(text, html []byte) (Part, error) {
	if len(text) > 0 && len(html) > 0 {
		mps := NewMultipartSimpleAlternative(text, html)
		return mps.Part()
	} else if len(html) > 0 {
		return Part{
			Type:         PartTypeRaw,
			ContentType:  httputilmore.ContentTypeTextHTMLUtf8,
			Base64Encode: false,
			RawBody:      html,
		}, nil
	} else {
		return Part{
			Type:         PartTypeRaw,
			ContentType:  httputilmore.ContentTypeTextPlainUtf8,
			Base64Encode: false,
			RawBody:      text,
		}, nil
	}
}

func NewMultipartSimpleAlternative(text, html []byte) MultipartSimple {
	return MultipartSimple{
		ContentType: httputilmore.ContentTypeMultipartAlternative,
		Parts: []Part{
			{
				Type:         PartTypeRaw,
				ContentType:  httputilmore.ContentTypeTextPlain,
				Base64Encode: false,
				RawBody:      text,
			}, {
				Type:         PartTypeRaw,
				ContentType:  httputilmore.ContentTypeTextHTMLUtf8,
				Base64Encode: false,
				RawBody:      html,
			},
		},
	}
}
