package fileutil

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// SerializeToFile serializes data to a JSON file.
func SerializeToFile(filePath string, data interface{}) error {
	// Convert struct to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Create or open file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write JSON data to file
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// DeserializeFromFile deserializes data from a JSON file.
func DeserializeFromFile(filePath string, data interface{}) error {
	// Read file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Parse JSON data into struct
	err = json.Unmarshal(fileData, data)
	if err != nil {
		return err
	}

	return nil
}

// ReadFile reads a file and returns a slice of strings representing lines.
func ReadFile(fp string) ([]string, error) {
	readFile, err := os.Open(fp)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	return fileLines, nil
}

// WriteToFile takes a filepath and a slice of strings, then writes the strings to the file.
func WriteToFile(filePath string, lines []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

// WriteToFileAppend appends a line to a file.
func WriteToFileAppend(filePath string, line string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(line + "\n")
	if err != nil {
		return err
	}

	return writer.Flush()
}
