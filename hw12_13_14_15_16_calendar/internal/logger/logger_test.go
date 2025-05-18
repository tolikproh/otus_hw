package logger

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/cnst"
)

func TestLogger(t *testing.T) {
	msg := "hello world!"
	key := "key_test"
	value := "value_test"

	message := func(level, msg, key, value string) string {
		level = strings.ToUpper(level)
		const jsonTimeRE = `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?(Z|[+-]\d{2}:\d{2})`
		return `{"time":"` + jsonTimeRE + `","level":"` + level + `","msg":"` + msg + `","` + key + `":"` + value + `"}`
	}

	testCase := []struct {
		name  string
		level string
		msg   string
		key   string
		value string
		exp   []string
	}{
		{
			name:  "Debug level",
			level: cnst.LoggerLevelDebug,
			msg:   msg,
			key:   key,
			value: value,
			exp: []string{
				message(cnst.LoggerLevelDebug, msg, key, value),
				message(cnst.LoggerLevelInfo, msg, key, value),
				message(cnst.LoggerLevelWarn, msg, key, value),
				message(cnst.LoggerLevelError, msg, key, value),
			},
		},
		{
			name:  "Info level",
			level: cnst.LoggerLevelInfo,
			msg:   msg,
			key:   key,
			value: value,
			exp: []string{
				message(cnst.LoggerLevelInfo, msg, key, value),
				message(cnst.LoggerLevelWarn, msg, key, value),
				message(cnst.LoggerLevelError, msg, key, value),
			},
		},
		{
			name:  "Warn level",
			level: cnst.LoggerLevelWarn,
			msg:   msg,
			key:   key,
			value: value,
			exp: []string{
				message(cnst.LoggerLevelWarn, msg, key, value),
				message(cnst.LoggerLevelError, msg, key, value),
			},
		},
		{
			name:  "Error level",
			level: cnst.LoggerLevelError,
			msg:   msg,
			key:   key,
			value: value,
			exp: []string{
				message(cnst.LoggerLevelError, msg, key, value),
			},
		},
	}

	for _, tc := range testCase {
		var buf bytes.Buffer
		defer buf.Reset()

		log := newLogger(&buf, tc.level)

		log.Debug(tc.msg, tc.key, tc.value)
		log.Info(tc.msg, tc.key, tc.value)
		log.Warn(tc.msg, tc.key, tc.value)
		log.Error(tc.msg, tc.key, tc.value)

		str := strings.Split(buf.String(), "\n")
		for i := 0; i < len(str)-1; i++ {
			assert.Regexpf(t, tc.exp[i], str[i], "log formatted error, message: %s", str[i])
		}
	}
}

func TestLoggerLevel(t *testing.T) {
	testCase := []struct {
		name  string
		level string
		exp   slog.Level
	}{
		{
			name:  "Upper Info",
			level: "INFO",
			exp:   slog.LevelInfo,
		},
		{
			name:  "Upper Warn",
			level: "WARN",
			exp:   slog.LevelWarn,
		},
		{
			name:  "Upper Error",
			level: "ERROR",
			exp:   slog.LevelError,
		},
		{
			name:  "Upper Debug",
			level: "DEBUG",
			exp:   slog.LevelDebug,
		},
		{
			name:  "Lower Info",
			level: cnst.LoggerLevelInfo,
			exp:   slog.LevelInfo,
		},
		{
			name:  "Lower Warn",
			level: cnst.LoggerLevelWarn,
			exp:   slog.LevelWarn,
		},
		{
			name:  "Lower Error",
			level: cnst.LoggerLevelError,
			exp:   slog.LevelError,
		},
		{
			name:  "Lower Debug",
			level: cnst.LoggerLevelDebug,
			exp:   slog.LevelDebug,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			sLogLevel := getLogLevel(tc.level)
			assert.EqualValues(t, tc.exp, sLogLevel)
		})
	}
}
