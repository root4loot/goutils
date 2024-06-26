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

// EnsureScheme ensures a URL has a scheme. If no scheme is provided, it defaults to "http".
func EnsureScheme(rawURL string, scheme ...string) string {
	if rawURL == "" {
		return rawURL
	}

	defaultScheme := "http"
	if len(scheme) > 0 && scheme[0] != "" {
		defaultScheme = scheme[0]
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		// If parsing fails, assume it's because of a missing scheme
		return defaultScheme + "://" + rawURL
	}

	if u.Scheme == "" {
		u.Scheme = defaultScheme
	}

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

// EnsureTrailingSlash appends a trailing slash to the URL path if it doesn't end in a file extension
// or with a non-alphanumeric symbol, and if it makes sense to do so.
func EnsureTrailingSlash(rawURL string) string {
	parsedURL, _ := url.Parse(rawURL)

	// Regex to check if the URL ends with a non-alphanumeric character
	re := regexp.MustCompile(`[\W_]$`)

	// Check if the path has a file extension, ends with a non-alphanumeric character, or already has a trailing slash
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
		// Replace all other occurrences of "//"
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
