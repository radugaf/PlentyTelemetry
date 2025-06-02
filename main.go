package main

import (
	"fmt"
	"time"
)

type LogLevel string

const (
	Debug   LogLevel = "DEBUG"
	Info    LogLevel = "INFO"
	Warning LogLevel = "WARNING"
	Error   LogLevel = "ERROR"
)

type LogEntry struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
}

type Driver interface {
	Write(logEntry LogEntry)
}

type Logger struct {
	drivers []Driver
}

func NewLogger(drivers ...Driver) *Logger {
	return &Logger{drivers: drivers}
}

func (l *Logger) log(level LogLevel, message string) {
	le := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
	}

	for _, driver := range l.drivers {
		driver.Write(le)
	}
}

func (l *Logger) Info(msg string) {
	l.log(Info, msg)
}

func (l *Logger) Debug(msg string) {
	l.log(Debug, msg)
}

func (l *Logger) Error(msg string) {
	l.log(Error, msg)
}

// CLI driver
type CLIDriver struct{}

func NewCLIDriver() *CLIDriver {
	return &CLIDriver{}
}

func (d *CLIDriver) Write(logEntry LogEntry) {
	fmt.Printf("[%s] %s: %s\n",
		logEntry.Timestamp.Format("2006-01-02 15:04:05"),
		logEntry.Level,
		logEntry.Message,
	)
}

func main() {
	cliAdapter := NewCLIDriver()
	logger := NewLogger(cliAdapter)
	logger.Info("Test Started")
	logger.Debug("Debug")
	logger.Error("Error")
}
