package strutil

import "strings"

// TrimWhitespace removes all whitespace from a string.
func TrimWhitespace(str string) string {
	fields := strings.Fields(str)
	return strings.Join(fields, "")
}
