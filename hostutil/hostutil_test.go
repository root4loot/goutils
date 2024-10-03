package hostutil

import (
	"testing"
)

func TestNormalizeHost(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"example.com", "example.com", false},
		{"Example.COM", "example.com", false},
		{"example.com:8080", "example.com:8080", false},
		{"example.com:80", "example.com", false},
		{"example.com:443", "example.com", false},
		{"subdomain.example.com:443", "subdomain.example.com", false},
		{"subdomain.example.com:80", "subdomain.example.com", false},
		{"  example.com  ", "example.com", false},
		{"invalid_host:port", "", true},
		{"subdomain:invalidport", "", true},
		{"http://example.com", "", true},
		{"https://example.com", "", true},
		{"ftp://example.com", "", true},
		{"example.com/path", "", true},
		{"http://example.com:443", "", true},
		{"example", "", true},
		{"localhost", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := NormalizeHost(tt.input)
			if (err != nil) != tt.hasError {
				t.Errorf("expected error status %v, got %v (error: %v)", tt.hasError, (err != nil), err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
