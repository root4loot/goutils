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

func TestEnsureScheme(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"example.com", "http://example.com"},
		{"http://example.com", "http://example.com"},
	}

	for _, test := range tests {
		result := EnsureScheme(test.input)
		if result != test.expected {
			t.Errorf("EnsureScheme(%s) = %s; want %s", test.input, result, test.expected)
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
