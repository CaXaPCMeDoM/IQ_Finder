package logger

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
	level string
}

func New(level string) *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", log.LstdFlags),
		level:  level,
	}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level == "debug" {
		l.Printf("[DEBUG] "+format, v...)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.level == "debug" || l.level == "info" {
		l.Printf("[INFO] "+format, v...)
	}
}

func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level == "debug" || l.level == "info" || l.level == "warn" {
		l.Printf("[WARN] "+format, v...)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.Printf("[ERROR] "+format, v...)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.Printf("[FATAL] "+format, v...)
	os.Exit(1)
}
