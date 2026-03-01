package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Sources []string `yaml:"sources"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	for i, src := range cfg.Sources {
		abs, err := filepath.Abs(src)
		if err != nil {
			return nil, fmt.Errorf("source[%d] invalid path %q: %w", i, src, err)
		}
		info, err := os.Stat(abs)
		if err != nil || !info.IsDir() {
			return nil, fmt.Errorf("source[%d] not a directory: %q", i, src)
		}
		cfg.Sources[i] = abs
	}

	return &cfg, nil
}
