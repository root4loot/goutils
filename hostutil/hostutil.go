package hostutil

import (
	"net"
	"strconv"
	"strings"
)

// IsValidHostname checks if the given hostname is a valid domain name based on RFC 1123.
// It ensures that the hostname does not consist entirely of numeric labels.
func IsValidHostname(hostname string) bool {
	if len(hostname) == 0 || len(hostname) > 255 {
		return false
	}
	labels := strings.Split(hostname, ".")
	allNumericLabels := true

	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 {
			return false
		}

		isNumericLabel := true
		for i := 0; i < len(label); i++ {
			c := label[i]
			if c < '0' || c > '9' {
				isNumericLabel = false
				break
			}
		}

		if !isNumericLabel {
			allNumericLabels = false
		}

		startChar := label[0]
		endChar := label[len(label)-1]
		if !((startChar >= 'a' && startChar <= 'z') || (startChar >= 'A' && startChar <= 'Z') || (startChar >= '0' && startChar <= '9')) {
			return false
		}
		if !((endChar >= 'a' && endChar <= 'z') || (endChar >= 'A' && endChar <= 'Z') || (endChar >= '0' && endChar <= '9')) {
			return false
		}

		for i := 0; i < len(label); i++ {
			c := label[i]
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' {
				continue
			}
			return false
		}
	}

	if allNumericLabels {
		return false
	}
	return true
}

// IsValidIP checks if the provided IP address is valid.
// It compares the parsed IP's string representation with the original input.
func IsValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	var formattedIP string

	if strings.Contains(ip, ".") {
		parsedIP = parsedIP.To4()
		if parsedIP == nil {
			return false
		}
		formattedIP = parsedIP.String()
	} else if strings.Contains(ip, ":") {
		parsedIP = parsedIP.To16()
		if parsedIP == nil {
			return false
		}
		formattedIP = parsedIP.String()
	} else {
		return false
	}

	return strings.EqualFold(ip, formattedIP)
}

// IsValidPort checks if the given port is a valid port number.
func IsValidPort(port string) bool {
	if p, err := strconv.Atoi(port); err == nil {
		return p > 0 && p <= 65535
	}
	return false
}

// IsValidHost checks if the given host is a valid hostname or IP address with an optional port.
func IsValidHost(host string) bool {
	hostPart, portPart, err := net.SplitHostPort(host)
	if err != nil {
		// If there's no port, validate the host as an IP address first, then as a hostname.
		return IsValidIP(host) || IsValidHostname(host)
	}

	if !IsValidPort(portPart) {
		return false
	}

	return IsValidIP(hostPart) || IsValidHostname(hostPart)
}
