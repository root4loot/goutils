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

	if entry.Level == logrus.TraceLevel {
		if customLevel, ok := entry.Data["custom_level"]; ok && customLevel == "result" {
			levelText = "RES"
		}
	} else {
		originalLevelText := strings.ToUpper(entry.Level.String())
		levelText = logLevelAbbreviations[originalLevelText] // Get abbreviation or default to original
	}

	var levelColor string
	if entry.Level == logrus.TraceLevel && levelText == "RES" {
		levelColor = color.Blue
	} else {
		levelColor = getColor(entry.Level)
	}

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
		return color.Blue
	case logrus.TraceLevel: // Handle Trace level separately
		return color.Blue
	default:
		return color.Reset
	}
}
