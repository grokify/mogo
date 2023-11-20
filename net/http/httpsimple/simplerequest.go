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
	// headerContentType = httputilmore.HeaderContentType
	BodyTypeForm = "form"
	BodyTypeJSON = "json"
	BodyTypeXML  = "xml"
)

type Request struct {
	Method        string
	URL           string
	Query         url.Values
	Headers       http.Header
	Body          any
	BodyType      string
	AddXMLDocType bool // only used if `Body` is a struct.
}

func NewRequest() Request {
	return Request{
		Query:   url.Values{},
		Headers: http.Header{},
	}
}

func (req *Request) Inflate() {
	req.Method = strings.ToUpper(strings.TrimSpace(req.Method))
	if req.Method == "" {
		req.Method = http.MethodGet
	}
	if req.Headers == nil {
		req.Headers = http.Header{}
	}
	if strings.TrimSpace(req.Headers.Get(httputilmore.HeaderContentType)) == "" {
		switch req.BodyType {
		case BodyTypeForm:
			req.Headers.Add(httputilmore.HeaderContentType, httputilmore.ContentTypeAppFormURLEncodedUtf8)
		case BodyTypeJSON:
			req.Headers.Add(httputilmore.HeaderContentType, httputilmore.ContentTypeAppJSONUtf8)
		case BodyTypeXML:
			req.Headers.Add(httputilmore.HeaderContentType, httputilmore.ContentTypeAppXMLUtf8)
		}
	}
}

func (req *Request) BodyBytes() ([]byte, error) {
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
	return []byte{}, fmt.Errorf("body type (%s) not supported", reflectutil.NameOf(req.Body, true))
}
