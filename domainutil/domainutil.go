package domainutil

import (
	"net/url"
	"regexp"
)

// IsDomainName checks if a string is a valid domain name.
func IsDomainName(str string) bool {
	// Regular expression pattern to match domain names with optional wildcards
	regex := regexp.MustCompile(`^(?i)(\*\.){0,1}(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.){1,}(?:[a-z]{2,})$`)

	return regex.MatchString(str)
}

// IsValidDomain checks if a domain is valid by attempting to parse it with a custom scheme.
func IsValidDomain(domain string) bool {
	u, err := url.Parse("valid://" + domain)
	if err != nil {
		return false
	}

	// Check if the host is non-empty and there's no Path or RawQuery
	return u.Host != "" && u.Path == "" && u.RawQuery == ""
}
