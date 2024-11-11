package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type yamlConfig struct {
	Token string
}

func NewYamlConfig(file string) *yamlConfig {
	cfg := &yamlConfig{}
	err := cfg.loadConfig(file)
	if err != nil {
		log.Fatalf("[FATAL] Failed to load config file: %v", err)
	}
	err = cfg.validateConfig()
	if err != nil {
		log.Fatalf("[FATAL] Failed to validate config file: %v", err)
	}

	return cfg
}

func (c *yamlConfig) loadConfig(file string) error {
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return err
	}
	return nil
}

func (c *yamlConfig) validateConfig() error {
	if c.Token == "" {
		return fmt.Errorf("empty token")
	}
	return nil
}
