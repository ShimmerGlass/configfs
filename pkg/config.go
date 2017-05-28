package configfs

import (
	"github.com/BurntSushi/toml"
	"os"
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

func (m *ConfigManager) Save(cfg Config) error {
	f, err := os.Create(m.path)
	if err != nil {
		return err
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(cfg)
}
