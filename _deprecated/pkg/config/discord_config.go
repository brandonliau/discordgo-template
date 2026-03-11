package config

import (
	"fmt"

	"discord-template/pkg/logger"
)

type DiscordConfig struct {
	Token          string `yaml:"token"`
	Guild          string `yaml:"guild"`
	AuthorizedRole string `yaml:"authorized_role"`
	logger         logger.Logger
}

func NewDiscordConfig(file string, logger logger.Logger) *DiscordConfig {
	cfg := &DiscordConfig{
		logger: logger,
	}
	err := load(file, cfg)
	if err != nil {
		logger.Fatal("Failed to load config file: %v", err)
	}
	err = cfg.validate()
	if err != nil {
		logger.Fatal("Failed to validate config file: %v", err)
	}
	return cfg
}

func (c *DiscordConfig) validate() error {
	if c.Token == "" {
		return fmt.Errorf("empty token")
	}
	if c.Guild == "" {
		return fmt.Errorf("empty guild")
	}
	if c.AuthorizedRole == "" {
		return fmt.Errorf("empty authorized role")
	}
	return nil
}
