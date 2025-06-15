package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/cnst"
)

func TestConfig(t *testing.T) {
	testCases := []struct {
		name string
		path string
		file string
		env  map[string]string
		exp  *Config
	}{
		{
			name: "default config",
			path: "",
			file: "",
			env:  map[string]string{},
			exp: &Config{
				HTTPServer: HTTPServer{
					Host: "localhost",
					Port: 8080,
				},
				GRPCServer: GRPCServer{
					Host: "localhost",
					Port: 5000,
				},
				Logger: Logger{
					Level: cnst.LoggerLevelInfo,
				},
				Storage: Storage{
					Type: cnst.StorageTypeMemory,
					Conn: "",
				},
			},
		},
		{
			name: "config file yaml",
			path: "./testdata",
			file: "config_test.yaml",
			env:  map[string]string{},
			exp: &Config{
				HTTPServer: HTTPServer{
					Host: "calendar.ru",
					Port: 8000,
				},
				GRPCServer: GRPCServer{
					Host: "calendar.ru",
					Port: 5001,
				},
				Logger: Logger{
					Level: cnst.LoggerLevelError,
				},
				Storage: Storage{
					Type: cnst.StorageTypePostgres,
					Conn: "postgres://user:password@localhost:5432/calendar_db?sslmode=disable",
				},
			},
		},
		{
			name: "config environment",
			path: "./testdata",
			file: "config_test.yaml",
			env: map[string]string{
				"CALENDAR_HTTP_SERVER_HOST": "env.net",
				"CALENDAR_HTTP_SERVER_PORT": "1234",
				"CALENDAR_GRPC_SERVER_HOST": "env.net",
				"CALENDAR_GRPC_SERVER_PORT": "5002",
				"CALENDAR_LOG_LEVEL":        "warn",
				"CALENDAR_STORAGE_TYPE":     "memory",
				"CALENDAR_STORAGE_CONN":     "http://memory.com",
			},
			exp: &Config{
				HTTPServer: HTTPServer{
					Host: "env.net",
					Port: 1234,
				},
				GRPCServer: GRPCServer{
					Host: "env.net",
					Port: 5002,
				},
				Logger: Logger{
					Level: cnst.LoggerLevelWarn,
				},
				Storage: Storage{
					Type: cnst.StorageTypeMemory,
					Conn: "http://memory.com",
				},
			},
		},
		{
			name: "config file and secret environment",
			path: "./testdata",
			file: "config_test.yaml",
			env: map[string]string{
				"CALENDAR_STORAGE_CONN": "postgres://user:secret@localhost:5555",
			},
			exp: &Config{
				HTTPServer: HTTPServer{
					Host: "calendar.ru",
					Port: 8000,
				},
				GRPCServer: GRPCServer{
					Host: "calendar.ru",
					Port: 5001,
				},
				Logger: Logger{
					Level: cnst.LoggerLevelError,
				},
				Storage: Storage{
					Type: cnst.StorageTypePostgres,
					Conn: "postgres://user:secret@localhost:5555",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for env, val := range tc.env {
				os.Setenv(env, val)
			}
			defer os.Clearenv()

			cfg := NewConfig(tc.path, tc.file)

			if !reflect.DeepEqual(cfg, tc.exp) {
				t.Errorf("wrong data %v, exp %v", cfg, tc.exp)
			}
		})
	}
}
