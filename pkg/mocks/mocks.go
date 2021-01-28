package mocks

import (
	"fmt"
)

type SpyLogger struct {
	Logs []string
}

func (s SpyLogger) Println(v ...interface{}) {
	newLogs := append(s.Logs, fmt.Sprint(v...))
	s.Logs = newLogs
}

func NewSpyLogger() *SpyLogger {
	spy := SpyLogger{
		Logs: make([]string, 0),
	}

	return &spy
}
