package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/root4loot/goutils/color"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

type Level int

const PipedOutputNotification = "Notice: Output is being piped. Results will be formatted accordingly."

const (
	DebugLevel Level = iota
	TraceLevel       // Custom level
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
)

var log *Logger

var levels = []string{
	"debug",
	"trace",
	"info",
	"warn",
	"error",
	"fatal",
	"panic",
}

var logLevelAbbreviations = map[string]string{
	"WARNING": "WAR",
	"INFO":    "INF",
	"DEBUG":   "DEB",
	"RESULT":  "RES",
	"ERROR":   "ERR",
	"FATAL":   "FAT",
	"PANIC":   "PAN",
	// Add other levels as needed
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

// Init initializes the global logger and sets the log level to Info.
func Init(name string) {
	log = NewLogger(name)
	log.SetLevel(logrus.InfoLevel)
}

// Notify logs a message using a specified or global logger.
func Notify(v ...interface{}) {
	var logger *Logger
	if len(v) > 0 {
		var ok bool
		if logger, ok = v[0].(*Logger); ok {
			v = v[1:]
		} else {
			logger = log
		}
	} else {
		logger = log
	}

	if IsOutputPiped() {
		message := fmt.Sprint(v...)
		packageNameFormatted := color.Colorize(color.LightGrey, fmt.Sprintf("[%s]", logger.packageName))
		infoLabelFormatted := color.Colorize(color.Yellow, fmt.Sprintf(" (%s)", logLevelAbbreviations["INFO"]))
		fmt.Fprint(os.Stderr, packageNameFormatted+infoLabelFormatted+" "+message+"\n")
	}
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

// Logf logs message at the Info level using the label's logger.
// If the label's logger is not set, it defaults to the global logger.
func (l *Label) Log(v ...interface{}) {
	if l.logger == nil {
		l.logger = log
	}

	entry := l.logger.WithFields(logrus.Fields{
		"label":      l.label,
		"labelColor": l.color,
	})

	entry.Info(v...)
}

// Logf logs a formatted message at the Info level using the label's logger.
func (l *Label) Logf(format string, v ...interface{}) {
	if l.logger == nil {
		l.logger = log // Default to the global logger
	}

	entry := l.logger.WithFields(logrus.Fields{
		"label":      l.label,
		"labelColor": l.color,
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

// Result logs a message using a specified or global logger.
// The log is printed to stdout, and metadata is printed to stderr if not piped.
func Result(v ...interface{}) {
	var logger *Logger
	if len(v) > 0 {
		var ok bool
		if logger, ok = v[0].(*Logger); ok {
			v = v[1:]
		} else {
			logger = log
		}
	} else {
		logger = log
	}

	message := fmt.Sprint(v...)
	if !IsOutputPiped() {
		levelAbbrev := logLevelAbbreviations["RESULT"]
		packageNameFormatted := color.Colorize(color.LightGrey, fmt.Sprintf("[%s]", logger.packageName))
		levelLabelFormatted := color.Colorize(color.Cyan, fmt.Sprintf(" (%s)", levelAbbrev))
		fmt.Fprint(os.Stderr, packageNameFormatted+levelLabelFormatted+" ")
	}
	// Print the message to stdout
	fmt.Fprintln(os.Stdout, message)
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

// Resultf logs a formatted message using a specified or global logger.
// The log is printed to stdout, and metadata is printed to stderr if not piped.
func Resultf(format string, v ...interface{}) {
	var logger *Logger
	if len(v) > 0 {
		var ok bool
		if logger, ok = v[0].(*Logger); ok {
			v = v[1:]
		} else {
			logger = log
		}
	} else {
		logger = log
	}

	message := fmt.Sprintf(format, v...)
	if !IsOutputPiped() {
		levelAbbrev := logLevelAbbreviations["RESULT"]
		packageNameFormatted := color.Colorize(color.LightGrey, fmt.Sprintf("[%s]", logger.packageName))
		levelLabelFormatted := color.Colorize(color.Cyan, fmt.Sprintf(" (%s)", levelAbbrev))
		fmt.Fprint(os.Stderr, packageNameFormatted+levelLabelFormatted+" ")
	}
	// Print the formatted message to stdout
	fmt.Fprintln(os.Stdout, message)
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

func IsOutputPiped() bool {
	return !terminal.IsTerminal(int(os.Stdout.Fd()))
}

// getLogLevel returns the log level based on passed string
func getLogLevel(level string) Level {
	switch level {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "trace":
		return TraceLevel
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
	case logrus.TraceLevel:
		return TraceLevel
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
	case TraceLevel:
		return logrus.TraceLevel
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
	l.color = color
}
