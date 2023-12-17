package log

import (
	"fmt"
	"strings"

	"github.com/root4loot/goutils/color"
	"github.com/sirupsen/logrus"
)

type CustomFormatter struct {
	packageName string
}

// Format formats the log output
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelText string

	// Check for 'Results' level using 'Trace' and 'custom_level'
	if entry.Level == logrus.TraceLevel {
		if customLevel, ok := entry.Data["custom_level"]; ok && customLevel == "result" {
			return []byte{}, nil // Optionally, handle Result logs differently
		}
	} else {
		levelText = strings.ToUpper(entry.Level.String())
	}

	// Determine the color for the level text
	var levelColor string
	if entry.Level == logrus.TraceLevel && levelText == "RESULT" {
		levelColor = color.Blue // Custom color for 'Results' level
	} else {
		levelColor = getColor(entry.Level)
	}

	// Formatting the log output
	logOutput := fmt.Sprintf("%s[%s]%s %s(%s)%s %s", color.LightGrey, f.packageName, color.Reset, levelColor, levelText, color.Reset, entry.Message)

	return []byte(logOutput + "\n"), nil
}

func getColor(level logrus.Level) string {
	switch level {
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return color.Red
	case logrus.WarnLevel:
		return color.Yellow
	case logrus.InfoLevel:
		return color.Cyan
	case logrus.TraceLevel: // Handle Trace level separately
		return color.Blue // Default color for Trace level
	default:
		return color.Reset
	}
}
