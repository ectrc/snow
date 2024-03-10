package aid

import (
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
)

type CS struct {
	Accounts struct {
		Gods []string
		Owners []string
	}
	Database struct {
		URI  string
		Type string
		DropAllTables bool
	}
	Discord struct {
		ID string
		Secret string
		Token string
		Guild string
		CallbackURL string
	}
	Amazon struct {
		Enabled bool
		BucketURI string
		AccessKeyID string
		SecretAccessKey string
		ClientSettingsBucket string
	}
	Output struct {
		Level string
	}
	API struct {
		Host string
		Port string
		FrontendPort string
		Debug bool
		XMPP struct {
			Host string
			Port string
		}
	}
	JWT struct {
		Secret string
	}
	Fortnite struct {
		Season int
		Build float64
		Everything bool
		Password bool
		DisableClientCredentials bool
		ShopSeed int
		EnableVBucks bool
	}
}

var (
	Config *CS
)

func LoadConfig(file []byte) {
	Config = &CS{}
	
	cfg, err := ini.Load(file)
	if err != nil {
		panic(err)
	}

	Config.Accounts.Gods = cfg.Section("accounts").Key("gods").Strings(",")
	Config.Accounts.Owners = cfg.Section("accounts").Key("owners").Strings(",")
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

	Config.Discord.ID = cfg.Section("discord").Key("id").String()
	if Config.Discord.ID == "" {
		panic("Discord Client ID is empty")
	}

	Config.Discord.Secret = cfg.Section("discord").Key("secret").String()
	if Config.Discord.Secret == "" {
		panic("Discord Client Secret is empty")
	}

	Config.Discord.Token = cfg.Section("discord").Key("token").String()
	if Config.Discord.Token == "" {
		panic("Discord Bot Token is empty")
	}

	Config.Discord.Guild = cfg.Section("discord").Key("guild").String()
	if Config.Discord.Guild == "" {
		panic("Discord Guild ID is empty")
	}

	Config.Discord.CallbackURL = cfg.Section("api").Key("discord_url").String()
	if Config.Discord.CallbackURL == "" {
		panic("Discord Callback URL is empty")
	}

	Config.Amazon.Enabled = true
	Config.Amazon.BucketURI = cfg.Section("amazon").Key("uri").String()
	if Config.Amazon.BucketURI == "" {
		Config.Amazon.Enabled = false
	}

	Config.Amazon.AccessKeyID = cfg.Section("amazon").Key("id").String()
	if Config.Amazon.AccessKeyID == "" {
		Config.Amazon.Enabled = false
	}

	Config.Amazon.SecretAccessKey = cfg.Section("amazon").Key("key").String()
	if Config.Amazon.SecretAccessKey == "" {
		Config.Amazon.Enabled = false
	}

	Config.Amazon.ClientSettingsBucket = cfg.Section("amazon").Key("bucket").String()
	if Config.Amazon.ClientSettingsBucket == "" {
		Config.Amazon.Enabled = false
	}

	Config.API.Host = cfg.Section("api").Key("host").String()
	if Config.API.Host == "" {
		panic("API Host is empty")
	}

	Config.API.Port = cfg.Section("api").Key("port").String()
	if Config.API.Port == "" {
		panic("API Port is empty")
	}

	Config.API.FrontendPort = cfg.Section("api").Key("frontend_port").String()
	if Config.API.FrontendPort == "" {
		Config.API.FrontendPort = Config.API.Port
	}

	Config.API.Debug = cfg.Section("api").Key("debug").MustBool(false)

	Config.API.XMPP.Host = cfg.Section("api").Key("xmpp_host").String()
	if Config.API.XMPP.Host == "" {
		panic("API XMPP Host is empty")
	}

	Config.API.XMPP.Port = cfg.Section("api").Key("xmpp_port").String()
	if Config.API.XMPP.Port == "" {
		panic("API XMPP Port is empty")
	}

	Config.JWT.Secret = cfg.Section("jwt").Key("secret").String()
	if Config.JWT.Secret == "" {
		panic("JWT Secret is empty")
	}

	build, err := cfg.Section("fortnite").Key("build").Float64()
	if err != nil {
		panic("Fortnite Build is empty")
	}

	Config.Fortnite.Build = build

	buildStr := strconv.FormatFloat(build, 'f', -1, 64)
	if buildStr == "" {
		panic("Fortnite Build is empty")
	}

	buildInfo := strings.Split(buildStr, ".")
	if len(buildInfo) < 2 {
		panic("Fortnite Build is invalid")
	}

	parsedSeason, err := strconv.Atoi(buildInfo[0])
	if err != nil {
		panic("Fortnite Season is invalid")
	}

	Config.Fortnite.Season = parsedSeason
	Config.Fortnite.Everything = cfg.Section("fortnite").Key("everything").MustBool(false)
	Config.Fortnite.Password = !(cfg.Section("fortnite").Key("disable_password").MustBool(false))
	Config.Fortnite.DisableClientCredentials = cfg.Section("fortnite").Key("disable_client_credentials").MustBool(false)
	Config.Fortnite.ShopSeed = cfg.Section("fortnite").Key("shop_seed").MustInt(0)
	Config.Fortnite.EnableVBucks = cfg.Section("fortnite").Key("enable_vbucks").MustBool(false)
}