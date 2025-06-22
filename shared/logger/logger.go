package logger

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

type Logger struct {
	service string
	logger  *log.Logger
}

func NewLogger(service string) *Logger {
	logFileName := fmt.Sprintf("%s.log", service)
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("failed to open or create log file '%s': %v", logFileName, err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	logger := log.New(multiWriter, "", log.LstdFlags)

	return &Logger{
		service,
		logger,
	}
}

func (l *Logger) string(variant string, action string, message string, args ...any) (out string) {
	head := fmt.Sprintf("%s [%s] (%s) ", variant, l.service, action)
	body := fmt.Sprintf(message, args...)
	out = fmt.Sprintf("%s %s", head, body)
	return
}

func (l *Logger) GetLogString(action string, message string, args ...any) (out string) {
	out = l.string("LOG", action, message, args...)
	return
}

func (l *Logger) GetErrorString(action string, message string, args ...any) (out string) {
	out = l.string("ERROR", action, message, args...)
	return
}

func (l *Logger) Error(action string, message string, args ...any) error {
	return errors.New(l.GetErrorString(action, message, args...))
}

func (l *Logger) Log(action string, message string, args ...any) {
	l.logger.Println(l.GetLogString(action, message, args...))
}

func (l *Logger) LogError(action string, message string, args ...any) {
	l.logger.Println(l.GetErrorString(action, message, args...))
}

func (l *Logger) SError(action string, message string, args ...any) error {
	content := l.GetErrorString(action, message, args...)
	l.logger.Println(content)
	return errors.New(content)
}

func (l *Logger) FatalLog(action string, message string, args ...any) {
	l.logger.Fatal(l.GetErrorString(action, message, args...))
}
