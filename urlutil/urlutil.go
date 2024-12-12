package urlutil

import (
	"fmt"
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

// IsValidURL checks if a string is a valid URL.
func IsValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// CanReachURL checks if a URL can be reached without a timeout.
func CanReachURL(rawURL string) error {
	var err error

	rawURL = EnsurePortIsSet(rawURL)
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

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

	if iputil.IsURLIP(rawURL) {
		if netutil.CanDialWithTimeout(u.Hostname(), u.Port(), timeout) {
			return err
		}
	}

	ip, err := domainutil.ResolveDomainWithTimeout(u.Hostname(), timeout)
	if err != nil {
		return err
	}

	if netutil.CanDialWithTimeout(ip, u.Port(), timeout) {
		return err
	}

	return err
}

// EnsurePortIsSet ensures a URL has a port set. If no port is provided, it defaults to 80 for HTTP and 443 for HTTPS.
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

// EnsureHTTP ensures that the URL has an HTTP scheme, replacing any existing scheme.
func EnsureHTTP(rawURL string) string {
	if !HasScheme(rawURL) {
		return "http://" + rawURL
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "http://" + rawURL
	}

	u.Scheme = "http"
	return u.String()
}

// EnsureHTTPS ensures that the URL has an HTTPS scheme, replacing any existing scheme.
func EnsureHTTPS(rawURL string) string {
	if !HasScheme(rawURL) {
		return "https://" + rawURL
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "https://" + rawURL
	}

	u.Scheme = "https"
	return u.String()
}

// HasFileExtension checks if the given rawURL string has a file extension in its path
func HasFileExtension(rawURL string) bool {
	u, _ := url.Parse(rawURL)
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
	return u.RawQuery != ""
}

// GetExt returns the file extension of a raw URL string
func GetExt(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	segments := strings.Split(u.Path, "/")
	if len(segments) > 0 {
		lastSegment := segments[len(segments)-1]
		return strings.ToLower(filepath.Ext(lastSegment))
	}
	return ""
}

// EnsureTrailingSlash ensures a URL has a trailing slash
func EnsureTrailingSlash(rawURL string) string {
	parsedURL, _ := url.Parse(rawURL)

	re := regexp.MustCompile(`[\W_]$`)
	if filepath.Ext(parsedURL.Path) == "" && !re.MatchString(parsedURL.Path) && !strings.HasSuffix(parsedURL.Path, "/") {
		parsedURL.Path += "/"
	}

	return parsedURL.String()
}

// IsMediaExt checks if a file extension is a media type
func IsMediaExt(ext string) bool {
	ext = strings.ToLower(ext)
	for _, mediaExt := range GetMediaExtensions() {
		if ext == mediaExt {
			return true
		}
	}
	return false
}

// GetOrigin returns the origin of a URL.
func GetOrigin(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	return u.Scheme + "://" + u.Host, nil
}

// NormalizeSlashes trims double slashes from a URL, preserving the initial scheme
func NormalizeSlashes(rawURL string) string {
	if strings.HasPrefix(rawURL, "http://") {
		rawURL = "http://" + strings.Replace(rawURL[len("http://"):], "//", "/", -1)
	} else if strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + strings.Replace(rawURL[len("https://"):], "//", "/", -1)
	} else {
		rawURL = strings.Replace(rawURL, "//", "/", -1)
	}
	return rawURL
}

// RemoveTrailingSlash removes an unwanted "/" suffix from a URL
func RemoveTrailingSlash(rawURL string) string {
	if strings.HasSuffix(rawURL, "/") {
		return rawURL[:len(rawURL)-1]
	}
	return rawURL
}

// GetMediaExtensions returns a slice of common media file extensions
func GetMediaExtensions() []string {
	return []string{
		".png", ".jpg", ".jpeg", ".woff", ".woff2", ".ttf", ".eot", ".svg", ".gif", ".ico", ".webp",
		".mp4", ".webm", ".mp3", ".wav", ".flac", ".aac", ".ogg", ".m4a", ".flv", ".avi", ".mov",
		".wmv", ".swf", ".mkv", ".m4v", ".3gp", ".3g2",
	}
}

// RemoveDefaultPort removes the default port from a URL based on its scheme.
func RemoveDefaultPort(urlStr string) (string, error) {
	isTempScheme := false

	if !strings.Contains(urlStr, "://") {
		urlStr = "temp://" + urlStr
		isTempScheme = true
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	host := u.Host
	var port string

	if strings.Contains(u.Host, ":") {
		host, port, err = net.SplitHostPort(u.Host)
		if err != nil {
			// Handle IPv6 addresses enclosed in brackets
			if strings.HasPrefix(u.Host, "[") && strings.Contains(u.Host, "]") {
				hostPort := strings.TrimPrefix(u.Host, "[")
				hostPort = strings.Replace(hostPort, "]", "", 1)
				host, port, err = net.SplitHostPort(hostPort)
				host = "[" + host + "]"
				if err != nil {
					return "", fmt.Errorf("invalid host: %w", err)
				}
			} else {
				return "", fmt.Errorf("invalid host: %w", err)
			}
		}
	}

	defaultPort := ""
	switch u.Scheme {
	case "http":
		defaultPort = "80"
	case "https":
		defaultPort = "443"
	case "ftp":
		defaultPort = "21"
	case "temp":
		defaultPort = ""
	default:
		return u.String(), nil
	}

	if port == defaultPort {
		u.Host = host
	}

	finalURL := u.String()
	if isTempScheme {
		finalURL = strings.TrimPrefix(finalURL, "temp://")
	}

	return finalURL, nil
}
