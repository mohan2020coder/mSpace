// internal/config/config.go
package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host string
		Port int
	}
	Database struct {
		Driver string
		DSN    string
	}
	Storage struct {
		Provider  string
		Endpoint  string
		Bucket    string
		AccessKey string
		SecretKey string
		SSL       bool
	}
	Search struct {
		Provider  string
		IndexPath string
	}
	Auth struct {
		JWTSecret   string
		TokenExpiry time.Duration
	}
	Logging struct {
		Level  string
		Format string
	}
}

func LoadConfig(path string) *Config {
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("❌ Config error: %s", err)
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("❌ Config unmarshal error: %s", err)
	}
	return &cfg
}
