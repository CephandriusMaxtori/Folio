package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port     int        `yaml:"port"`
	Database Database   `yaml:"database"`
	JWT      JWTConfig  `yaml:"jwt"`
	Library  LibraryCfg `yaml:"library"`
}

type Database struct {
	Path string `yaml:"path"`
}

type JWTConfig struct {
	Secret     string `yaml:"secret"`
	ExpireHour int    `yaml:"expire_hours"`
}

type LibraryCfg struct {
	Paths []string `yaml:"paths"`
}

func Load(path string) (*Config, error) {
	cfg := &Config{
		Port: 8080,
		Database: Database{
			Path: "./data/folio.db",
		},
		JWT: JWTConfig{
			Secret:     "change-me-in-production",
			ExpireHour: 72,
		},
	}

	if path == "" {
		paths := []string{"config.yaml", "config/config.yaml"}
		for _, p := range paths {
			if _, err := os.Stat(p); err == nil {
				path = p
				break
			}
		}
	}

	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}

	if v := os.Getenv("FOLIO_PORT"); v != "" {
		fmt.Sscanf(v, "%d", &cfg.Port)
	}
	if v := os.Getenv("FOLIO_DB_PATH"); v != "" {
		cfg.Database.Path = v
	}
	if v := os.Getenv("FOLIO_JWT_SECRET"); v != "" {
		cfg.JWT.Secret = v
	}

	dir := filepath.Dir(cfg.Database.Path)
	if dir != "" && dir != "." {
		os.MkdirAll(dir, 0755)
	}

	return cfg, nil
}
