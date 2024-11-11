package config

type Config interface {
	loadConfig(string) error
	validateConfig() error
}
