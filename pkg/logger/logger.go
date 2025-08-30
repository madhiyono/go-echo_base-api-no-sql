package logger

import (
	"log"
	"os"
)

type Logger struct {
	infoLogger *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
}

func New(level string) *Logger {
	return &Logger{
		infoLogger: log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		fatalLogger: log.New(os.Stdout, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.infoLogger.Printf(format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.errorLogger.Printf(format, v...)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.fatalLogger.Printf(format, v...)
}