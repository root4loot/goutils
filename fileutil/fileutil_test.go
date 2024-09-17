package fileutil

import (
	"os"
	"testing"
)

func TestSerializeToFileAndDeserializeFromFile(t *testing.T) {
	type TestData struct {
		Name string
		Age  int
	}

	filePath := "test.json"
	data := TestData{Name: "John Doe", Age: 30}

	err := SerializeToFile(filePath, data)
	if err != nil {
		t.Fatalf("SerializeToFile failed: %v", err)
	}
	defer os.Remove(filePath)

	var result TestData
	err = DeserializeFromFile(filePath, &result)
	if err != nil {
		t.Fatalf("DeserializeFromFile failed: %v", err)
	}

	if result != data {
		t.Errorf("Expected %v, got %v", data, result)
	}
}

func TestReadFileAndWriteSliceToFile(t *testing.T) {
	filePath := "test.txt"
	lines := []string{"line1", "line2", "line3"}

	err := WriteSliceToFile(filePath, lines)
	if err != nil {
		t.Fatalf("WriteSliceToFile failed: %v", err)
	}
	defer os.Remove(filePath)

	readLines, err := ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	for i, line := range lines {
		if readLines[i] != line {
			t.Errorf("Expected %s, got %s", line, readLines[i])
		}
	}
}

func TestWriteStringToFile(t *testing.T) {
	filePath := "test_string.txt"
	content := "Hello, World!"

	err := WriteStringToFile(filePath, content)
	if err != nil {
		t.Fatalf("WriteStringToFile failed: %v", err)
	}
	defer os.Remove(filePath)

	readLines, err := ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if readLines[0] != content {
		t.Errorf("Expected %s, got %s", content, readLines[0])
	}
}

func TestWriteToFileAppend(t *testing.T) {
	filePath := "test_append.txt"
	initialContent := "First line"
	appendContent := "Second line"

	err := WriteStringToFile(filePath, initialContent)
	if err != nil {
		t.Fatalf("WriteStringToFile failed: %v", err)
	}
	defer os.Remove(filePath)

	err = WriteToFileAppend(filePath, appendContent)
	if err != nil {
		t.Fatalf("WriteToFileAppend failed: %v", err)
	}

	readLines, err := ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if readLines[0] != initialContent || readLines[1] != appendContent {
		t.Errorf("Expected [%s, %s], got [%s, %s]", initialContent, appendContent, readLines[0], readLines[1])
	}
}
