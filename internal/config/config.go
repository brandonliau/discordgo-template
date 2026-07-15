package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalYAML(node *yaml.Node) error {
	var raw string
	if err := node.Decode(&raw); err != nil {
		return fmt.Errorf("duration must be a string: %w", err)
	}
	parsed, err := time.ParseDuration(raw)
	if err != nil {
		return fmt.Errorf("parse duration %q: %w", raw, err)
	}
	d.Duration = parsed
	return nil
}

type Config struct {
	Discord      DiscordConfig      `yaml:"discord"`
	Database     DatabaseConfig     `yaml:"database"`
	SampleWorker SampleWorkerConfig `yaml:"sample_worker"`
}

type DiscordConfig struct {
	Token         string `yaml:"token"`
	ApplicationID string `yaml:"application_id"`
	GuildID       string `yaml:"guild_id"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type SampleWorkerConfig struct {
	Interval Duration `yaml:"interval"`
}

func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config: %w", err)
	}
	defer file.Close()

	cfg, err := Decode(file)
	if err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}
	return cfg, nil
}

func Decode(reader io.Reader) (*Config, error) {
	var cfg Config
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c Config) Validate() error {
	var errs []error
	required := []struct {
		name  string
		value string
	}{
		{"discord.token", c.Discord.Token},
		{"discord.application_id", c.Discord.ApplicationID},
		{"discord.guild_id", c.Discord.GuildID},
		{"database.path", c.Database.Path},
	}
	for _, field := range required {
		if strings.TrimSpace(field.value) == "" {
			errs = append(errs, fmt.Errorf("%s is required", field.name))
		}
	}
	if c.SampleWorker.Interval.Duration <= 0 {
		errs = append(errs, errors.New("sample_worker.interval must be greater than zero"))
	}
	return errors.Join(errs...)
}
