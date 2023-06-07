package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
)

func main() {
	resolvers := []string{"8.8.8.8:53", "1.1.1.1:53"}
	client, err := ClientWithCustomResolvers(resolvers)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := client.Get("https://google.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	// Process the response
	// ...
}

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
