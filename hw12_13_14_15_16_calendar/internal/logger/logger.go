package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/cnst"
)

type Logger struct {
	*slog.Logger
}

func New(level string) *Logger {
	return newLogger(os.Stdout, level)
}

func newLogger(w io.Writer, level string) *Logger {
	l := slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: getLogLevel(level),
	}))

	return &Logger{l}
}

func getLogLevel(level string) slog.Level {
	level = strings.ToLower(level)

	switch level {
	case cnst.LoggerLevelDebug:
		return slog.LevelDebug // DebugLevel = -4
	case cnst.LoggerLevelWarn:
		return slog.LevelWarn // WarnLevel = 4
	case cnst.LoggerLevelError:
		return slog.LevelError // ErrorLevel = 8
	default:
		return slog.LevelInfo // InfoLevel = 0
	}
}
