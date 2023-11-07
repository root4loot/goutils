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
