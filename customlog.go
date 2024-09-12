package tqwp

import (
	"fmt"
	"time"
)

type customLogger struct {
	level string
}

func NewCustomLogger() *customLogger {
	return &customLogger{}
}

func (l *customLogger) log(message string) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	fmt.Printf("%s: %s %s\n", l.level, timestamp, message)
}

func (l *customLogger) CustomTag(tag, message string) {
	l.level = tag
	l.log(message)
}

func (l *customLogger) Info(level, message string) {
	l.level = "[INFO] "
	l.log(message)
}

func (l *customLogger) Warn(message string) {
	l.level = "[WARN] "
	l.log(message)
}

func (l *customLogger) Error(message string) {
	l.level = "[ERROR] "
	l.log(message)
}

func (l *customLogger) Success(message string) {
	l.level = "[SUCCESS] "
	l.log(message)
}
