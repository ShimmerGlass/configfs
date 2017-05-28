package main

import (
	"github.com/BurntSushi/toml"
)

type ConfigFn func() (Config, error)

type Config struct {
	Env         string            `toml:"env"`
	EnvPatterns map[string]string `toml:"env_patterns"`
}

type ConfigManager struct {
	path string
}

func NewConfigManager(path string) *ConfigManager {
	return &ConfigManager{
		path: path,
	}
}

func (m *ConfigManager) Load() (Config, error) {
	c := Config{}
	_, err := toml.DecodeFile(m.path, &c)
	return c, err
}
