// Package config provides application configuration loading and structures.
package config

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config represents the application configuration.
type Config struct {
	Addr   string `mapstructure:"addr"`   // listen address, e.g. "0.0.0.0:9000"
	Driver string `mapstructure:"driver"` // database driver: "postgres" (default) or "mysql"
	DSN    string `mapstructure:"dsn"`    // database connection string
}

// Load reads configuration from the given file path.
func Load(path string) (*Config, error) {
	cfgFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(bytes.NewReader(cfgFile)); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if cfg.Addr == "" {
		return nil, fmt.Errorf("addr is required")
	}
	if cfg.DSN == "" {
		return nil, fmt.Errorf("dsn is required")
	}
	if cfg.Driver == "" {
		cfg.Driver = "postgres"
	}
	if cfg.Driver != "postgres" && cfg.Driver != "mysql" {
		return nil, fmt.Errorf("unsupported driver: %s (supported: postgres, mysql)", cfg.Driver)
	}

	return &cfg, nil
}
