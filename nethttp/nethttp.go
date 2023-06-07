package nethttp

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
)

// ClientWithCustomResolvers creates an HTTP client with custom resolvers list
func ClientWithCustomResolvers(resolvers []string) (*http.Client, error) {
	// Check if the resolvers contain the port :53, if not, add it
	for i, resolver := range resolvers {
		_, _, err := net.SplitHostPort(resolver)
		if err != nil {
			// Port is not specified, add :53
			resolvers[i] = resolver + ":53"
		}
	}

	dialer := &net.Dialer{
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				for _, resolver := range resolvers {
					d := net.Dialer{}
					conn, err := d.DialContext(ctx, "udp", resolver)
					if err != nil {
						continue // Try the next resolver
					}
					return conn, nil // Success, return the connection
				}
				return nil, fmt.Errorf("failed to resolve domain using custom resolvers")
			},
		},
	}

	transport := &http.Transport{
		DialContext:         dialer.DialContext,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives:   true,
		MaxIdleConnsPerHost: -1,
	}

	client := &http.Client{
		Transport: transport,
	}

	return client, nil
}

// ClientWithDefaultAndCustomResolvers creates an HTTP client with default and custom resolvers
func ClientWithDefaultAndCustomResolvers(resolvers []string) (*http.Client, error) {
	// Check if the resolvers contain the port :53, if not, add it
	for i, resolver := range resolvers {
		_, _, err := net.SplitHostPort(resolver)
		if err != nil {
			// Port is not specified, add :53
			resolvers[i] = resolver + ":53"
		}
	}

	defaultResolver := &net.Resolver{PreferGo: true}

	dialer := &net.Dialer{
		Resolver: defaultResolver,
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			conn, err := dialer.DialContext(ctx, network, address)
			if err != nil {
				for _, resolver := range resolvers {
					d := net.Dialer{}
					conn, err = d.DialContext(ctx, network, resolver)
					if err != nil {
						continue // Try the next resolver
					}
					return conn, nil // Success, return the connection
				}
				return nil, fmt.Errorf("failed to resolve domain using custom resolvers")
			}
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
