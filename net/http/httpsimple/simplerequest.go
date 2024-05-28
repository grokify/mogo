package httpsimple

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/grokify/mogo/encoding/xmlutil"
	"github.com/grokify/mogo/net/http/httputilmore"
	"github.com/grokify/mogo/net/urlutil"
	"github.com/grokify/mogo/reflect/reflectutil"
)

const (
	BodyTypeFile = "file" // Body must be an `io.Reader`. Used for streaming.
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
		if ct := DefaultContentTypeBodyType(req.BodyType); ct != "" {
			req.Headers.Add(httputilmore.HeaderContentType, ct)
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

func (req *Request) FullURL() (*url.URL, error) {
	return urlutil.URLStringAddQuery(req.URL, req.Query, true)
}

func (req *Request) HTTPRequest(ctx context.Context) (*http.Request, error) {
	bodyBytes, err := req.BodyBytes()
	if err != nil {
		return nil, err
	}
	u, err := req.FullURL()
	if err != nil {
		return nil, err
	}
	httpreq, err := http.NewRequestWithContext(ctx, req.Method, u.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	for k, vals := range req.Headers {
		for _, v := range vals {
			httpreq.Header.Add(k, v)
		}
	}
	if httpreq.Header.Get(httputilmore.HeaderContentType) == "" {
		if ct := DefaultContentTypeBodyType(req.BodyType); ct != "" {
			httpreq.Header.Add(httputilmore.HeaderContentType, ct)
		}
	}
	return httpreq, nil
}

func DefaultContentTypeBodyType(bt string) string {
	bt = strings.ToLower(strings.TrimSpace(bt))
	switch bt {
	case BodyTypeForm:
		return httputilmore.ContentTypeAppFormURLEncodedUtf8
	case BodyTypeJSON:
		return httputilmore.ContentTypeAppJSONUtf8
	case BodyTypeXML:
		return httputilmore.ContentTypeAppXMLUtf8
	}
	return ""
}
