package config

import (
	"discordgo-skeleton/pkg/utils"
)

type Config struct {
	Discord DiscordConfig `yaml:"discord"`
}

func Load(path string) (*Config, error) {
	var cfg Config

	err := utils.DecodeYAML(path, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c Config) Validate() error {
	if err := c.Discord.Validate(); err != nil {
		return err
	}
	return nil
}
