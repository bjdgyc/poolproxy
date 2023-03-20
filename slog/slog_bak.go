package slog

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"strings"
// 	"sync"
// 	"time"
// )
//
// type LogLevel int
//
// const (
// 	FATAL LogLevel = iota
// 	ERROR
// 	WARN
// 	INFO
// 	DEBUG
// )
//
// var levelName = map[string]LogLevel{
// 	"FATAL": FATAL,
// 	"ERROR": ERROR,
// 	"WARN":  WARN,
// 	"INFO":  INFO,
// 	"DEBUG": DEBUG,
// }
//
// var (
// 	maxLogLevel = DEBUG
// 	logflag     = log.LstdFlags | log.Lshortfile
// 	logout      *log.Logger
// 	logAccess   = make(map[string]*AccessLog)
// 	dateFormat  = "2006-01-02"
// )
//
// // request 按天分割日志
// type AccessLog struct {
// 	lock    sync.Mutex
// 	oldDate string
// 	logfile string
// 	fd      *os.File
// 	logger  *log.Logger
// }
//
// func init() {
// 	logout = log.New(os.Stdout, "", logflag)
// }
//
// func SetLogfile(outfile string) {
// 	fileWriter, err := os.OpenFile(outfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
// 	if err != nil {
// 		Fatal(outfile, err)
// 	}
// 	logout = log.New(fileWriter, "", logflag)
// }
//
// // 设置Access对象
// func SetAccessFile(name, accessFile string) {
// 	accessWriter := createAccessLogger(accessFile)
// 	access := &AccessLog{
// 		logfile: accessFile,
// 		fd:      accessWriter,
// 		oldDate: time.Now().Format(dateFormat),
// 	}
// 	access.logger = log.New(accessWriter, "", log.LstdFlags)
// 	logAccess[name] = access
// }
//
// func createAccessLogger(accessFile string) *os.File {
// 	requestWriter, err := os.OpenFile(accessFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
// 	if err != nil {
// 		Fatal(accessFile, err)
// 	}
// 	return requestWriter
// }
//
// // SetLogLevel sets MaxLogLevel based on the provided string
// func SetLogLevel(level string) {
// 	level = strings.ToUpper(level)
// 	lev, ok := levelName[level]
// 	if !ok {
// 		log.Fatalf("Unknown log level requested: %v", level)
// 	}
// 	maxLogLevel = lev
// }
//
// func GetLogLevel() LogLevel {
// 	return maxLogLevel
// }
//
// func Fatal(args ...interface{}) {
// 	if FATAL > maxLogLevel {
// 		return
// 	}
// 	output("FATAL", fmt.Sprint(args...))
// 	os.Exit(1)
// }
//
// func Fatalf(msg string, args ...interface{}) {
// 	if FATAL > maxLogLevel {
// 		return
// 	}
// 	output("FATAL", fmt.Sprintf(msg, args...))
// 	os.Exit(1)
// }
//
// // Error logs a message to the 'standard' Logger (always)
// func Error(args ...interface{}) {
// 	if ERROR > maxLogLevel {
// 		return
// 	}
// 	output("ERROR", fmt.Sprint(args...))
// }
//
// func Errorf(msg string, args ...interface{}) {
// 	if ERROR > maxLogLevel {
// 		return
// 	}
// 	output("ERROR", fmt.Sprintf(msg, args...))
// }
//
// // Warn logs a message to the 'standard' Logger if MaxLogLevel is >= WARN
// func Warn(args ...interface{}) {
// 	if WARN > maxLogLevel {
// 		return
// 	}
// 	output("WARN", fmt.Sprint(args...))
// }
//
// func Warnf(msg string, args ...interface{}) {
// 	if WARN > maxLogLevel {
// 		return
// 	}
// 	output("WARN", fmt.Sprintf(msg, args...))
// }
//
// // Info logs a message to the 'standard' Logger if MaxLogLevel is >= INFO
// func Info(args ...interface{}) {
// 	if INFO > maxLogLevel {
// 		return
// 	}
// 	output("INFO", fmt.Sprint(args...))
// }
//
// func Infof(msg string, args ...interface{}) {
// 	if INFO > maxLogLevel {
// 		return
// 	}
// 	output("INFO", fmt.Sprintf(msg, args...))
// }
//
// // Trace logs a message to the 'standard' Logger if MaxLogLevel is >= DEBUG
// func Debug(args ...interface{}) {
// 	if DEBUG > maxLogLevel {
// 		return
// 	}
// 	output("DEBUG", fmt.Sprint(args...))
// }
//
// func Debugf(msg string, args ...interface{}) {
// 	if DEBUG > maxLogLevel {
// 		return
// 	}
// 	output("DEBUG", fmt.Sprintf(msg, args...))
// }
//
// func Access(name string, args ...interface{}) {
// 	var (
// 		access *AccessLog
// 		ok     bool
// 	)
// 	if access, ok = logAccess[name]; !ok {
// 		return
// 	}
//
// 	outputAccess(access, fmt.Sprint(args...))
// }
//
// func output(mode, msg string) {
// 	logout.Output(3, "["+mode+"] "+msg)
// }
//
// func outputAccess(access *AccessLog, msg string) {
// 	nowDate := time.Now().Format(dateFormat)
// 	if access.oldDate != nowDate {
// 		access.lock.Lock()
// 		defer access.lock.Unlock()
// 		oldDate := access.oldDate
// 		access.oldDate = nowDate
// 		access.fd.Close()
// 		err := os.Rename(access.logfile, access.logfile+oldDate)
// 		if err != nil {
// 			Error(err)
// 		}
// 		requestWriter := createAccessLogger(access.logfile)
// 		access.fd = requestWriter
// 		access.logger = log.New(requestWriter, "", log.LstdFlags)
// 	}
//
// 	access.logger.Output(3, msg)
// }
