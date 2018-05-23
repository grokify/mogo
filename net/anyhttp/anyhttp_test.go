package anyhttp

import (
	"testing"
)

// TestInterface ensures following interface.
func TestInterface(t *testing.T) {
	nethttpReq := &RequestNetHttp{}
	nethttpRes := ResponseNetHttp{}

	MockRequest(nethttpReq)
	MockResponse(nethttpRes)
	MockHandler(nethttpRes, nethttpReq)

	fasthttpReq := RequestFastHttp{}
	fasthttpRes := ResponseFastHttp{}

	MockRequest(fasthttpReq)
	MockResponse(fasthttpRes)
	MockHandler(fasthttpRes, fasthttpReq)
}

func MockRequest(aReq Request)                {}
func MockResponse(aReq Response)              {}
func MockHandler(aRes Response, aReq Request) {}
