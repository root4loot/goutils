package fileutil

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func setupTestDir(t *testing.T) string {
	t.Helper()
	testDir := "test_data"
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	return testDir
}

func cleanupTestDir(t *testing.T, testDir string) {
	t.Helper()
	if err := os.RemoveAll(testDir); err != nil {
		t.Fatalf("Failed to clean up test directory: %v", err)
	}
}

func TestFileExists(t *testing.T) {
	testDir := setupTestDir(t)
	defer cleanupTestDir(t, testDir)

	testFile := filepath.Join(testDir, "testfile.txt")
	os.WriteFile(testFile, []byte("hello"), 0644)

	if !FileExists(testFile) {
		t.Errorf("FileExists() returned false for an existing file")
	}

	nonExistentFile := filepath.Join(testDir, "nonexistent.txt")
	if FileExists(nonExistentFile) {
		t.Errorf("FileExists() returned true for a non-existing file")
	}
}

func TestEnsureDir(t *testing.T) {
	testDir := setupTestDir(t)
	defer cleanupTestDir(t, testDir)

	newDir := filepath.Join(testDir, "subdir")
	err := EnsureDir(newDir)
	if err != nil {
		t.Fatalf("EnsureDir() failed: %v", err)
	}

	if _, err := os.Stat(newDir); os.IsNotExist(err) {
		t.Errorf("EnsureDir() did not create the directory")
	}
}

func TestCreateFileIfNotExists(t *testing.T) {
	testDir := setupTestDir(t)
	defer cleanupTestDir(t, testDir)

	testFile := filepath.Join(testDir, "newfile.txt")
	err := CreateFileIfNotExists(testFile)
	if err != nil {
		t.Fatalf("CreateFileIfNotExists() failed: %v", err)
	}

	if !FileExists(testFile) {
		t.Errorf("CreateFileIfNotExists() did not create the file")
	}

	err = CreateFileIfNotExists(testFile)
	if err != nil {
		t.Fatalf("CreateFileIfNotExists() failed on existing file: %v", err)
	}
}

func TestWriteStringToFile(t *testing.T) {
	testDir := setupTestDir(t)
	defer cleanupTestDir(t, testDir)

	testFile := filepath.Join(testDir, "write_test.txt")
	os.WriteFile(testFile, []byte(""), 0644)

	err := WriteStringToFile(testFile, "Hello, Go!")
	if err != nil {
		t.Fatalf("WriteStringToFile() failed: %v", err)
	}

	content, _ := os.ReadFile(testFile)
	if string(content) != "Hello, Go!\n" {
		t.Errorf("WriteStringToFile() content mismatch: got %q", string(content))
	}
}

func TestWriteSliceToFile(t *testing.T) {
	testDir := setupTestDir(t)
	defer cleanupTestDir(t, testDir)

	testFile := filepath.Join(testDir, "slice_test.txt")
	os.WriteFile(testFile, []byte(""), 0644)

	lines := []string{"Line 1", "Line 2", "Line 3"}
	err := WriteSliceToFile(testFile, lines)
	if err != nil {
		t.Fatalf("WriteSliceToFile() failed: %v", err)
	}

	content, _ := os.ReadFile(testFile)
	expected := "Line 1\nLine 2\nLine 3\n"
	if string(content) != expected {
		t.Errorf("WriteSliceToFile() content mismatch: got %q, want %q", string(content), expected)
	}
}

func TestWriteToFileAppend(t *testing.T) {
	testDir := setupTestDir(t)
	defer cleanupTestDir(t, testDir)

	testFile := filepath.Join(testDir, "append_test.txt")
	os.WriteFile(testFile, []byte("Initial\n"), 0644)

	err := WriteToFileAppend(testFile, "Appended line")
	if err != nil {
		t.Fatalf("WriteToFileAppend() failed: %v", err)
	}

	content, _ := os.ReadFile(testFile)
	expected := "Initial\nAppended line\n"
	if string(content) != expected {
		t.Errorf("WriteToFileAppend() content mismatch: got %q, want %q", string(content), expected)
	}
}

func TestReadFile(t *testing.T) {
	testDir := setupTestDir(t)
	defer cleanupTestDir(t, testDir)

	testFile := filepath.Join(testDir, "read_test.txt")
	content := "Line A\nLine B\nLine C\n"
	os.WriteFile(testFile, []byte(content), 0644)

	lines, err := ReadFile(testFile)
	if err != nil {
		t.Fatalf("ReadFile() failed: %v", err)
	}

	expected := []string{"Line A", "Line B", "Line C"}
	for i, line := range expected {
		if lines[i] != line {
			t.Errorf("ReadFile() mismatch at line %d: got %q, want %q", i, lines[i], line)
		}
	}
}

func TestSerializeDeserialize(t *testing.T) {
	testDir := setupTestDir(t)
	defer cleanupTestDir(t, testDir)

	testFile := filepath.Join(testDir, "json_test.json")
	os.WriteFile(testFile, []byte("{}"), 0644)

	type Sample struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	writeData := Sample{Name: "Alice", Age: 30}
	err := SerializeToFile(testFile, writeData)
	if err != nil {
		t.Fatalf("SerializeToFile() failed: %v", err)
	}

	var readData Sample
	err = DeserializeFromFile(testFile, &readData)
	if err != nil {
		t.Fatalf("DeserializeFromFile() failed: %v", err)
	}

	if readData != writeData {
		t.Errorf("DeserializeFromFile() returned incorrect data: got %+v, want %+v", readData, writeData)
	}
}

func TestSaveJSONLines(t *testing.T) {
	// Temporary test file
	testFile := "test_data.jsonl"
	defer os.Remove(testFile)

	data1 := map[string]interface{}{
		"id":   "123",
		"name": "test1",
	}
	data2 := map[string]interface{}{
		"id":   "456",
		"name": "test2",
	}

	err := SaveJSONLines(testFile, data1)
	if err != nil {
		t.Fatalf("Failed to save first JSON line: %v", err)
	}

	err = SaveJSONLines(testFile, data2)
	if err != nil {
		t.Fatalf("Failed to save second JSON line: %v", err)
	}

	results, err := LoadJSONLines(testFile)
	if err != nil {
		t.Fatalf("Failed to load JSON lines: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(results))
	}

	if results[0]["id"] != "123" || results[0]["name"] != "test1" {
		t.Errorf("Unexpected first entry: %v", results[0])
	}

	if results[1]["id"] != "456" || results[1]["name"] != "test2" {
		t.Errorf("Unexpected second entry: %v", results[1])
	}
}

func TestLoadJSONLines_FileNotExist(t *testing.T) {
	_, err := LoadJSONLines("non_existent.jsonl")
	if err == nil {
		t.Fatal("Expected an error for a non-existent file, but got nil")
	}
}

func TestWriteJSONToFile(t *testing.T) {
	testDir := "testdata"
	testFile := filepath.Join(testDir, "test.json")
	defer os.RemoveAll(testDir) // Cleanup after test

	testData := map[string]string{"key": "value"}

	err := WriteJSONToFile(testFile, testData)
	if err != nil {
		t.Fatalf("WriteJSONToFile failed: %v", err)
	}

	// Verify file exists
	if !FileExists(testFile) {
		t.Fatalf("Expected file %q to exist", testFile)
	}

	// Verify file contents
	var readData map[string]string
	err = DeserializeFromFile(testFile, &readData)
	if err != nil {
		t.Fatalf("Failed to read JSON from file: %v", err)
	}

	// Compare written and read data
	expected, _ := json.Marshal(testData)
	actual, _ := json.Marshal(readData)
	if string(expected) != string(actual) {
		t.Errorf("Expected %q but got %q", expected, actual)
	}
}
