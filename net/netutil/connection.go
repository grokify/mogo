package netutil

import (
	"bufio"
	"errors"
	"net"
	"net/http"

	"github.com/grokify/mogo/errors/errorsutil"
)

// ModifyConnectionRequest updates the HTTP request for network connection.
func ModifyConnectionRequest(conn net.Conn, modRequest func(r *http.Request) error) error {
	// Code adapted from: https://stackoverflow.com/a/76684845/1908967
	if conn == nil {
		return errors.New("net.Conn cannot be nil")
	}
	defer conn.Close()

	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		return errorsutil.Wrap(err, "error reading request")
	}

	// Modify the request as needed
	if modRequest != nil {
		if err := modRequest(req); err != nil {
			return errorsutil.Wrap(err, "error modifying request")
		}
	}

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return errorsutil.Wrap(err, "error sending request")
	}
	defer resp.Body.Close()

	return resp.Write(conn)
}
