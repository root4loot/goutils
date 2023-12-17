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
	green     = "\033[92m"
	purple    = "\033[95m"
	cyan      = "\033[96m"
	orange    = "\033[38;5;214m" // Adding a different orange color
)

type CustomFormatter struct {
	packageName string
}

// Format formats the log output
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var logOutput string

	// Check for a custom tag and use it if present
	levelText := strings.ToUpper(entry.Level.String())
	if label, ok := entry.Data["label"]; ok {
		levelText = label.(string) // Casting to string, ensure label is always a string
	}

	// Check if a custom color ("labelColor") is defined in the log entry's data.
	if color, ok := entry.Data["labelColor"].(string); ok {
		// If a custom color is defined, use it for log output
		logOutput = fmt.Sprintf("%s[%s]%s%s (%s)%s %s", lightGrey, f.packageName, reset, color, levelText, reset, entry.Message)
	} else {
		// If no custom color is defined, use the default color based on the log level.
		logOutput = fmt.Sprintf("%s[%s]%s%s (%s)%s %s", lightGrey, f.packageName, reset, getColor(entry.Level), levelText, reset, entry.Message)
	}

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

// getColor returns the color for the log level
func getColor(level logrus.Level) string {
	switch level {
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return red
	case logrus.WarnLevel:
		return yellow
	case logrus.InfoLevel:
		return blue
	default:
		return reset
	}
}
