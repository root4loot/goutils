package iputil

import (
	"fmt"
	"net"
	"reflect"
	"testing"
)

func TestParseIPRange(t *testing.T) {
	tests := []struct {
		name    string
		ipRange string
		want    []net.IP
		wantErr bool
	}{
		{
			name:    "valid range with full IPs",
			ipRange: "10.0.0.1-10.0.0.5",
			want:    []net.IP{net.ParseIP("10.0.0.1"), net.ParseIP("10.0.0.2"), net.ParseIP("10.0.0.3"), net.ParseIP("10.0.0.4"), net.ParseIP("10.0.0.5")},
			wantErr: false,
		},
		{
			name:    "valid range with abbreviated end IP",
			ipRange: "10.0.0.1-5",
			want:    []net.IP{net.ParseIP("10.0.0.1"), net.ParseIP("10.0.0.2"), net.ParseIP("10.0.0.3"), net.ParseIP("10.0.0.4"), net.ParseIP("10.0.0.5")},
			wantErr: false,
		},
		{
			name:    "single IP range",
			ipRange: "10.0.0.1-10.0.0.1",
			want:    []net.IP{net.ParseIP("10.0.0.1")},
			wantErr: false,
		},
		{
			name:    "start IP greater than end IP",
			ipRange: "10.0.0.5-10.0.0.1",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid IP format",
			ipRange: "10.0.0.256-10.0.0.260",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "non-numeric abbreviated end IP",
			ipRange: "10.0.0.1-XYZ",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative abbreviated end IP",
			ipRange: "10.0.0.1--1",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "out of range abbreviated end IP",
			ipRange: "10.0.0.1-300",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "extra hyphens in range",
			ipRange: "10.0.0.1--10.0.0.5",
			want:    nil,
			wantErr: true,
		},
		// Add more test cases as needed.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIPRange(tt.ipRange)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseIPRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseIPRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestIPRangeToCIDR(t *testing.T) {
	t.Run("valid IP range", func(t *testing.T) {
		ipRange := "192.168.1.1-192.168.1.3"
		want := []string{"192.168.1.1/32", "192.168.1.2/32", "192.168.1.3/32"}
		got, err := IPRangeToCIDR(ipRange)
		if err != nil {
			t.Fatalf("IPRangeToCIDR() error = %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("IPRangeToCIDR() = %v, want %v", got, want)
		}
	})
}

func TestCIDRtoIPRange(t *testing.T) {
	t.Run("valid CIDR", func(t *testing.T) {
		cidr := "192.168.1.0/30"
		want := "192.168.1.1 - 192.168.1.2"
		got, err := CIDRtoIPRange(cidr)
		if err != nil {
			t.Fatalf("CIDRtoIPRange() error = %v", err)
		}
		if got != want {
			t.Errorf("CIDRtoIPRange() = %v, want %v", got, want)
		}
	})
}

func TestIsValidNetworkInput(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"192.168.1.1", true},
		{"192.168.1.0/24", true},
		{"192.168.1.1-192.168.1.5", true},
		{"invalid", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsValidNetworkInput(tt.input); got != tt.want {
				t.Errorf("IsValidNetworkInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIP(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"192.168.1.1", true},
		{"192.168.1.1:8080", true},
		{"invalid", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsIP(tt.input); got != tt.want {
				t.Errorf("IsIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsURLIP(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"http://192.168.1.1", true},
		{"http://example.com", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsURLIP(tt.input); got != tt.want {
				t.Errorf("IsURLIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsCIDR(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"192.168.1.0/24", true},
		{"192.168.1.1", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsCIDR(tt.input); got != tt.want {
				t.Errorf("IsCIDR() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIPRange(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"192.168.1.1-192.168.1.5", true},
		{"192.168.1.1", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsIPRange(tt.input); got != tt.want {
				t.Errorf("IsIPRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidIP(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"192.168.1.1", true},
		{"invalid", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsValidIP(tt.input); got != tt.want {
				t.Errorf("IsValidIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidIPRange(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"192.168.1.1-192.168.1.5", true},
		{"192.168.1.5-192.168.1.1", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsValidIPRange(tt.input); got != tt.want {
				t.Errorf("IsValidIPRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidCIDR(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"192.168.1.0/24", true},
		{"192.168.1.1", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsValidCIDR(tt.input); got != tt.want {
				t.Errorf("IsValidCIDR() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIPv4(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"192.168.1.1", true},
		{"::1", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsIPv4(tt.input); got != tt.want {
				t.Errorf("IsIPv4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIPv6(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"::1", true},
		{"192.168.1.1", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsIPv6(tt.input); got != tt.want {
				t.Errorf("IsIPv6() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPrivateIP(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"192.168.1.1", true},
		{"8.8.8.8", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsPrivateIP(tt.input); got != tt.want {
				t.Errorf("IsPrivateIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPublicIP(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"8.8.8.8", true},
		{"192.168.1.1", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsPublicIP(tt.input); got != tt.want {
				t.Errorf("IsPublicIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIPInCIDR(t *testing.T) {
	tests := []struct {
		ip   string
		cidr string
		want bool
	}{
		{"192.168.1.1", "192.168.1.0/24", true},
		{"8.8.8.8", "192.168.1.0/24", false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s in %s", tt.ip, tt.cidr), func(t *testing.T) {
			if got := IsIPInCIDR(tt.ip, tt.cidr); got != tt.want {
				t.Errorf("IsIPInCIDR() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIPInRange(t *testing.T) {
	tests := []struct {
		ip      string
		ipRange string
		want    bool
	}{
		{"192.168.1.1", "192.168.1.0-192.168.1.5", true},
		{"8.8.8.8", "192.168.1.0-192.168.1.5", false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s in %s", tt.ip, tt.ipRange), func(t *testing.T) {
			if got := IsIPInRange(tt.ip, tt.ipRange); got != tt.want {
				t.Errorf("IsIPInRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReverseDNSLookup(t *testing.T) {
	t.Run("valid IP", func(t *testing.T) {
		ip := "8.8.8.8"
		names, err := ReverseDNSLookup(ip)
		if err != nil {
			t.Fatalf("ReverseDNSLookup() error = %v", err)
		}
		if len(names) == 0 {
			t.Errorf("ReverseDNSLookup() returned no names")
		}
	})
}

func TestParseCIDR(t *testing.T) {
	t.Run("valid CIDR", func(t *testing.T) {
		cidr := "192.168.1.0/30"
		expectedIPs := []string{
			"192.168.1.0",
			"192.168.1.1",
			"192.168.1.2",
			"192.168.1.3",
		}
		ips, err := ParseCIDR(cidr)
		if err != nil {
			t.Fatalf("ParseCIDR() error = %v, wantErr false", err)
		}

		var ipsStrings []string
		for _, ip := range ips {
			ipsStrings = append(ipsStrings, ip.String())
		}

		if !reflect.DeepEqual(ipsStrings, expectedIPs) {
			t.Errorf("ParseCIDR() = %v, want %v", ipsStrings, expectedIPs)
		}
	})

	// ... other test cases ...
}

func TestIPsToCIDR(t *testing.T) {
	t.Run("convert IPs to /32 CIDR blocks", func(t *testing.T) {
		ips := []string{"192.168.1.0", "192.168.1.1", "192.168.1.2", "192.168.1.3"}
		got, err := IPsToCIDR(ips)
		if err != nil {
			t.Fatalf("IPsToCIDR() error = %v", err)
		}

		want := []string{"192.168.1.0/32", "192.168.1.1/32", "192.168.1.2/32", "192.168.1.3/32"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("IPsToCIDR() = %v, want %v", got, want)
		}
	})

}

func TestIPsToRange(t *testing.T) {
	tests := []struct {
		name    string
		IPs     []string
		want    []string
		wantErr bool
	}{
		{
			name: "unordered IPs",
			IPs:  []string{"192.168.1.3", "192.168.1.1", "192.168.1.5", "192.168.1.2"},
			want: []string{"192.168.1.1 - 192.168.1.3", "192.168.1.5 - 192.168.1.5"},
		},
		// other test cases...
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IPsToRange(tt.IPs)
			if (err != nil) != tt.wantErr {
				t.Errorf("IPsToRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IPsToRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
