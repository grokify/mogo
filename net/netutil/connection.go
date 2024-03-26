package netutil

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

// ModifyConnectionRequest updates the HTTP request for network connection.
func ModifyConnectionRequest(conn net.Conn, modRequest func(r *http.Request)) {
	// Code adapted from: https://stackoverflow.com/a/76684845/1908967
	defer conn.Close()

	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		fmt.Println("Error reading request:", err)
		return
	}

	// Modify the request as needed
	if modRequest != nil {
		modRequest(req)
	}

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	resp.Write(conn)
}
