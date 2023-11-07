package iputil

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// ParseIPRange parses the provided IP range and returns a slice of IP addresses.
func ParseIPRange(ipRange string) ([]net.IP, error) {
	parts := strings.Split(ipRange, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid IP range format: %s", ipRange)
	}

	startIP := net.ParseIP(strings.TrimSpace(parts[0]))
	if startIP == nil {
		return nil, fmt.Errorf("invalid start IP address: %s", parts[0])
	}

	var endIP net.IP
	if strings.Contains(parts[1], ".") {
		endIP = net.ParseIP(strings.TrimSpace(parts[1]))
		if endIP == nil {
			return nil, fmt.Errorf("invalid end IP address: %s", parts[1])
		}
	} else {
		inc, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid IP range format: %s", ipRange)
		}
		endIP = calculateEndIP(startIP, inc)
	}

	var ips []net.IP
	ip := startIP
	for ipLessThan(ip, endIP) {
		clone := make(net.IP, len(ip))
		copy(clone, ip)
		ips = append(ips, clone)

		inc(ip)
	}

	return ips, nil
}

// ParseCIDR parses the provided CIDR and returns a slice of IP addresses.
func ParseCIDR(cidr string) ([]net.IP, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []net.IP
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		clone := make(net.IP, len(ip))
		copy(clone, ip)
		ips = append(ips, clone)
	}

	return ips, nil
}

// IsValidNetworkInput checks if the provided string is a valid IP address, CIDR or IP range.
func IsValidNetworkInput(str string) bool {
	return IsIP(str) || IsCIDR(str) || IsIPRange(str)
}

// IsIPAddress checks if the provided string is an IP address with optional port.
func IsIP(str string) bool {
	ipPattern := `^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}(:\d+)?$`
	match, _ := regexp.MatchString(ipPattern, str)
	return match
}

// IsURLIP checks if the provided string is a URL with an IP address.
func IsURLIP(str string) (bool, error) {
	parsedURL, err := url.Parse(str)
	if err != nil {
		return false, err // Parse error
	}

	host := parsedURL.Hostname()
	if host == "" {
		return false, nil // Return false if the hostname is empty
	}

	// Check if the host is a valid IP address (IPv4 or IPv6)
	if net.ParseIP(host) != nil {
		return true, nil
	}

	// If the host is not a valid IP address, check if it has a port
	if strings.Contains(host, ":") {
		host, _, err = net.SplitHostPort(host)
		if err != nil {
			return false, err // Split error
		}

		// Check if the modified host is a valid IP address (IPv4 or IPv6)
		if net.ParseIP(host) != nil {
			return true, nil
		}
	}

	return false, nil
}

// IsCIDR checks if the provided string is a CIDR.
func IsCIDR(str string) bool {
	cidrPattern := `^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}/\d{1,2}$`
	match, _ := regexp.MatchString(cidrPattern, str)
	return match
}

// IsIPRange checks if the provided string is a valid IP range.
func IsIPRange(str string) bool {
	ipRangePattern := `^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}-\d{1,3}(\.\d{1,3}){0,3}$`
	match, _ := regexp.MatchString(ipRangePattern, str)
	return match
}

// IsValidIP checks if the provided IP address is valid.
func IsValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}

// IsValidIPRange checks if the provided IP range is valid.
func IsValidIPRange(ipRange string) bool {
	parts := strings.Split(ipRange, "-")
	if len(parts) != 2 {
		return false
	}

	startIP := net.ParseIP(strings.TrimSpace(parts[0]))
	if startIP == nil {
		return false
	}

	endIP := net.ParseIP(strings.TrimSpace(parts[1]))
	if endIP == nil {
		endIP = net.ParseIP(strings.TrimSpace(parts[0])) // Treat it as a single IP
		if endIP == nil {
			return false
		}
	}

	start := ipToUint32(startIP)
	end := ipToUint32(endIP)

	return start <= end
}

// IsValidCIDR checks if the provided CIDR is valid.
func IsValidCIDR(cidr string) bool {
	_, _, err := net.ParseCIDR(cidr)
	return err == nil
}

// IsIPv4 checks if the provided IP address is an IPv4 address.
func IsIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() != nil
}

// IsIPv6 checks if the provided IP address is an IPv6 address.
func IsIPv6(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() == nil
}

// IsPublicIP checks if the provided IP address is a public IP address.
func IsPublicIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// Check if the IP address is within the private address ranges
	privateBlocks := []*net.IPNet{
		{IP: net.ParseIP("10.0.0.0"), Mask: net.CIDRMask(8, 32)},
		{IP: net.ParseIP("172.16.0.0"), Mask: net.CIDRMask(12, 32)},
		{IP: net.ParseIP("192.168.0.0"), Mask: net.CIDRMask(16, 32)},
	}

	for _, block := range privateBlocks {
		if block.Contains(parsedIP) {
			return false
		}
	}

	return true
}

// IsIPInCIDR checks if the provided IP address is within the provided CIDR.
func IsIPInCIDR(ip string, cidr string) (bool, error) {
	ipAddr := net.ParseIP(ip)
	_, network, err := net.ParseCIDR(cidr)
	return err == nil && network.Contains(ipAddr), err
}

// IsIPInRange checks if the provided IP address is within the provided IP range.
func IsIPInRange(ip string, ipRange string) (bool, error) {
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return false, fmt.Errorf("invalid IP address: %s", ip)
	}

	ips, err := ParseIPRange(ipRange)
	if err != nil {
		return false, err
	}

	for _, addr := range ips {
		if addr.Equal(ipAddr) {
			return true, nil
		}
	}

	return false, err
}

// calculateEndIP calculates the end IP address based on the given increment value.
func calculateEndIP(startIP net.IP, increment int) net.IP {
	endIP := make(net.IP, len(startIP))
	copy(endIP, startIP)

	for i := len(endIP) - 1; i >= 0; i-- {
		sum := int(endIP[i]) + increment
		endIP[i] = byte(sum)
		increment = sum >> 8
	}

	return endIP
}

// ipLessThan checks if the first IP is less than the second IP.
func ipLessThan(ip1, ip2 net.IP) bool {
	return bytesCompare(ip1, ip2) < 0
}

// ipLessThanOrEqual checks if the first IP is less than or equal to the second IP.
func ipLessThanOrEqual(ip1, ip2 net.IP) bool {
	return bytesCompare(ip1, ip2) <= 0
}

// ipToUint32 converts an IP address to a uint32.
func ipToUint32(ip net.IP) uint32 {
	ip = ip.To4()
	return (uint32(ip[0]) << 24) | (uint32(ip[1]) << 16) | (uint32(ip[2]) << 8) | uint32(ip[3])
}

// uint32ToIP converts a uint32 to an IP address.
func uint32ToIP(ip uint32) net.IP {
	return net.IPv4(byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

// ipGreaterThan checks if the first IP is greater than the second IP.
func ipGreaterThan(ip1, ip2 net.IP) bool {
	return bytesCompare(ip1, ip2) > 0
}

// bytesCompare compares two byte slices lexicographically.
func bytesCompare(a, b []byte) int {
	lenA, lenB := len(a), len(b)
	for i := 0; i < lenA && i < lenB; i++ {
		if a[i] != b[i] {
			if a[i] < b[i] {
				return -1
			}
			return 1
		}
	}
	if lenA == lenB {
		return 0
	} else if lenA < lenB {
		return -1
	}
	return 1
}

// ipEqual checks if two IP addresses are equal.
func ipEqual(ip1, ip2 net.IP) bool {
	return ip1.Equal(ip2)
}

// inc increments the IP address.
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
