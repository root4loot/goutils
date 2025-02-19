package fileutil

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// FileExists checks whether a file exists.
func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	return err == nil && !info.IsDir()
}

// EnsureDir ensures that the parent directory exists. If it does not, it creates it.
func EnsureDir(dirPath string) error {
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %q: %w", dirPath, err)
	}
	return nil
}

// CreateFileIfNotExists explicitly creates a file only if it does not already exist.
func CreateFileIfNotExists(filePath string) error {
	if FileExists(filePath) {
		return nil
	}

	dir := filepath.Dir(filePath)
	if err := EnsureDir(dir); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	return nil
}

// SerializeToFile writes JSON data to an existing file. It does NOT create the file or directories.
func SerializeToFile(filePath string, data interface{}) error {
	if !FileExists(filePath) {
		return fmt.Errorf("file %q does not exist", filePath)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(jsonData); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

// DeserializeFromFile reads JSON data from a file and unmarshals it.
func DeserializeFromFile(filePath string, data interface{}) error {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if err := json.Unmarshal(fileData, data); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

// ReadFile reads an existing file and returns a slice of strings.
func ReadFile(filePath string) ([]string, error) {
	if !FileExists(filePath) {
		return nil, fmt.Errorf("file %q does not exist", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// WriteSliceToFile writes a slice of strings to an existing file. It does NOT create the file.
func WriteSliceToFile(filePath string, lines []string) error {
	if !FileExists(filePath) {
		return fmt.Errorf("file %q does not exist", filePath)
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write line: %w", err)
		}
	}

	return writer.Flush()
}

// WriteStringToFile writes a string to an existing file. It does NOT create the file.
func WriteStringToFile(filePath string, content string) error {
	if !FileExists(filePath) {
		return fmt.Errorf("file %q does not exist", filePath)
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if _, err := writer.WriteString(content + "\n"); err != nil {
		return fmt.Errorf("failed to write content: %w", err)
	}

	return writer.Flush()
}

// WriteToFileAppend appends a line to an existing file. It does NOT create the file.
func WriteToFileAppend(filePath string, line string) error {
	if !FileExists(filePath) {
		return fmt.Errorf("file %q does not exist", filePath)
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for appending: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if _, err := writer.WriteString(line + "\n"); err != nil {
		return fmt.Errorf("failed to append line: %w", err)
	}

	return writer.Flush()
}
