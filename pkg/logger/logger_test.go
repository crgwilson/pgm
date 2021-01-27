package logger

import (
	"testing"

	"github.com/crgwilson/pgm/pkg/mocks"
)

func TestCliLogger(t *testing.T) {
	cases := []struct {
		Name         string
		LogLevel     LogLevel
		ExpectedLogs []string
	}{
		{
			"debug log level",
			DebugLogLevel(),
			[]string{"debug", "info", "warn", "error"},
		},
		{
			"info log level",
			InfoLogLevel(),
			[]string{"info", "warn", "error"},
		},
		{
			"warn log level",
			WarnLogLevel(),
			[]string{"warn", "error"},
		},
		{
			"error log level",
			ErrorLogLevel(),
			[]string{"error"},
		},
	}

	for _, test := range cases {
		spyLogger := mocks.NewSpyLogger()
		testLogger := CliLogger{
			Logger:   spyLogger,
			LogLevel: test.LogLevel,
		}
		testLogger.Debug("debug")
		testLogger.Info("info")
		testLogger.Warn("warn")
		testLogger.Error("error")

		for i := range spyLogger.Logs {
			if spyLogger.Logs[i] != test.ExpectedLogs[i] {
				t.Errorf("got %q, want %q", spyLogger.Logs[i], test.ExpectedLogs[i])
			}
		}
	}
}
