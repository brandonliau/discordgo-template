package config

import (
	"fmt"
	"os"

	"discord-template/pkg/logger"

	"gopkg.in/yaml.v3"
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
	err := cfg.load(file)
	if err != nil {
		logger.Fatal("Failed to load config file: %v", err)
	}
	err = cfg.validate()
	if err != nil {
		logger.Fatal("Failed to validate config file: %v", err)
	}
	return cfg
}

func (c *DiscordConfig) load(file string) error {
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("readfile: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return fmt.Errorf("unmarshal: %v", err)
	}
	return nil
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
