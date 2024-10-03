package hostutil

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// NormalizeHost takes a host input (either a domain name or hostname) and returns it in a standardized format.
// It converts the input to lowercase, removes redundant ports (e.g., 443 for HTTPS and 80 for HTTP),
// validates that the input is a valid fully qualified domain name (FQDN), and ensures that the hostname is properly formatted.
func NormalizeHost(host string) (string, error) {
	host = strings.TrimSpace(strings.ToLower(host))

	if strings.Contains(host, "://") || strings.Contains(host, "/") {
		return "", fmt.Errorf("input contains URL scheme or path: %s", host)
	}

	hostname, port, err := net.SplitHostPort(host)
	if err != nil {
		if strings.Contains(host, ":") {
			return "", fmt.Errorf("failed to parse host input: %s", host)
		}
		hostname = host
	}

	if port != "" {
		if _, err := strconv.Atoi(port); err != nil {
			return "", fmt.Errorf("invalid port number: %s", port)
		}
	}

	if net.ParseIP(hostname) == nil && !IsValidHostname(hostname) {
		return "", fmt.Errorf("invalid hostname or IP address: %s", hostname)
	}

	if !strings.Contains(hostname, ".") {
		return "", fmt.Errorf("invalid FQDN: %s", hostname)
	}

	switch port {
	case "443", "80":
		port = ""
	}

	if port != "" {
		hostname = net.JoinHostPort(hostname, port)
	}

	return hostname, nil
}

// IsValidHostname checks if the given hostname is valid based on RFC 1123.
func IsValidHostname(hostname string) bool {
	if len(hostname) == 0 || len(hostname) > 255 {
		return false
	}
	labels := strings.Split(hostname, ".")
	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 {
			return false
		}
		startChar := label[0]
		endChar := label[len(label)-1]
		if !((startChar >= 'a' && startChar <= 'z') || (startChar >= '0' && startChar <= '9')) {
			return false
		}
		if !((endChar >= 'a' && endChar <= 'z') || (endChar >= '0' && endChar <= '9')) {
			return false
		}
		for i := 0; i < len(label); i++ {
			c := label[i]
			if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' {
				continue
			}
			return false
		}
	}
	return true
}
