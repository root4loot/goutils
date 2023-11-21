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

func Init(name string) {
	log = NewLogger(name)
}

func (lvl Level) String() string {
	return levels[lvl]
}

type Fields logrus.Fields

type Logger struct {
	*logrus.Logger
	packageName string
}

func NewLogger(packageName string) *Logger {
	l := logrus.New()
	l.SetFormatter(&CustomFormatter{packageName: packageName})
	return &Logger{l, packageName}
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

func GetLevel() Level {
	return convertLogrusLevelToRelogLevel(log.Logger.GetLevel())
}

func WithFields(fields Fields) *logrus.Entry {
	return log.Logger.WithFields(logrus.Fields(fields))
}

// Generalized logging functions that accept any type as an argument
func Debug(v ...interface{}) {
	log.Debug(strings.TrimSpace(fmt.Sprintln(v...)))
}

func Info(v ...interface{}) {
	log.Info(strings.TrimSpace(fmt.Sprintln(v...)))
}

func Warn(v ...interface{}) {
	log.Warn(strings.TrimSpace(fmt.Sprintln(v...)))
}

func Error(v ...interface{}) {
	log.Error(strings.TrimSpace(fmt.Sprintln(v...)))
}

func Fatal(v interface{}) {
	log.Fatal(strings.TrimSpace(fmt.Sprintln(v)))
}

// Formatted logging functions
func Debugf(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

func Infof(format string, v ...interface{}) {
	log.Infof(format, v...)
}

func Warnf(format string, v ...interface{}) {
	log.Warnf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	log.Errorf(format, v...)
}

func Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
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
