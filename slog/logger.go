package slog

import (
	"log"
	"os"
	"strings"
)

type Logger struct {
	logout      *log.Logger
	maxLogLevel LogLevel
	logflag     int
	dateFormat  string
}

func GetStdLog() *Logger {
	logflag := log.LstdFlags | log.Lshortfile
	return &Logger{
		logout:      log.New(os.Stderr, "", logflag),
		maxLogLevel: DEBUG,
		logflag:     logflag,
		dateFormat:  "2006-01-02",
	}
}


func (l *Logger) SetLogfile(outfile string) {
	fileWriter, err := os.OpenFile(outfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		l.Fatal(outfile, err)
	}
	l.logout = log.New(fileWriter, "", l.logflag)
}

// SetLogLevel sets MaxLogLevel based on the provided string
func (l *Logger) SetLogLevel(level string) {
	level = strings.ToUpper(level)
	lev, ok := levelName[level]
	if !ok {
		log.Fatalf("Unknown log level requested: %v", level)
	}
	l.maxLogLevel = lev
}

func (l *Logger) GetLogLevel() LogLevel {
	return l.maxLogLevel
}
