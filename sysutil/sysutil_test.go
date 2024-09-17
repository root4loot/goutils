package sysutil

import (
	"testing"
)

func TestIsRoot(t *testing.T) {
	_, err := IsRoot()
	if err != nil {
		t.Errorf("IsRoot failed, expected no error, got %v", err)
	}
}

func TestIsTextContentType(t *testing.T) {
	tests := []struct {
		contentType string
		expected    bool
	}{
		{"text/plain", true},
		{"application/json", true},
		{"application/octet-stream", false},
		{"image/png", false},
		{"audio/mpeg", false},
		{"video/mp4", false},
	}

	for _, test := range tests {
		result := IsTextContentType(test.contentType)
		if result != test.expected {
			t.Errorf("IsTextContentType(%v) failed, expected %v, got %v", test.contentType, test.expected, result)
		}
	}
}
