package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Discord *DiscordConfig `yaml:"discord"`
}

// question: rename this to something else to be more representative of the configs it holds
func NewYamlConfig(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	err = cfg.Validate()
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) Validate() error {
	err := c.Discord.Validate()
	if err != nil {
		return fmt.Errorf("Discord: %v", err)
	}
	return nil
}
