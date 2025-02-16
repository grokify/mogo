package httputilmore

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/errors/errorsutil"
)

func ResponseIsContentType(ct string, r *http.Response) bool {
	if r == nil {
		return false
	}
	ct = strings.ToLower(strings.TrimSpace(ct))
	ctv := strings.ToLower(r.Header.Get(HeaderContentType))
	if ct == "" {
		if ct == ctv {
			return true
		} else {
			return false
		}
	}
	return strings.Index(ctv, ct) == 0
}

// ProxyResponse copies the information from a `*http.Response` to a
// `http.ResponseWriter`.
func ProxyResponse(w http.ResponseWriter, resp *http.Response) ([]byte, error) {
	if resp == nil {
		return nil, errors.New("E_NIL_HTTP_RESPONSE")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	_, err = w.Write(body)
	if err != nil {
		return body, err
	}
	headers := []string{
		HeaderContentEncoding,
		HeaderContentLanguage,
		HeaderContentTransferEncoding,
		HeaderContentType}
	for _, header := range headers {
		if len(resp.Header.Get(header)) > 0 {
			w.Header().Add(header, resp.Header.Get(header))
		}
	}
	w.WriteHeader(resp.StatusCode)
	return body, nil
}

func CondenseResponseNot2xxToError(resp *http.Response, err error, msg string) error {
	if err != nil {
		if len(msg) > 0 {
			return errorsutil.Wrap(err, msg)
		} else {
			return err
		}
	} else if resp == nil {
		return errors.New("*http.Response_is_nil")
	} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
		if len(msg) > 0 {
			msg += ": "
		}
		more := []string{}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			more = append(more, err.Error())
		} else {
			more = append(more, string(body))
		}
		moreString := ""
		jsn, err := json.Marshal(more)
		if err != nil {
			moreString = err.Error()
		} else {
			moreString = string(jsn)
		}
		return fmt.Errorf("non_2xx_status_code [%d] [%s] [%s]", resp.StatusCode, msg, moreString)
	}
	return nil
}

func ResponseBodyMore(r *http.Response, jsonPrefix, jsonIndent string) ([]byte, error) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return []byte{}, err
	} else if !HeaderContentTypeContains(r.Header, ContentTypeAppJSON) {
		return b, nil
	} else if jsonutil.IsJSON(b) {
		if jsonPrefix != "" || jsonIndent != "" {
			return jsonutil.IndentBytes(b, jsonPrefix, jsonIndent)
		} else {
			return b, nil
		}
	} else {
		return b, nil
	}
}

// ResponseInfo is a generic struct to handle response info.
type ResponseInfo struct {
	Name       string            `json:"name,omitempty"` // to distinguish from other requests
	Method     string            `json:"method,omitempty"`
	URL        string            `json:"url,omitempty"`
	StatusCode int               `json:"statusCode,omitempty"`
	Time       time.Time         `json:"time,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
}

// ToJSON returns ResponseInfo as a JSON byte array, embedding json.Marshal
// errors if encountered.
func (resIn *ResponseInfo) ToJSON() []byte {
	bytes, err := json.Marshal(resIn)
	if err != nil {
		resIn2 := ResponseInfo{StatusCode: 500, Body: err.Error()}
		bytes, _ := json.Marshal(resIn2)
		return bytes
	}
	return bytes
}

func ResponseWriterWriteJSON(w http.ResponseWriter, statusCode int, body any, prefix, indent string) error {
	if w == nil {
		return errors.New("nil response writer")
	}
	wroteStatus := false
	var errs []error
	if body != nil {
		bytes, err := jsonutil.MarshalSimple(body, "", "  ")
		if err != nil {
			errs = append(errs, err)
			gr := ResponseInfo{
				Body: err.Error()}
			bytes, err = jsonutil.MarshalSimple(gr, "", "  ")
			if err != nil {
				errs = append(errs, err)
			}
			w.WriteHeader(500)
			_, err2 := w.Write(bytes)
			errs = append(errs, err2)
			return errorsutil.Join(false, errs...)
		}
		_, err = w.Write(bytes)
		if err != nil {
			errs = append(errs, err)
			gr := ResponseInfo{
				Body: err.Error()}
			bytes, err = jsonutil.MarshalSimple(gr, "", "  ")
			if err != nil {
				errs = append(errs, err)
			}
			w.WriteHeader(500)
			_, err2 := w.Write(bytes)
			errs = append(errs, err2)
			return errorsutil.Join(false, errs...)
		}
	}
	if !wroteStatus && statusCode != 0 {
		w.WriteHeader(statusCode)
	}
	return nil
}
