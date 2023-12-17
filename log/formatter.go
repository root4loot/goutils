package log

import (
	"fmt"
	"strings"

	"github.com/root4loot/goutils/color"
	"github.com/sirupsen/logrus"
	"golang.org/x/term"
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
	if labelColor, ok := entry.Data["labelColor"].(string); ok {
		// If a custom color is defined, use it for log output
		logOutput = fmt.Sprintf("%s[%s]%s%s (%s)%s %s", color.LightGrey, f.packageName, color.Reset, labelColor, levelText, color.Reset, entry.Message)
	} else {
		// If no custom color is defined, use the default color based on the log level.
		logOutput = fmt.Sprintf("%s[%s]%s%s (%s)%s %s", color.LightGrey, f.packageName, color.Reset, getColor(entry.Level), levelText, color.Reset, entry.Message)
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
		return color.Red
	case logrus.WarnLevel:
		return color.Yellow
	case logrus.InfoLevel:
		return color.Cyan
	default:
		return color.Reset
	}
}
