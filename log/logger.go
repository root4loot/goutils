package log

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
)

var log *Logger

var levels = []string{
	"debug",
	"info",
	"warn",
	"error",
	"fatal",
	"panic",
}

type Fields logrus.Fields

type Logger struct {
	*logrus.Logger
	packageName string
}

type Label struct {
	label  string
	logger *Logger
	color  string
}

// Init initializes the global logger
func Init(name string) {
	log = NewLogger(name)
}

// String returns the string representation of the log level
func (lvl Level) String() string {
	return levels[lvl]
}

// NewLogger creates a new logger
func NewLogger(packageName string) *Logger {
	l := logrus.New()
	l.SetFormatter(&CustomFormatter{packageName: packageName})
	return &Logger{l, packageName}
}

// NewLabel creates a new label
func NewLabel(label string) *Label {
	return &Label{label: label}
}

// NewLabel creates a new label
func (l *Logger) NewLabel(label string) *Label {
	return &Label{logger: l, label: label}
}

// SetLevel sets the logger level
func SetLevel(level Level) {
	log.Logger.SetLevel(convertRelogLevelToLogrusLevel(level))
}

// SetLevelStr sets the logger level with a string
func SetLevelStr(level string) {
	level = strings.ToLower(level)
	log.Logger.SetLevel(convertRelogLevelToLogrusLevel(getLogLevel(level)))
}

// GetLevel returns the logger level
func GetLevel() Level {
	return convertLogrusLevelToRelogLevel(log.Logger.GetLevel())
}

// WithFields returns a logrus entry with the fields set
func WithFields(fields Fields) *logrus.Entry {
	return log.Logger.WithFields(logrus.Fields(fields))
}

// Log logs a message using the logger.
// If no logger is set, it uses the global logger.
func (l *Label) Log(v ...interface{}) {
	// Check if a logger is set, if not, use the global logger
	if l.logger == nil {
		l.logger = log
	}

	// Create a log entry with label and tagColor fields
	entry := l.logger.WithFields(logrus.Fields{
		"label":    l.label, // Label associated with this log
		"tagColor": l.color, // Color associated with the label
	})

	// Log the provided values as an Info message
	entry.Info(v...)
}

func (l *Label) Logf(format string, v ...interface{}) {
	if l.logger == nil {
		l.logger = log
	}

	entry := l.logger.WithFields(logrus.Fields{
		"label":    l.label,
		"tagColor": l.color,
	})
	entry.Infof(format, v...)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(v ...interface{}) {
	if log != nil {
		log.Debug(strings.TrimSpace(fmt.Sprintln(v...)))
	}
}

// Info logs a message at level Info on the standard logger.
func Info(v ...interface{}) {
	if log != nil {
		log.Info(strings.TrimSpace(fmt.Sprintln(v...)))
	}
}

// Warn logs a message at level Warn on the standard logger.
func Warn(v ...interface{}) {
	if log != nil {
		log.Warn(strings.TrimSpace(fmt.Sprintln(v...)))
	}
}

// Error logs a message at level Error on the standard logger.
func Error(v ...interface{}) {
	if log != nil {
		log.Error(strings.TrimSpace(fmt.Sprintln(v...)))
	}
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(v interface{}) {
	if log != nil {
		log.Fatal(strings.TrimSpace(fmt.Sprintln(v)))
	}
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, v ...interface{}) {
	if log != nil {
		log.Debugf(format, v...)
	}
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, v ...interface{}) {
	if log != nil {
		log.Infof(format, v...)
	}
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, v ...interface{}) {
	if log != nil {
		log.Warnf(format, v...)
	}
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, v ...interface{}) {
	if log != nil {
		log.Errorf(format, v...)
	}
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, v ...interface{}) {
	if log != nil {
		log.Fatalf(format, v...)
	}
}

// getLogLevel returns the log level based on passed string
func getLogLevel(level string) Level {
	switch level {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	case "panic":
		return PanicLevel
	default:
		return InfoLevel
	}
}

// convertLogrusLevelToRelogLevel converts a logrus level to a relog level
func convertLogrusLevelToRelogLevel(level logrus.Level) Level {
	switch level {
	case logrus.DebugLevel:
		return DebugLevel
	case logrus.InfoLevel:
		return InfoLevel
	case logrus.WarnLevel:
		return WarnLevel
	case logrus.ErrorLevel:
		return ErrorLevel
	case logrus.FatalLevel:
		return FatalLevel
	case logrus.PanicLevel:
		return PanicLevel
	default:
		return InfoLevel
	}
}

// convertRelogLevelToLogrusLevel converts a relog level to a logrus level
func convertRelogLevelToLogrusLevel(level Level) logrus.Level {
	switch level {
	case DebugLevel:
		return logrus.DebugLevel
	case InfoLevel:
		return logrus.InfoLevel
	case WarnLevel:
		return logrus.WarnLevel
	case ErrorLevel:
		return logrus.ErrorLevel
	case FatalLevel:
		return logrus.FatalLevel
	case PanicLevel:
		return logrus.PanicLevel
	default:
		return logrus.InfoLevel
	}
}

// SetColor sets the label color
func (l *Label) SetColor(color string) {
	if color == "grey" {
		l.color = lightGrey
	}
	if color == "red" {
		l.color = red
	}
	if color == "blue" {
		l.color = blue
	}
	if color == "yellow" {
		l.color = yellow
	}
	if color == "grey" {
		l.color = lightGrey
	}
}
