package domainutil

import (
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

// IsURL checks if a string is a URL.
func IsURL(url string) bool {
	regex := regexp.MustCompile(`^(https?|ftp)://[^\s/$.?#].[^\s]*$`)
	return regex.MatchString(url)
}

// EnsureTrailingSlash appends a trailing slash to the URL path if it doesn't end in a file extension
// or with a symbol, and if it makes sense to do so.
func EnsureTrailingSlash(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// Regex to check if the URL ends with a non-alphanumeric character
	re := regexp.MustCompile(`[\W_]$`)

	// Check if the path has a file extension or ends with a symbol
	if filepath.Ext(parsedURL.Path) == "" && !re.MatchString(parsedURL.Path) && !strings.HasSuffix(parsedURL.Path, "/") {
		parsedURL.Path += "/"
	}

	return parsedURL.String(), nil
}

// IsDomainName checks if a string is a valid domain name.
func IsDomainName(str string) bool {
	regex := regexp.MustCompile(`^(?i)(\*\.){0,1}(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.){1,}(?:[a-z]{2,})$`)

	return regex.MatchString(str)
}

// IsHostname checks if the given target is a valid hostname
func IsHostname(host string) bool {
	// Regular expression for validating a hostname (simple version)
	var hostnameRegex = regexp.MustCompile(`^(?i)([a-z0-9]([-a-z0-9]*[a-z0-9])?\.)*[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
	return hostnameRegex.MatchString(host)
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
