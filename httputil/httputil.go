package httputil

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"sync"
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
