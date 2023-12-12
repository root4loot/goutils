package netutil

import (
	"net"
	"time"
)

// CanDial checks if a connection can be established to the given address.
func CanDial(ip string, port string) bool {
	conn, err := net.Dial("tcp", net.JoinHostPort(ip, port))
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// CanDialWithTimeout checks if a connection can be established to the given address with a timeout.
func CanDialWithTimeout(ip string, port string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
