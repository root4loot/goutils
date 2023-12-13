package domainutil

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// ResolveDomain resolves a domain name to an IP address (IPv4 or IPv6).
// It returns the IP address and any error encountered.
func ResolveDomain(domain string) (string, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return "", err
	}

	for _, ip := range ips {
		if ip4 := ip.To4(); ip4 != nil {
			return ip4.String(), nil // IPv4
		} else {
			return ip.String(), nil // IPv6
		}
	}

	return "", fmt.Errorf("no IP addresses found for domain: %s", domain)
}

// ResolveDomainWithTimeout resolves a domain name to an IP address (IPv4 or IPv6) with a specified timeout.
// It returns the IP address and any error encountered.
func ResolveDomainWithTimeout(domain string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ips, err := net.DefaultResolver.LookupIP(ctx, "ip", domain)
	if err != nil {
		return "", err
	}

	for _, ip := range ips {
		if ip4 := ip.To4(); ip4 != nil {
			return ip4.String(), nil // IPv4
		} else {
			return ip.String(), nil // IPv6
		}
	}

	return "", fmt.Errorf("no IP addresses found for domain: %s", domain)
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
