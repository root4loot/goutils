package domainutil

import (
	"testing"
)

func TestIsValidDomain(t *testing.T) {
	tests := []struct {
		domain string
		valid  bool
	}{
		{"example.com", true},
		{"sub.example.com", true},
		{"example", false},
		{"example..com", false},
		{"-example.com", false},
		{"example-.com", false},
		{"example.com-", false},
		{"*.example.com", true},
		{"*.com", false},
	}

	for _, test := range tests {
		if IsValidDomain(test.domain) != test.valid {
			t.Errorf("IsValidDomain(%s) = %v; want %v", test.domain, !test.valid, test.valid)
		}
	}
}

func TestGetRootDomain(t *testing.T) {
	tests := []struct {
		domain     string
		rootDomain string
	}{
		{"example.com", "example.com"},
		{"sub.example.com", "example.com"},
		{"sub.sub.example.com", "example.com"},
		{"example.co.uk", "co.uk"},
		{"sub.example.co.uk", "co.uk"},
	}

	for _, test := range tests {
		if root := GetRootDomain(test.domain); root != test.rootDomain {
			t.Errorf("GetRootDomain(%s) = %s; want %s", test.domain, root, test.rootDomain)
		}
	}
}
