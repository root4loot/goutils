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
)

type CustomFormatter struct {
	packageName string
}

// Format formats the log output
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	color := getColor(entry.Level)

	// Start with the basic log format
	logOutput := fmt.Sprintf("%s[%s]%s%s (%s)%s %s", lightGrey, f.packageName, reset, color, strings.ToUpper(entry.Level.String()), reset, entry.Message)

	// Prepare fields output
	fieldsOutput := ""
	if len(entry.Data) > 0 {
		for key, value := range entry.Data {
			fieldsOutput += fmt.Sprintf(" %s=%v", key, value) // Adjust this to colorize key/values
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
