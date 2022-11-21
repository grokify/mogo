package httputilmore

import (
	"fmt"
	"net/http"
	"strings"
)

type HTTPMethod string

const (
	MethodConnect HTTPMethod = http.MethodConnect
	MethodDelete             = http.MethodDelete
	MethodGet                = http.MethodGet
	MethodHead               = http.MethodHead
	MethodOptions            = http.MethodOptions
	MethodPatch              = http.MethodPatch
	MethodPost               = http.MethodPost
	MethodPut                = http.MethodPut
	MethodTrace              = http.MethodTrace
)

// ParseHTTPMethod returns a HTTPMethod type for a string.
func ParseHTTPMethod(method string) (HTTPMethod, error) {
	methodCanonical := strings.ToUpper(strings.TrimSpace(method))
	switch methodCanonical {
	case http.MethodConnect:
		return MethodConnect, nil
	case http.MethodDelete:
		return MethodDelete, nil
	case http.MethodGet:
		return MethodGet, nil
	case http.MethodHead:
		return MethodHead, nil
	case http.MethodOptions:
		return MethodOptions, nil
	case http.MethodPatch:
		return MethodPatch, nil
	case http.MethodPost:
		return MethodPost, nil
	case http.MethodPut:
		return MethodPut, nil
	case http.MethodTrace:
		return MethodTrace, nil
	}
	return MethodConnect, fmt.Errorf("cannot parse method [%v]", method)
}

// ParseHTTPMethodString returns a HTTPMethod as a string for a string.
func ParseHTTPMethodString(method string) (string, error) {
	methodCanonical, err := ParseHTTPMethod(method)
	return string(methodCanonical), err
}

func MethodsMap() map[string]int {
	return map[string]int{
		http.MethodConnect: 1,
		http.MethodDelete:  1,
		http.MethodGet:     1,
		http.MethodHead:    1,
		http.MethodOptions: 1,
		http.MethodPatch:   1,
		http.MethodPost:    1,
		http.MethodPut:     1,
		http.MethodTrace:   1,
	}
}
