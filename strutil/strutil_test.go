package strutil

import (
	"testing"
)

func TestTrimWhitespace(t *testing.T) {
	input := "  Hello   World  "
	expected := "HelloWorld"
	result := TrimWhitespace(input)

	if result != expected {
		t.Errorf("TrimWhitespace failed, expected %v, got %v", expected, result)
	}
}

func TestIsPrintable(t *testing.T) {
	printable := "Hello, World!"
	nonPrintable := "Hello\x00World"

	if !IsPrintable(printable) {
		t.Error("IsPrintable failed, expected true, got false")
	}

	if IsPrintable(nonPrintable) {
		t.Error("IsPrintable failed, expected false, got true")
	}
}

func TestIsBinaryString(t *testing.T) {
	binary := "Hello\x00World"
	nonBinary := "Hello, World!"

	if !IsBinaryString(binary) {
		t.Error("IsBinaryString failed, expected true, got false")
	}

	if IsBinaryString(nonBinary) {
		t.Error("IsBinaryString failed, expected false, got true")
	}
}
