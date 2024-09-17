package netutil

import (
	"net"
	"strconv"
	"testing"
	"time"
)

func TestCanDial(t *testing.T) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to set up mock server: %v", err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port
	if !CanDial("localhost", strconv.Itoa(port)) {
		t.Errorf("Expected CanDial to return true for open port, got false")
	}

	// Test with a known closed port
	if CanDial("localhost", "9999") {
		t.Errorf("Expected CanDial to return false for closed port, got true")
	}
}

func TestCanDialWithTimeout(t *testing.T) {
	// Test with a known open port
	ln, err := net.Listen("tcp", ":0") // Listen on a random open port
	if err != nil {
		t.Fatalf("Failed to set up mock server: %v", err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port
	if !CanDialWithTimeout("localhost", strconv.Itoa(port), 1*time.Second) {
		t.Errorf("Expected CanDialWithTimeout to return true for open port, got false")
	}

	// Test with a known closed port
	if CanDialWithTimeout("localhost", "9999", 1*time.Second) {
		t.Errorf("Expected CanDialWithTimeout to return false for closed port, got true")
	}
}
