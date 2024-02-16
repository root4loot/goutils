package netutil

import (
	"net"
	"time"
)

// CanDial checks if a connection can be established to the given host.
func CanDial(host string, port string) bool {
	conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// CanDialWithTimeout checks if a connection can be established to the given host with a timeout.
func CanDialWithTimeout(host string, port string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
