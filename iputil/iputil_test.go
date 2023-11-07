package iputil

import (
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
