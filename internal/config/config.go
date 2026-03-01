package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	PollingInterval time.Duration   `mapstructure:"polling_interval"`
	Ntfy            NtfyConfig      `mapstructure:"ntfy"`
	Packages        []PackageConfig `mapstructure:"packages"`
	CountryCode     string          `mapstructure:"country_code"`
	Log             LogConfig       `mapstructure:"log"`
	DB              DBConfig        `mapstructure:"db"`
}

type NtfyConfig struct {
	URL   string `mapstructure:"url"`
	Topic string `mapstructure:"topic"`
	Token string `mapstructure:"token"`
}

type PackageConfig struct {
	ID   int    `mapstructure:"id"`
	Name string `mapstructure:"name"`
}

type LogConfig struct {
	Path       string `mapstructure:"path"`
	MaxSizeMB  int    `mapstructure:"max_size_mb"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type DBConfig struct {
	Path string `mapstructure:"path"`
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	return &cfg, nil
}
