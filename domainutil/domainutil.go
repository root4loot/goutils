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

// ResolveDomainWithTimeout resolves a domain name to an IP address (IPv4 or IPv6) with a timeout.
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
	var hostnameRegex = regexp.MustCompile(`^(?i)([a-z0-9]([-a-z0-9]*[a-z0-9])?\.)*[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
	return hostnameRegex.MatchString(host)
}

// IsValidDomain checks if a string is a valid domain name.
func IsValidDomain(domain string) bool {
	if !IsDomainName(domain) {
		return false
	}
	u, err := url.Parse("valid://" + domain)
	if err != nil {
		return false
	}

	return u.Host != "" && u.Path == "" && u.RawQuery == ""
}

// GetRootDomain returns the root domain of a domain
func GetRootDomain(domain string) string {
	r, _ := regexp.Compile(`\w+\.\w+$`)
	m := r.FindString(domain)
	return strings.ToLower(m)
}

// GetRootDomains returns a list of unique root domains for a slice of domains
func GetRootDomains(items []string) (rootDomains []string) {
	for _, item := range items {
		root := GetRootDomain(item)
		rootDomains = append(rootDomains, root)
	}
	return unique(rootDomains)
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
