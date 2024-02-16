package httputil

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

var ipCache sync.Map

// ClientWithOptionalResolvers creates an HTTP client with default and custom resolvers that are
// optionally provided. If no resolvers are provided, the default resolver is used,
// utilizing an IP cache to store and reuse resolved IP addresses.
func ClientWithOptionalResolvers(resolvers ...string) (*http.Client, error) {
	// Check if resolvers are provided
	if len(resolvers) > 0 {
		// Check if the resolvers contain the port :53, if not, add it
		for i, resolver := range resolvers {
			_, _, err := net.SplitHostPort(resolver)
			if err != nil {
				// Port is not specified, add :53
				resolvers[i] = resolver + ":53"
			}
		}
	} else {
		// No resolvers provided, use an empty slice
		resolvers = []string{}
	}

	defaultResolver := &net.Resolver{PreferGo: true}

	dialer := &net.Dialer{
		Resolver: defaultResolver,
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			// Check if the resolved IP address exists in the cache
			if cachedIP, ok := ipCache.Load(address); ok {
				if conn, ok := cachedIP.(net.Conn); ok {
					return conn, nil // Reuse the connection
				}
			}

			// IP address not found in cache, perform DNS lookup
			conn, err := dialer.DialContext(ctx, network, address)
			if err != nil {
				for _, resolver := range resolvers {
					d := net.Dialer{}
					conn, err = d.DialContext(ctx, network, resolver)
					if err != nil {
						continue // Try the next resolver
					}

					// Store the resolved IP address in the cache for future reuse
					ipCache.Store(address, conn)

					return conn, nil // Success, return the connection
				}
				return nil, fmt.Errorf("failed to resolve domain using custom resolvers")
			}

			// Store the resolved IP address in the cache for future reuse
			ipCache.Store(address, conn)

			return conn, nil // Success, return the connection
		},
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives:   true,
		MaxIdleConnsPerHost: -1,
	}

	client := &http.Client{
		Transport: transport,
	}

	return client, nil
}

// RedirectsToHTTPS checks for various types of redirects (Location header, Refresh header, and meta refresh tags)
// and returns true if there's a redirect to an HTTPS URL along with the final URL.
func RedirectsToHTTPS(httpURL string) (bool, string, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Don't follow redirects automatically
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(httpURL)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	finalURL := httpURL

	// Check for Location header redirect
	if location, ok := resp.Header["Location"]; ok && len(location) > 0 {
		finalURL = location[0]
		if strings.HasPrefix(location[0], "https://") {
			return true, finalURL, nil
		}
	}

	// Check for Refresh header redirect
	if refresh, ok := resp.Header["Refresh"]; ok && len(refresh) > 0 {
		// Refresh header format: "5;url=https://example.com/"
		parts := strings.SplitN(refresh[0], "url=", 2)
		if len(parts) == 2 && strings.HasPrefix(parts[1], "https://") {
			return true, parts[1], nil
		}
	}

	// For non-redirect responses or HTTP redirects, check for meta refresh tags in the HTML body
	if resp.StatusCode == 200 || (resp.StatusCode >= 300 && resp.StatusCode < 400) {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return false, "", err
		}

		// Parse the HTML body to find meta refresh tags
		metaURL, found := extractMetaRefreshURL(string(body)) // Adjusted to convert body to string and match the expected function signature
		if found && strings.HasPrefix(metaURL, "https://") {
			return true, metaURL, nil
		}
	}
	return false, finalURL, nil
}

func IsBinaryResponse(resp *http.Response) bool {
	if resp == nil || resp.Header == nil {
		return false
	}

	contentType := strings.ToLower(resp.Header.Get("Content-Type"))

	// List of common binary MIME types
	binaryMimes := []string{
		"application/octet-stream", "application/pdf", "application/zip",
		"application/x-rar-compressed", "application/x-7z-compressed",
		"application/x-tar", "application/gzip", "application/msword",
		"application/vnd.ms-excel", "application/vnd.ms-powerpoint",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"image/jpeg", "image/png", "image/gif", "image/webp",
		"image/tiff", "image/bmp", "video/mp4", "video/mpeg",
		"video/quicktime", "video/x-msvideo", "video/x-ms-wmv", "video/webm",
		"audio/mpeg", "audio/x-wav", "audio/ogg", "audio/mp4", "audio/webm",
		"application/x-binary", "application/x-shockwave-flash",
	}

	for _, mime := range binaryMimes {
		if strings.HasPrefix(contentType, mime) {
			return true
		}
	}
	return false
}

// extractMetaRefreshURL searches the HTML content for a meta refresh tag and extracts the redirect URL if present.
func extractMetaRefreshURL(htmlContent string) (string, bool) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		// Handle error if necessary
		return "", false
	}

	var findMetaRefresh func(*html.Node) (string, bool)
	findMetaRefresh = func(n *html.Node) (string, bool) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			httpEquivPresent := false
			contentValue := ""

			for _, a := range n.Attr {
				if strings.EqualFold(a.Key, "http-equiv") && strings.EqualFold(a.Val, "refresh") {
					httpEquivPresent = true
				} else if a.Key == "content" {
					contentValue = a.Val
				}
			}

			if httpEquivPresent && contentValue != "" {
				// Extract URL from content, expected format: "0; URL='http://example.com/'"
				parts := strings.Split(contentValue, "URL=")
				if len(parts) > 1 {
					url := strings.TrimSpace(parts[1])
					// Remove potential surrounding quotes
					url = strings.Trim(url, `"'`)
					return url, true
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			url, found := findMetaRefresh(c)
			if found {
				return url, true
			}
		}
		return "", false
	}

	return findMetaRefresh(doc)
}
