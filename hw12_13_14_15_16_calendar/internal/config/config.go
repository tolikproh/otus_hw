package config

import (
	"net/url"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/viper"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/internal/cnst"
)

type Config struct {
	HTTPServer `yaml:"httpserver"`
	Logger     `yaml:"logger"`
	Storage    `yaml:"storage"`
}

type HTTPServer struct {
	Host string `yaml:"host" env:"CALENDAR_HTTP_SERVER_HOST"`
	Port int    `yaml:"port" env:"CALENDAR_HTTP_SERVER_PORT"`
}

type Logger struct {
	Level string `yaml:"level" env:"CALENDAR_LOG_LEVEL"`
}

type Storage struct {
	Type string `yaml:"type" env:"CALENDAR_STORAGE_TYPE"` // MEM or SQL
	Conn string `yaml:"conn" env:"CALENDAR_STORAGE_CONN"`
}

func defaultConfig() *Config {
	var cfg Config
	// Http Server
	cfg.HTTPServer.Host = "example.com"
	cfg.HTTPServer.Port = 8080

	// Logger
	cfg.Logger.Level = cnst.LoggerLevelInfo

	// Storage
	cfg.Storage.Type = cnst.StorageTypeMemory
	cfg.Storage.Conn = ""

	return &cfg
}

// NewConfig Set Default.
func NewConfig(path, file string) *Config {
	cfg := defaultConfig()

	viper.SetConfigName(file)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err == nil {
		viper.Unmarshal(&cfg)
	}

	cleanenv.ReadEnv(cfg)

	if !cfg.checkStorageType() {
		cfg.Storage.Type = cnst.StorageTypeMemory
		cfg.Storage.Conn = ""
	}

	return cfg
}

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// return (true = ok type), (false = type not constant).
func (cfg *Config) checkStorageType() bool {
	switch cfg.Storage.Type {
	case cnst.StorageTypePostgres:
		if !isURL(cfg.Storage.Conn) {
			return false
		}
		return true
	case cnst.StorageTypeMemory:
		return true

	default:
		return false
	}
}
