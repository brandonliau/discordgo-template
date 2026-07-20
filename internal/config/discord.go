package config

import (
	"errors"
)

type DiscordConfig struct {
	Token         string `yaml:"token"`
	ApplicationID string `yaml:"application_id"`
	GuildID       string `yaml:"guild_id"`
}

func (c DiscordConfig) Validate() error {
	if c.Token == "" {
		return errors.New("token is required")
	}
	if c.ApplicationID == "" {
		return errors.New("application ID is required")
	}
	if c.GuildID == "" {
		return errors.New("guild ID is required")
	}
	return nil
}
