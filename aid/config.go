package aid

import (
	"os"

	"gopkg.in/ini.v1"
)

type CS struct {
	Database struct {
		URI  string
		Type string
		DropAllTables bool
	}
	Output struct {
		Level string
	}
	API struct {
		Host string
		Port string
	}
	JWT struct {
		Secret string
	}
}

var (
	Config *CS
)

func LoadConfig() {
	Config = &CS{}

	configPath := "config.ini"
	if _, err := os.Stat(configPath); err != nil {
		panic("config.ini not found! please rename default.config.ini to config.ini and complete")
	}

	cfg, err := ini.Load("config.ini")
	if err != nil {
		panic(err)
	}

	Config.Database.DropAllTables = cfg.Section("database").Key("drop").MustBool(false)
	Config.Database.URI = cfg.Section("database").Key("uri").String()
	if Config.Database.URI == "" {
		panic("Database URI is empty")
	}
	Config.Database.Type = cfg.Section("database").Key("type").String()
	if Config.Database.Type == "" {
		panic("Database Type is empty")
	}

	Config.Output.Level = cfg.Section("output").Key("level").String()
	if Config.Output.Level == "" {
		panic("Output Level is empty")
	}

	if Config.Output.Level != "dev" && Config.Output.Level != "prod" && Config.Output.Level != "time" && Config.Output.Level != "info" {
		panic("Output Level must be either dev or prod")
	}

	Config.API.Host = cfg.Section("api").Key("host").String()
	if Config.API.Host == "" {
		panic("API Host is empty")
	}

	Config.API.Port = cfg.Section("api").Key("port").String()
	if Config.API.Port == "" {
		panic("API Port is empty")
	}

	Config.JWT.Secret = cfg.Section("jwt").Key("secret").String()
	if Config.JWT.Secret == "" {
		panic("JWT Secret is empty")
	}
}