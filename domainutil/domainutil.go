package domainutil

import (
	"net/url"
	"regexp"
	"strings"
)

// IsDomainName checks if a string is a valid domain name.
func IsDomainName(str string) bool {
	// Regular expression pattern to match domain names with optional wildcards
	regex := regexp.MustCompile(`^(?i)(\*\.){0,1}(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.){1,}(?:[a-z]{2,})$`)

	return regex.MatchString(str)
}

// IsValidDomain checks if a domain is valid by attempting to parse it with a custom scheme.
func IsValidDomain(domain string) bool {
	if !IsDomainName(domain) {
		return false // If the domain doesn't match the regex, return false immediately
	}

	// Parse the domain with a custom scheme
	u, err := url.Parse("valid://" + domain)
	if err != nil {
		return false // If parsing fails, return false
	}

	// Check if the host is non-empty and there's no Path or RawQuery
	return u.Host != "" && u.Path == "" && u.RawQuery == ""
}

// GetDomainRoot returns the root domain of a domain
func GetDomainRoot(domain string) string {
	r, _ := regexp.Compile(`\w+\.\w+$`)
	m := r.FindString(domain)
	return strings.ToLower(m)
}

// DomainRoots returns a list of unique root domains for a slice of domains
func DomainRoots(items []string) (roots []string) {
	for _, item := range items {
		root := GetDomainRoot(item)
		roots = append(roots, root)
	}
	return unique(roots)
}

func unique(items []string) (uniqueItems []string) {
	uniqueMap := make(map[string]bool)
	for _, item := range items {
		uniqueMap[item] = true
	}
	for item := range uniqueMap {
		uniqueItems = append(uniqueItems, item)
	}
	return
}
