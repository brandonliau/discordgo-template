package config

import (
	"fmt"
)

type DiscordConfig struct {
	Token         string `yaml:"token"`
	ApplicationID string `yaml:"application_id"`
	GuildID       string `yaml:"guild_id"`
}

func (c *DiscordConfig) Validate() error {
	if c.Token == "" {
		return fmt.Errorf("token is required")
	}
	if c.ApplicationID == "" {
		return fmt.Errorf("application ID is required")
	}
	if c.GuildID == "" {
		return fmt.Errorf("guild ID is required")
	}
	return nil
}
