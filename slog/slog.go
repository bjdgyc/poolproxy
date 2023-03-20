package slog

import (
	"fmt"
	"os"
)

type LogLevel int

const (
	FATAL LogLevel = iota
	ERROR
	WARN
	INFO
	DEBUG
)

var levelName = map[string]LogLevel{
	"FATAL": FATAL,
	"ERROR": ERROR,
	"WARN":  WARN,
	"INFO":  INFO,
	"DEBUG": DEBUG,
}

func (l *Logger) Fatal(args ...interface{}) {
	if FATAL > l.maxLogLevel {
		return
	}
	l.output("FATAL", fmt.Sprint(args...))
	os.Exit(1)
}

func (l *Logger) Fatalf(msg string, args ...interface{}) {
	if FATAL > l.maxLogLevel {
		return
	}
	l.output("FATAL", fmt.Sprintf(msg, args...))
	os.Exit(1)
}

// Error logs a message to the 'standard' Logger (always)
func (l *Logger) Error(args ...interface{}) {
	if ERROR > l.maxLogLevel {
		return
	}
	l.output("ERROR", fmt.Sprint(args...))
}

func (l *Logger) Errorf(msg string, args ...interface{}) {
	if ERROR > l.maxLogLevel {
		return
	}
	l.output("ERROR", fmt.Sprintf(msg, args...))
}

// Warn logs a message to the 'standard' Logger if MaxLogLevel is >= WARN
func (l *Logger) Warn(args ...interface{}) {
	if WARN > l.maxLogLevel {
		return
	}
	l.output("WARN", fmt.Sprint(args...))
}

func (l *Logger) Warnf(msg string, args ...interface{}) {
	if WARN > l.maxLogLevel {
		return
	}
	l.output("WARN", fmt.Sprintf(msg, args...))
}

// Info logs a message to the 'standard' Logger if MaxLogLevel is >= INFO
func (l *Logger) Info(args ...interface{}) {
	if INFO > l.maxLogLevel {
		return
	}
	l.output("INFO", fmt.Sprint(args...))
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	if INFO > l.maxLogLevel {
		return
	}
	l.output("INFO", fmt.Sprintf(msg, args...))
}

// Trace logs a message to the 'standard' Logger if MaxLogLevel is >= DEBUG
func (l *Logger) Debug(args ...interface{}) {
	if DEBUG > l.maxLogLevel {
		return
	}
	l.output("DEBUG", fmt.Sprint(args...))
}

func (l *Logger) Debugf(msg string, args ...interface{}) {
	if DEBUG > l.maxLogLevel {
		return
	}
	l.output("DEBUG", fmt.Sprintf(msg, args...))
}

func (l *Logger) output(mode, msg string) {
	l.logout.Output(3, "["+mode+"] "+msg)
}
