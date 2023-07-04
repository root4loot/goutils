package domainutil

import "regexp"

// IsDomainName checks if a string is a valid domain name.
func IsDomainName(str string) bool {
	// Regular expression pattern to match domain names with optional wildcards
	regex := regexp.MustCompile(`^(?i)(\*\.){0,1}(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.){1,}(?:[a-z]{2,})$`)

	return regex.MatchString(str)
}
