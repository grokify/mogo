package anyhttp

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"

	hum "github.com/grokify/gotilla/net/httputilmore"
)

type RequestNetHttp struct {
	Raw                 *http.Request
	allArgs             *ArgsUrlValues
	postArgs            *ArgsUrlValues
	multipartForm       *multipart.Form
	parsedMultipartForm bool
	parsedFormArgs      bool
}

func NewRequestNetHttp(req *http.Request) *RequestNetHttp {
	return &RequestNetHttp{
		Raw:      req,
		allArgs:  &ArgsUrlValues{Raw: req.Form},
		postArgs: &ArgsUrlValues{Raw: req.PostForm}}
}

func (r *RequestNetHttp) ParseForm() error {
	if r.parsedFormArgs {
		return nil
	}
	r.parsedFormArgs = true
	if err := r.Raw.ParseForm(); err != nil {
		return err
	}
	r.allArgs = &ArgsUrlValues{r.Raw.Form}
	r.postArgs = &ArgsUrlValues{r.Raw.PostForm}
	return nil
}

func (r RequestNetHttp) RemoteAddr() net.Addr {
	return Addr{Protocol: "tcp", Address: r.Raw.RemoteAddr}
}
func (r RequestNetHttp) RemoteAddress() string { return r.Raw.RemoteAddr }
func (r RequestNetHttp) UserAgent() []byte     { return []byte(r.Raw.UserAgent()) }
func (r RequestNetHttp) AllArgs() Args         { return r.allArgs }
func (r RequestNetHttp) QueryArgs() Args       { return r.allArgs }
func (r RequestNetHttp) PostArgs() Args        { return r.postArgs }
func (r RequestNetHttp) Method() []byte        { return []byte(r.Raw.Method) }
func (r RequestNetHttp) Header() http.Header   { return r.Raw.Header }
func (r RequestNetHttp) Form() url.Values      { return r.Raw.Form }
func (r RequestNetHttp) RequestURI() []byte    { return []byte(r.Raw.RequestURI) }

func (r *RequestNetHttp) MultipartForm() (*multipart.Form, error) {
	if !r.parsedMultipartForm {
		r.parsedMultipartForm = true
		if err := r.Raw.ParseMultipartForm(100000); err != nil {
			return nil, err
		}
	}
	return r.Raw.MultipartForm, nil
}

type ResponseNetHttp struct {
	Raw http.ResponseWriter
}

func NewResponseNetHttp(w http.ResponseWriter) ResponseNetHttp {
	return ResponseNetHttp{Raw: w}
}

func (w ResponseNetHttp) GetHeader(k string) []byte { return []byte(w.Raw.Header().Get(k)) }
func (w ResponseNetHttp) SetHeader(k, v string)     { w.Raw.Header().Set(k, v) }
func (w ResponseNetHttp) SetStatusCode(code int)    { w.Raw.WriteHeader(code) }
func (w ResponseNetHttp) SetContentType(ct string) {
	w.Raw.Header().Set(hum.HeaderContentType, ct)
}

func (w ResponseNetHttp) SetBodyBytes(body []byte) (int, error) {
	w.Raw.Write(body)
	return -1, nil
}
func (w ResponseNetHttp) SetBodyStream(bodyStream io.Reader, bodySize int) error {
	bytes, err := ioutil.ReadAll(bodyStream)
	if err != nil {
		return err
	}
	w.Raw.Write(bytes)
	return nil
}

func (w ResponseNetHttp) SetCookie(cookie *Cookie) {
	http.SetCookie(w.Raw, cookie.ToNetHttp())
}

func NewResReqNetHttp(res http.ResponseWriter, req *http.Request) (ResponseNetHttp, *RequestNetHttp) {
	return NewResponseNetHttp(res), NewRequestNetHttp(req)
}
