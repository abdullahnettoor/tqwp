package tqwp

import (
	"fmt"
	"time"
)

var CusLogger *CustomLogger

type CustomLogger struct {
	level string
}

func NewCustomLogger() *CustomLogger {
	return &CustomLogger{}
}

func (l *CustomLogger) log(message string) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	fmt.Printf("%s: %s %s\n", l.level, timestamp, message)
}

func (l *CustomLogger) CustomeTag(tag, message string) {
	l.level = tag
	l.log(message)
}

func (l *CustomLogger) Info(level, message string) {
	l.level = "[INFO] "
	l.log(message)
}

func (l *CustomLogger) Warn(message string) {
	l.level = "[WARN] "
	l.log(message)
}

func (l *CustomLogger) Error(message string) {
	l.level = "[ERROR] "
	l.log(message)
}

func (l *CustomLogger) Success(message string) {
	l.level = "[SUCCESS] "
	l.log(message)
}
