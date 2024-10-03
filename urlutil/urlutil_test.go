package urlutil

import (
	"testing"
	"time"
)

func TestIsURL(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"http://example.com", true},
		{"https://example.com", true},
		{"ftp://example.com", true},
		{"example.com", false},
	}

	for _, test := range tests {
		result := IsURL(test.input)
		if result != test.expected {
			t.Errorf("IsURL(%s) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestCanReachURL(t *testing.T) {
	// This test assumes that "http://example.com" is reachable.
	err := CanReachURL("http://example.com")
	if err != nil {
		t.Errorf("CanReachURL(http://example.com) returned error: %v", err)
	}
}

func TestCanReachURLWithTimeout(t *testing.T) {
	// This test assumes that "http://example.com" is reachable.
	err := CanReachURLWithTimeout("http://example.com", 5*time.Second)
	if err != nil {
		t.Errorf("CanReachURLWithTimeout(http://example.com) returned error: %v", err)
	}
}

func TestEnsurePortIsSet(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"http://example.com", "http://example.com:80"},
		{"https://example.com", "https://example.com:443"},
		{"http://example.com:8080", "http://example.com:8080"},
	}

	for _, test := range tests {
		result := EnsurePortIsSet(test.input)
		if result != test.expected {
			t.Errorf("EnsurePortIsSet(%s) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestHasScheme(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"http://example.com", true},
		{"example.com", false},
	}

	for _, test := range tests {
		result := HasScheme(test.input)
		if result != test.expected {
			t.Errorf("HasScheme(%s) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestHasFileExtension(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"http://example.com/file.txt", true},
		{"http://example.com/file", false},
	}

	for _, test := range tests {
		result := HasFileExtension(test.input)
		if result != test.expected {
			t.Errorf("HasFileExtension(%s) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestHasParam(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"http://example.com?param=value", true},
		{"http://example.com", false},
	}

	for _, test := range tests {
		result := HasParam(test.input)
		if result != test.expected {
			t.Errorf("HasParam(%s) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestGetExt(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"http://example.com/file.txt", ".txt"},
		{"http://example.com/file", ""},
	}

	for _, test := range tests {
		result := GetExt(test.input)
		if result != test.expected {
			t.Errorf("GetExt(%s) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestEnsureTrailingSlash(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"http://example.com/path", "http://example.com/path/"},
		{"http://example.com/path/", "http://example.com/path/"},
	}

	for _, test := range tests {
		result := EnsureTrailingSlash(test.input)
		if result != test.expected {
			t.Errorf("EnsureTrailingSlash(%s) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestRemoveTrailingSlash(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"http://example.com/path/", "http://example.com/path"},
		{"http://example.com/path", "http://example.com/path"},
	}

	for _, test := range tests {
		result := RemoveTrailingSlash(test.input)
		if result != test.expected {
			t.Errorf("RemoveTrailingSlash(%s) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestNormalizeSlashes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"http://example.com//path", "http://example.com/path"},
		{"https://example.com//path", "https://example.com/path"},
		{"example.com//path", "example.com/path"},
	}

	for _, test := range tests {
		result := NormalizeSlashes(test.input)
		if result != test.expected {
			t.Errorf("NormalizeSlashes(%s) = %s; want %s", test.input, result, test.expected)
		}
	}
}

func TestGetOrigin(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"http://example.com/path", "http://example.com"},
		{"https://example.com/path", "https://example.com"},
	}

	for _, test := range tests {
		result, err := GetOrigin(test.input)
		if err != nil || result != test.expected {
			t.Errorf("GetOrigin(%s) = %s, %v; want %s, nil", test.input, result, err, test.expected)
		}
	}
}

func TestIsMediaExt(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{".jpg", true},
		{".txt", false},
	}

	for _, test := range tests {
		result := IsMediaExt(test.input)
		if result != test.expected {
			t.Errorf("IsMediaExt(%s) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestRemoveDefaultPort(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"http://example.com:80", "http://example.com", false},
		{"https://example.com:443", "https://example.com", false},
		{"http://example.com:8080", "http://example.com:8080", false},
		{"https://example.com:8443", "https://example.com:8443", false},
		{"ftp://example.com:21", "ftp://example.com", false},
		{"ftp://example.com:2121", "ftp://example.com:2121", false},
		{"http://example.com", "http://example.com", false},
		{"https://example.com", "https://example.com", false},
		{"invalid_url", "", true},
		{"example.com", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := RemoveDefaultPort(tt.input)
			if (err != nil) != tt.hasError {
				t.Errorf("expected error status %v, got %v (error: %v)", tt.hasError, (err != nil), err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEnsureHTTP(t *testing.T) {
	tests := []struct {
		rawURL   string
		expected string
	}{
		{"example.com", "http://example.com"},
		{"http://example.com", "http://example.com"},
		{"https://example.com", "http://example.com"},
		{"example.com/path", "http://example.com/path"},
		{"localhost", "http://localhost"},
		{"http://localhost", "http://localhost"},
		{"https://localhost", "http://localhost"},
		{"192.168.0.1", "http://192.168.0.1"},
		{"https://192.168.0.1", "http://192.168.0.1"},
		{"ftp://example.com", "http://example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.rawURL, func(t *testing.T) {
			result := EnsureHTTP(tt.rawURL)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEnsureHTTPS(t *testing.T) {
	tests := []struct {
		rawURL   string
		expected string
	}{
		{"example.com", "https://example.com"},
		{"http://example.com", "https://example.com"},
		{"https://example.com", "https://example.com"},
		{"example.com/path", "https://example.com/path"},
		{"localhost", "https://localhost"},
		{"http://localhost", "https://localhost"},
		{"192.168.0.1", "https://192.168.0.1"},
		{"http://192.168.0.1", "https://192.168.0.1"},
		{"ftp://example.com", "https://example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.rawURL, func(t *testing.T) {
			result := EnsureHTTPS(tt.rawURL)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
