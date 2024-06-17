package urlutil

import (
	"net"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/root4loot/goutils/domainutil"
	"github.com/root4loot/goutils/iputil"
	"github.com/root4loot/goutils/netutil"
)

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
	ip, err := domainutil.ResolveDomainWithTimeout(u.Hostname(), timeout)
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
func HasScheme(rawURL string) bool {
	re := regexp.MustCompile(`^\w+?:\/\/\w+`)
	return re.MatchString(rawURL)
}

// HasFileExtension checks if the given rawURL string has a file extension in its path
func HasFileExtension(rawURL string) bool {
	u, _ := url.Parse(rawURL)
	return HasFileExtensionParsed(u)
}

// HasFileExtensionParsed checks if the given parsed URL has a file extension in its path
func HasFileExtensionParsed(u *url.URL) bool {
	// Split the path and look for the first instance of "."
	segments := strings.Split(u.Path, "/")
	for _, segment := range segments {
		if strings.Contains(segment, ".") {
			ext := filepath.Ext(segment)
			if ext != "" {
				return true
			}
		}
	}

	return false
}

// HasParam checks if a rawURL string has parameters
func HasParam(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	return HasParamParsed(u)
}

// HasParamParsed checks if a parsed URL has parameters
func HasParamParsed(u *url.URL) bool {
	return u.RawQuery != ""
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

// GetOrigin returns the origin of a URL.
func GetOrigin(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	return u.Scheme + "://" + u.Host, nil
}
