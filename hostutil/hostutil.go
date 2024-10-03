package hostutil

import (
	"strings"
)

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
