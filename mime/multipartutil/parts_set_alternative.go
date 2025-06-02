package multipartutil

import (
	"github.com/grokify/mogo/net/http/httputilmore"
)

// NewPartAlternativeOrNot can be used for email bodies which come with
// text and HTML alternatives.
func NewPartAlternativeOrNot(textBody, htmlBody []byte) (Part, error) {
	if len(textBody) > 0 && len(htmlBody) > 0 {
		mps := NewPartsSetAlternative(textBody, htmlBody)
		return mps.Part()
	} else if len(htmlBody) > 0 {
		return Part{
			Type:             PartTypeRaw,
			DispositionType:  httputilmore.DispositionTypeInline,
			ContentType:      httputilmore.ContentTypeTextHTMLUtf8,
			BodyEncodeBase64: false,
			BodyDataRaw:      htmlBody,
		}, nil
	} else {
		return Part{
			Type:             PartTypeRaw,
			DispositionType:  httputilmore.DispositionTypeInline,
			ContentType:      httputilmore.ContentTypeTextPlainUtf8,
			BodyEncodeBase64: false,
			BodyDataRaw:      textBody,
		}, nil
	}
}

func NewPartsSetAlternative(text, html []byte) PartsSet {
	return PartsSet{
		ContentType: httputilmore.ContentTypeMultipartAlternative,
		Parts: []Part{
			{
				Type:             PartTypeRaw,
				DispositionType:  httputilmore.DispositionTypeInline,
				ContentType:      httputilmore.ContentTypeTextPlain,
				BodyEncodeBase64: false,
				BodyDataRaw:      text,
			}, {
				Type:             PartTypeRaw,
				DispositionType:  httputilmore.DispositionTypeInline,
				ContentType:      httputilmore.ContentTypeTextHTMLUtf8,
				BodyEncodeBase64: false,
				BodyDataRaw:      html,
			},
		},
	}
}
