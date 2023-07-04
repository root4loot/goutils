package netdomain

import "regexp"

func IsWildcardHostname(hostname string) bool {
	hostnamePattern := `^(\*\.)?[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)*$`
	match, _ := regexp.MatchString(hostnamePattern, hostname)
	return match
}
