package logger

import (
	"log"
	"os"
)

const (
	LogLevelDebugName  = "DEBUG"
	LogLevelDebugValue = 4

	LogLevelInfoName  = "INFO"
	LogLevelInfoValue = 3

	LogLevelWarnName  = "WARN"
	LogLevelWarnValue = 2

	LogLevelErrorName  = "ERROR"
	LogLevelErrorValue = 1
)

// There is probably a better way to handle log levels
type LogLevel struct {
	Name  string
	Value int
}

func DebugLogLevel() LogLevel {
	debug := LogLevel{
		Name:  LogLevelDebugName,
		Value: LogLevelDebugValue,
	}

	return debug
}

func InfoLogLevel() LogLevel {
	info := LogLevel{
		Name:  LogLevelInfoName,
		Value: LogLevelInfoValue,
	}

	return info
}

func WarnLogLevel() LogLevel {
	warn := LogLevel{
		Name:  LogLevelWarnName,
		Value: LogLevelWarnValue,
	}
	return warn
}

func ErrorLogLevel() LogLevel {
	error := LogLevel{
		Name:  LogLevelErrorName,
		Value: LogLevelErrorValue,
	}
	return error
}

type BaseLogger interface {
	Println(v ...interface{})
}

type CliLogger struct {
	Logger   BaseLogger
	LogLevel LogLevel
}

func (c CliLogger) log(message string, logLevelValue int) {
	if c.LogLevel.Value >= logLevelValue {
		c.Logger.Println(message)
	}
}

func (c CliLogger) Debug(message string) {
	c.log(message, LogLevelDebugValue)
}

func (c CliLogger) Info(message string) {
	c.log(message, LogLevelInfoValue)
}

func (c CliLogger) Warn(message string) {
	c.log(message, LogLevelWarnValue)
}

func (c CliLogger) Error(message string) {
	c.log(message, LogLevelErrorValue)
}

func NewCliLogger(logLevel LogLevel) CliLogger {
	writer := os.Stdout
	loggerFlags := 0

	logger := log.New(writer, "", loggerFlags)

	cliLogger := CliLogger{
		Logger:   logger,
		LogLevel: logLevel,
	}

	return cliLogger
}
