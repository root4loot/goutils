package domainutil

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/root4loot/goutils/iputil"
	"github.com/root4loot/goutils/netutil"
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

// IsURL checks if a string is a URL.
func IsURL(url string) bool {
	regex := regexp.MustCompile(`^(https?|ftp)://[^\s/$.?#].[^\s]*$`)
	return regex.MatchString(url)
}

// CanReachURL checks if a URL can be reached without a timeout.
func CanReachURL(rawURL string) error {
	var err error

	rawURL = EnsurePortIsSet(rawURL)
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	// Dial the host to check connectivity
	conn, err := net.Dial("tcp", u.Host)
	if err != nil {
		return err
	}
	defer conn.Close()

	return err
}

// CanReachURLWithTimeout checks if a URL can be reached with a specified timeout.
func CanReachURLWithTimeout(rawURL string, timeout time.Duration) error {
	var err error

	rawURL = EnsurePortIsSet(rawURL)
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	// check if URL is an IP address
	// if so, check if it can be dialed
	if iputil.IsURLIP(rawURL) {
		if netutil.CanDialWithTimeout(u.Hostname(), u.Port(), timeout) {
			return err
		}
	}

	// resolve the domain
	ip, err := ResolveDomainWithTimeout(u.Hostname(), timeout)
	if err != nil {
		return err
	}

	// check if the port can be dialed
	if netutil.CanDialWithTimeout(ip, u.Port(), timeout) {
		return err
	}

	return err
}

// EnsurePortIsSet takes a URL and ensures that a port is set, depending on the scheme.
// It returns the URL with the port set (if it was missing).
func EnsurePortIsSet(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	if u.Port() == "" {
		switch u.Scheme {
		case "http":
			u.Host = u.Hostname() + ":80"
		case "https":
			u.Host = u.Hostname() + ":443"
		}
	}

	return u.String()
}

// HasScheme checks if a URL has a scheme
func HasScheme(url string) bool {
	re := regexp.MustCompile(`^\w+?:\/\/\w+`)
	return re.MatchString(url)
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
