package sysutil

import (
	"os/user"
	"strings"
)

// IsRoot checks if the current user is root.
func IsRoot() (bool, error) {
	currentUser, err := user.Current()
	if err != nil {
		return false, err
	}
	return currentUser.Username == "root", nil
}

// IsTextContentType checks if a content type is text-based
func IsTextContentType(contentType string) bool {
	var nonTextContentTypes = []string{
		"application/octet-stream",
		"application/pdf",
		"application/zip",
		"application/x-gzip",
		"application/vnd.ms-excel",
		"application/vnd.ms-powerpoint",
		"application/vnd.ms-word",
	}

	for _, nonTextType := range nonTextContentTypes {
		if nonTextType == contentType {
			return false
		}
	}

	// Check for generic types
	part := strings.Split(contentType, "/")[0]
	if part == "image" || part == "audio" || part == "video" {
		return false
	}

	return true
}
