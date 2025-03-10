package iputil

import (
	"bytes"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/yl2chen/cidranger"
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

	endIP := net.ParseIP(strings.TrimSpace(parts[1]))
	if endIP == nil {
		increment, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid IP range format: %s", ipRange)
		}
		endIP = make(net.IP, len(startIP))
		copy(endIP, startIP)
		lastOctet := int(startIP[len(startIP)-1]) + increment
		if lastOctet >= 256 {
			return nil, fmt.Errorf("invalid end IP value: results in an octet greater than 255")
		}
		endIP[len(endIP)-1] = byte(lastOctet - 1)
	} else {
		if bytes.Compare(startIP, endIP) > 0 {
			return nil, fmt.Errorf("start IP address %s is greater than end IP address %s", startIP, endIP)
		}
	}

	var ips []net.IP
	ip := make(net.IP, len(startIP))
	copy(ip, startIP)

	for {
		ips = append(ips, ip)

		if ip.Equal(endIP) {
			break
		}

		nextIP := make(net.IP, len(ip))
		copy(nextIP, ip)
		if !inc(nextIP) {
			return nil, fmt.Errorf("IP address overflowed while incrementing")
		}
		ip = nextIP
	}

	return ips, nil
}

// ParseCIDR parses the provided CIDR block and returns a slice of IP addresses.
func ParseCIDR(cidr string) ([]net.IP, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []net.IP
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); {
		clone := make(net.IP, len(ip))
		copy(clone, ip)
		ips = append(ips, clone)

		if !inc(ip) {
			break
		}
	}

	return ips, nil
}

// IPsToCIDR converts a slice of IPs to a slice of CIDR blocks.
func IPsToCIDR(IPs []string) (cidrs []string, err error) {
	ranger := cidranger.NewPCTrieRanger()

	for _, ipStr := range IPs {
		ip, _, err := net.ParseCIDR(ipStr + "/32")
		if err != nil {
			fmt.Printf("Skipping invalid IP: %s, error: %v\n", ipStr, err)
			continue
		}
		_, network, _ := net.ParseCIDR(ip.String() + "/32")
		err = ranger.Insert(cidranger.NewBasicRangerEntry(*network))
		if err != nil {
			return nil, fmt.Errorf("error inserting IP to ranger: %s, error: %v", ipStr, err)
		}
	}

	_, allIPv4Net, err := net.ParseCIDR("0.0.0.0/0")
	if err != nil {
		return nil, fmt.Errorf("error parsing all-encompassing IPv4 CIDR: %v", err)
	}
	entries, err := ranger.CoveredNetworks(*allIPv4Net)
	if err != nil {
		return nil, fmt.Errorf("error getting covered networks: %v", err)
	}
	for _, e := range entries {
		ones, _ := e.Network().Mask.Size() // This gives you the prefix length
		cidr := fmt.Sprintf("%s/%d", e.Network().IP.String(), ones)
		cidrs = append(cidrs, cidr)
	}

	return
}

// IPsToRange takes a slice of IP strings and returns a slice of IP ranges.
func IPsToRange(IPs []string) ([]string, error) {
	if len(IPs) == 0 {
		return nil, fmt.Errorf("no IPs provided")
	}

	sort.Slice(IPs, func(i, j int) bool {
		return bytes.Compare(net.ParseIP(IPs[i]), net.ParseIP(IPs[j])) < 0
	})

	startIP := IPs[0]
	endIP := IPs[0]

	var ipRanges []string
	for i := 1; i < len(IPs); i++ {
		currentIP := net.ParseIP(IPs[i])
		if !isConsecutive(net.ParseIP(endIP), currentIP) {
			ipRanges = append(ipRanges, startIP+" - "+endIP)
			startIP = IPs[i]
		}
		endIP = IPs[i]
	}

	ipRanges = append(ipRanges, startIP+" - "+endIP)

	return ipRanges, nil
}

// IPRangeToCIDR converts an IP range to a slice of CIDR blocks.
func IPRangeToCIDR(ipRange string) ([]string, error) {
	ips, err := ParseIPRange(ipRange)
	if err != nil {
		return nil, err
	}

	var cidrs []string
	for _, ip := range ips {
		cidrs = append(cidrs, ip.String()+"/32")
	}
	return cidrs, nil
}

// CIDRtoIPRange converts a CIDR block to an IP range.
func CIDRtoIPRange(cidr string) (string, error) {
	ips, err := ParseCIDR(cidr)
	if err != nil {
		return "", err
	}
	if len(ips) == 0 {
		return "", fmt.Errorf("no IPs in the CIDR")
	}

	startIP := ips[1]
	endIP := ips[len(ips)-2]

	return fmt.Sprintf("%s - %s", startIP.String(), endIP.String()), nil
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
func IsURLIP(str string) bool {
	parsedURL, err := url.Parse(str)
	if err != nil {
		return false // Parse error
	}

	host := parsedURL.Hostname()
	if host == "" {
		return false
	}

	if net.ParseIP(host) != nil {
		return true
	}

	if strings.Contains(host, ":") {
		host, _, err = net.SplitHostPort(host)
		if err != nil {
			return false
		}
		if net.ParseIP(host) != nil {
			return true
		}
	}

	return false
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
	if parsedIP == nil {
		return false
	}

	if strings.Count(ip, ".") == 3 {
		segments := strings.Split(ip, ".")
		for _, segment := range segments {
			num, err := strconv.Atoi(segment)
			if err != nil || num < 0 || num > 255 {
				return false
			}
		}
	}

	return true
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

// IsPrivateIP checks if the given IP is a private address.
func IsPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	privateBlocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	for _, block := range privateBlocks {
		_, cidr, err := net.ParseCIDR(block)
		if err != nil {
			continue
		}
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}

// IsPublicIP checks if the provided IP address is a public IP address.
func IsPublicIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	specialBlocks := []*net.IPNet{
		{IP: net.ParseIP("10.0.0.0"), Mask: net.CIDRMask(8, 32)},
		{IP: net.ParseIP("172.16.0.0"), Mask: net.CIDRMask(12, 32)},
		{IP: net.ParseIP("192.168.0.0"), Mask: net.CIDRMask(16, 32)},
		{IP: net.ParseIP("127.0.0.0"), Mask: net.CIDRMask(8, 32)},    // Loopback
		{IP: net.ParseIP("169.254.0.0"), Mask: net.CIDRMask(16, 32)}, // Link-local
		{IP: net.ParseIP("224.0.0.0"), Mask: net.CIDRMask(4, 32)},    // Multicast
		{IP: net.ParseIP("240.0.0.0"), Mask: net.CIDRMask(4, 32)},    // Future use/reserved
	}

	for _, block := range specialBlocks {
		if block.Contains(parsedIP) {
			return false
		}
	}

	return true
}

// IsIPInCIDR checks if the provided IP address is within the provided CIDR.
func IsIPInCIDR(ip string, cidr string) bool {
	ipAddr := net.ParseIP(ip)
	_, network, _ := net.ParseCIDR(cidr)
	return network.Contains(ipAddr)
}

// IsIPInRange checks if the provided IP address is within the provided IP range.
func IsIPInRange(ip string, ipRange string) bool {
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return false
	}

	ips, err := ParseIPRange(ipRange)
	if err != nil {
		return false
	}

	for _, addr := range ips {
		if addr.Equal(ipAddr) {
			return true
		}
	}

	return false
}

// GetPTRs returns the PTR record for the given IP address.
func GetPTRs(ip string) ([]string, error) {
	names, err := net.LookupAddr(ip)
	if err != nil {
		return nil, err
	}

	var validNames []string
	for _, name := range names {
		trimmedName := strings.TrimSuffix(name, ".")

		ips, err := net.LookupIP(trimmedName)
		if err == nil {
			for _, resolvedIP := range ips {
				if resolvedIP.String() == ip {
					validNames = append(validNames, trimmedName)
					break
				}
			}
		}
	}

	return validNames, nil
}

// calculateEndIP calculates the end IP address based on the given increment value.
func calculateEndIP(startIP net.IP, inc int) net.IP {
	endIP := make(net.IP, len(startIP))
	copy(endIP, startIP)

	lastOctet := int(endIP[len(endIP)-1]) + inc
	if lastOctet > 255 {
		return nil // This indicates an error condition.
	}

	endIP[len(endIP)-1] = byte(lastOctet)
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

// inc increments the IP address and returns false if there is an overflow.
func inc(ip net.IP) bool {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] != 0 {
			return true
		}
	}
	return false
}

// isConsecutive checks if two IP addresses are consecutive.
func isConsecutive(ip1, ip2 net.IP) bool {
	ip1 = ip1.To4()
	ip2 = ip2.To4()

	for i := len(ip1) - 1; i >= 0; i-- {
		ip1[i]++
		if ip1[i] != 0 {
			break
		}
	}

	return ip1.Equal(ip2)
}

// isSameSubnet checks if two IP addresses are in the same private subnet.
func isSameSubnet(ip1, ip2 net.IP) bool {
	return IsPrivateIP(ip1.String()) && IsPrivateIP(ip2.String())
}
