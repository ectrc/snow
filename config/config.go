package config

import (
	"os"

	"gopkg.in/ini.v1"
)

var (
	local *Config
)

type Config struct {
	Database struct {
		URI  string
		Type string
	}
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load() {
	configPath := "config.ini"
	if _, err := os.Stat(configPath); err != nil {
		panic("config.ini not found! please rename default.config.ini to config.ini and complete")
	}

	cfg, err := ini.Load("config.ini")
	if err != nil {
		panic(err)
	}

	c.Database.URI = cfg.Section("database").Key("uri").String()
	if c.Database.URI == "" {
		panic("Database URI is empty")
	}
	c.Database.Type = cfg.Section("database").Key("type").String()
	if c.Database.Type == "" {
		panic("Database Type is empty")
	}
}

func Get() *Config {
	if local == nil {
		local = NewConfig()
		local.Load()
	}

	return local
}