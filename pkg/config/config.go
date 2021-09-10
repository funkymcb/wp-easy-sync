package config

import (
	"os"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Easyverein EasyvereinCfg
}

var config *Config

// LoadConfig from ConfigPath
func LoadConfig(configPath string) error {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(file, &config); err != nil {
		return err
	}

	return nil
}

// GetConfig returns the actual config
func GetConfig() *Config {
	return config
}

// UnmarshalYAML overrides the default implementation with the ability to set default
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	defaults.Set(c)

	type plain Config
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}

	return nil
}
