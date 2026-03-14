package observability

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Level represents a log level.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Logger is a structured logger.
type Logger struct {
	level  Level
	service string
	logger *log.Logger
}

// NewLogger creates a new logger for a service.
func NewLogger(service string, level Level) *Logger {
	return &Logger{
		level:  level,
		service: service,
		logger: log.New(os.Stderr, "", 0),
	}
}

func (l *Logger) log(level Level, levelStr, msg string, fields map[string]interface{}) {
	if level < l.level {
		return
	}
	ts := time.Now().UTC().Format(time.RFC3339)
	extra := ""
	for k, v := range fields {
		extra += fmt.Sprintf(" %s=%v", k, v)
	}
	l.logger.Printf("%s level=%s service=%s msg=%q%s", ts, levelStr, l.service, msg, extra)
}

// Info logs at info level.
func (l *Logger) Info(msg string, fields map[string]interface{}) {
	l.log(LevelInfo, "INFO", msg, fields)
}

// Error logs at error level.
func (l *Logger) Error(msg string, fields map[string]interface{}) {
	l.log(LevelError, "ERROR", msg, fields)
}

// Warn logs at warn level.
func (l *Logger) Warn(msg string, fields map[string]interface{}) {
	l.log(LevelWarn, "WARN", msg, fields)
}

// Debug logs at debug level.
func (l *Logger) Debug(msg string, fields map[string]interface{}) {
	l.log(LevelDebug, "DEBUG", msg, fields)
}
