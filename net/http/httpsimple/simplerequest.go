package httpsimple

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/grokify/mogo/encoding/xmlutil"
	"github.com/grokify/mogo/net/http/httputilmore"
	"github.com/grokify/mogo/reflect/reflectutil"
)

const (
	headerContentType = httputilmore.HeaderContentType

	BodyTypeForm = "form"
	BodyTypeJSON = "json"
	BodyTypeXML  = "xml"
)

type SimpleRequest struct {
	Method        string
	URL           string
	Query         map[string][]string
	Headers       map[string][]string
	Body          any
	BodyType      string
	AddXMLDocType bool // only used if `Body` is a struct.
}

func (req *SimpleRequest) Inflate() {
	req.Method = strings.ToUpper(strings.TrimSpace(req.Method))
	if len(strings.TrimSpace(req.Method)) == 0 {
		req.Method = http.MethodGet
	} else {
		req.Method = strings.ToUpper(strings.TrimSpace(req.Method))
	}
	if req.Headers == nil {
		req.Headers = map[string][]string{}
	}
	if len(strings.TrimSpace(req.BodyType)) == 0 {
		headerCTLc := strings.ToLower(httputilmore.HeaderContentType)
		haveCT := false
		for hkey := range req.Headers {
			hkeyCTLc := strings.ToLower(hkey)
			if hkeyCTLc == headerCTLc {
				haveCT = true
				break
			}
		}
		if !haveCT {
			switch req.BodyType {
			case BodyTypeForm:
				req.Headers[headerContentType] = []string{httputilmore.ContentTypeAppFormURLEncodedUtf8}
			case BodyTypeJSON:
				req.Headers[headerContentType] = []string{httputilmore.ContentTypeAppJSONUtf8}
			case BodyTypeXML:
				req.Headers[headerContentType] = []string{httputilmore.ContentTypeAppXMLUtf8}
			}
		}
	}
}

func (req *SimpleRequest) BodyBytes() ([]byte, error) {
	if req.Body == nil {
		return []byte{}, nil
	} else if reqBodyBytes, ok := req.Body.([]byte); ok {
		return reqBodyBytes, nil
	} else if reqBodyString, ok := req.Body.(string); ok {
		return []byte(reqBodyString), nil
	} else if req.BodyType == BodyTypeJSON {
		return json.Marshal(req.Body)
	} else if req.BodyType == BodyTypeXML {
		return xmlutil.MarshalIndent(req.Body, "", "", req.AddXMLDocType)
	} else if req.BodyType == BodyTypeForm {
		if v, ok := req.Body.(url.Values); ok {
			return []byte(v.Encode()), nil
		}
	}
	return []byte{}, fmt.Errorf("body type (%s) not supported", reflectutil.TypeName(req.Body))
}
