package hostutil

import (
	"strings"
	"testing"
)

func TestIsValidHostname(t *testing.T) {
	tests := []struct {
		hostname string
		valid    bool
	}{
		{"example.com", true},
		{"localhost", true},
		{"sub.domain.example.com", true},
		{"example", true},
		{"example123.com", true},
		{"123example.com", true},
		{"example-com", true},
		{"example.com-", false},
		{"-example.com", false},
		{"exa_mple.com", false},
		{"example..com", false},
		{"", false},
		{strings.Repeat("a", 256), false},
		{"example!.com", false},
		{"example .com", false},
		{".example.com", false},
		{"example.com.", false},
		{"ex%ample.com", false},
		{"example.com/", false},
		{"example..com", false},
		{"-example-.com", false},
		{"ex--ample.com", true},
		{"example.-com", false},
		{"example.com-", false},
		{"example-.com", false},
		{"exa*mple.com", false},
		{"example@com", false},
		{"example,com", false},
	}

	for _, tt := range tests {
		t.Run(tt.hostname, func(t *testing.T) {
			result := IsValidHostname(tt.hostname)
			if result != tt.valid {
				t.Errorf("IsValidHostname(%q) = %v; want %v", tt.hostname, result, tt.valid)
			}
		})
	}
}

func TestIsValidHost(t *testing.T) {
	tests := []struct {
		host  string
		valid bool
	}{
		{"example.com", true},
		{"localhost", true},
		{"sub.domain.example.com", true},
		{"exa_mple.com", false},
		{"example..com", false},
		{"example.com:8080", true},
		{"localhost:3000", true},
		{"127.0.0.1:8080", true},
		{"exa_mple.com:8080", false},
		{"example.com:99999", false},
		{"example.com:-1", false},
		{"127.0.0.1", true},
		{"192.168.1.1", true},
		{"::1", true},
		{"2001:db8::ff00:42:8329", true},
		{"256.256.256.256", false},
		{"1234.123.123.123", false},
		{"2001:db8:::ff00:42:8329", false},
		{"127.0.0.1:80", true},
		{"192.168.1.1:65535", true},
		{"::1:8080", true},
		{"256.256.256.256:80", false},
		{"127.0.0.1:99999", false},
		{"127.0.0.1:-80", false},
	}

	for _, tt := range tests {
		t.Run(tt.host, func(t *testing.T) {
			result := IsValidHost(tt.host)
			if result != tt.valid {
				t.Errorf("IsValidHost(%q) = %v; want %v", tt.host, result, tt.valid)
			}
		})
	}
}
