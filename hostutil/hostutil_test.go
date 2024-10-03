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
