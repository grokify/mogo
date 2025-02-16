package httputilmore

import (
	"errors"
	"net/http"
	"strconv"
)

var (
	ErrHTTPResponseCannotBeNil = errors.New("'*http.Response' cannot be nil. Expected a non-nil value")
)

var ErrStatus404 = errors.New("status " + strconv.Itoa(http.StatusNotFound) + " " + http.StatusText(http.StatusNotFound))

type HTTPError struct {
	HTTPStatus int    `json:"httpStatus"`
	Stage      string `json:"preOpPost"`
	Message    string `json:"errorMessage"`
}

func NewHTTPError(message string, httpStatus int, stage string) *HTTPError {
	return &HTTPError{
		Message:    message,
		HTTPStatus: httpStatus,
		Stage:      stage}
}

func (httperr *HTTPError) Bytes() []byte {
	bytes, err := json.Marshal(httperr)
	if err != nil {
		panic(err)
	}
	return bytes
}
