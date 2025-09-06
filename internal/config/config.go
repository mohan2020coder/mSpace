// internal/config/config.go
package config

import (
	"time"

	"github.com/spf13/viper"
)

type ServerCfg struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type DatabaseCfg struct {
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
}

type StorageCfg struct {
	Provider  string `mapstructure:"provider"`
	Endpoint  string `mapstructure:"endpoint"`
	Bucket    string `mapstructure:"bucket"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	SSL       bool   `mapstructure:"ssl"`
}

type AuthCfg struct {
	JWTSecret   string        `mapstructure:"jwt_secret"`
	TokenExpiry time.Duration `mapstructure:"token_expiry"`
}

type LoggingCfg struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type Config struct {
	Server   ServerCfg   `mapstructure:"server"`
	Database DatabaseCfg `mapstructure:"database"`
	Storage  StorageCfg  `mapstructure:"storage"`
	Auth     AuthCfg     `mapstructure:"auth"`
	Logging  LoggingCfg  `mapstructure:"logging"`
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Ensure token expiry parsed if provided as string like "24h"
	if val := v.GetString("auth.token_expiry"); val != "" {
		d, err := time.ParseDuration(val)
		if err == nil {
			cfg.Auth.TokenExpiry = d
		}
	}

	return &cfg, nil
}
