package log

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/term"
)

const (
	reset     = "\033[0m"
	lightGrey = "\033[90m"
	red       = "\033[91m"
	yellow    = "\033[93m"
	blue      = "\033[94m"
	white     = "\033[37m"
)

type CustomFormatter struct {
	packageName string
}

// Format formats the log output
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Check for a custom tag and use it if present
	levelText := strings.ToUpper(entry.Level.String())
	if label, ok := entry.Data["label"]; ok {
		levelText = label.(string) // Casting to string, ensure label is always a string
	}

	// Get color from entry fields, fall back to default if not set
	color, ok := entry.Data["tagColor"].(string)
	if !ok {
		color = white
	}

	logOutput := fmt.Sprintf("%s[%s]%s%s (%s)%s %s", lightGrey, f.packageName, reset, color, levelText, reset, entry.Message)

	// Prepare fields output
	fieldsOutput := ""
	if len(entry.Data) > 0 {
		for key, value := range entry.Data {
			if _, ok := entry.Data["tag"]; ok {
				fieldsOutput += fmt.Sprintf(" %s=%v", key, value) // Adjust this to colorize key/values
			}
		}
	}

	// Get terminal width
	width, _, _ := term.GetSize(0)

	// Calculate the spacing needed to push fields to the right
	paddingLength := width - len(logOutput) - len(fieldsOutput) - 1 // -1 for newline character
	if paddingLength > 0 {
		padding := strings.Repeat(" ", paddingLength)
		logOutput += padding
	}

	return []byte(logOutput + fieldsOutput + "\n"), nil
}
