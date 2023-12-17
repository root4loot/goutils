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

type Tag struct {
	tag    string
	logger *Logger
}

// Init initializes the global logger
func Init(name string) {
	log = NewLogger(name)
}

func (lvl Level) String() string {
	return levels[lvl]
}

// NewLogger creates a new logger
func NewLogger(packageName string) *Logger {
	l := logrus.New()
	l.SetFormatter(&CustomFormatter{packageName: packageName})
	return &Logger{l, packageName}
}

func NewTag(tag string) *Tag {
	return &Tag{tag: tag}
}

func (l *Logger) NewTag(tag string) *Tag {
	return &Tag{logger: l, tag: tag}
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

// entry returns a logrus entry with the tag field set
func (t *Tag) entry() *logrus.Entry {
	if t.logger != nil {
		return t.logger.WithFields(logrus.Fields{"tag": t.tag})
	} else {
		return log.Logger.WithFields(logrus.Fields{"tag": t.tag})
	}
}

func (t *Tag) Log(v ...interface{}) {
	entry := t.entry()
	if entry != nil {
		entry.Info(v...)
	}
}

func (t *Tag) Logf(format string, v ...interface{}) {
	entry := t.entry()
	if entry != nil {
		entry.Infof(format, v...)
	}
}

// Generalized logging functions that accept any type as an argument
func Debug(v ...interface{}) {
	if log != nil {
		log.Debug(strings.TrimSpace(fmt.Sprintln(v...)))
	}
}
func Info(v ...interface{}) {
	if log != nil {
		log.Info(strings.TrimSpace(fmt.Sprintln(v...)))
	}
}

func Warn(v ...interface{}) {
	if log != nil {
		log.Warn(strings.TrimSpace(fmt.Sprintln(v...)))
	}
}

func Error(v ...interface{}) {
	if log != nil {
		log.Error(strings.TrimSpace(fmt.Sprintln(v...)))
	}
}

func Fatal(v interface{}) {
	if log != nil {
		log.Fatal(strings.TrimSpace(fmt.Sprintln(v)))
	}
}

// Formatted logging functions
func Debugf(format string, v ...interface{}) {
	if log != nil {
		log.Debugf(format, v...)
	}
}

func Infof(format string, v ...interface{}) {
	if log != nil {
		log.Infof(format, v...)
	}
}

func Warnf(format string, v ...interface{}) {
	if log != nil {
		log.Warnf(format, v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if log != nil {
		log.Errorf(format, v...)
	}
}

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
