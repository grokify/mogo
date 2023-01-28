package netutil

import "net"

const (
	HostLocalhost         = "localhost"
	HostLoopbackIPv4      = "127.0.0.1"
	HostLoopbackIPv6      = "0:0:0:0:0:0:0:1"
	HostLoopbackIPv6Short = "::1"
)

var (
	IPv4loopback = net.IPv4(127, 0, 0, 1)
)
