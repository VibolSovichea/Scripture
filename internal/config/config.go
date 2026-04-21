package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	configDir  = ".scripture"
	configFile = "config.yaml"
)

type Config struct {
	VaultPath string `yaml:"vault_path"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDir, configFile), nil
}

func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func Init(vaultPath string) (*Config, error) {
	if len(vaultPath) >= 2 && vaultPath[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		vaultPath = filepath.Join(home, vaultPath[2:])
	}

	dirs := []string{
		vaultPath,
		filepath.Join(vaultPath, "notes"),
		filepath.Join(vaultPath, "daily"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return nil, err
		}
	}

	internal := filepath.Join(vaultPath, ".scripture")
	if err := os.MkdirAll(filepath.Join(internal, "trash"), 0755); err != nil {
		return nil, err
	}

	cfg := &Config{VaultPath: vaultPath}

	path, err := configPath()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return nil, err
	}

	return cfg, nil
}
