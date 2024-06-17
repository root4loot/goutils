package strutil

import (
	"strconv"
	"strings"
	"unicode"
)

// TrimWhitespace removes all whitespace from a string.
func TrimWhitespace(str string) string {
	fields := strings.Fields(str)
	return strings.Join(fields, "")
}

// IsPrintable checks if a string is printable
func IsPrintable(str string) bool {
	for _, r := range str {
		if !unicode.IsPrint(r) && !strconv.IsPrint(r) {
			return false
		}
	}
	return true
}

// IsBinaryString checks if a string is binary
func IsBinaryString(str string) bool {
	for _, b := range []byte(str) {
		if b < 32 || b > 126 {
			return true
		}
	}
	return false
}
