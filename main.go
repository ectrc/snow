package main

import (
	_ "embed"
	"fmt"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/discord"
	"github.com/ectrc/snow/fortnite"
	"github.com/ectrc/snow/handlers"
	"github.com/ectrc/snow/person"
	"github.com/ectrc/snow/storage"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

//go:embed config.ini
var configFile []byte

func init() {
	aid.LoadConfig(configFile)
	
	var device storage.Storage
	switch aid.Config.Database.Type {
	case "postgres":
		postgresStorage := storage.NewPostgresStorage()
		if aid.Config.Database.DropAllTables {
			aid.Print("Dropping all tables")
			postgresStorage.DropTables()
		}
		postgresStorage.MigrateAll()
		device = postgresStorage
	default:
		panic("Invalid database type: " + aid.Config.Database.Type)
	}

	storage.Repo = storage.NewStorage(device)
}

func init() {
	discord.IntialiseClient()
	fortnite.PreloadCosmetics(aid.Config.Fortnite.Season)
	fortnite.GenerateRandomStorefront()
	fortnite.GeneratePlaylistImages()

	if found := person.FindByDisplay("god"); found == nil {
		fortnite.NewFortnitePerson("god", true)
	}
}
func main() {
	r := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	r.Use(aid.FiberLogger())
	r.Use(aid.FiberLimiter())
	r.Use(aid.FiberCors())

	r.Get("/content/api/pages/fortnite-game", handlers.GetContentPages)
	r.Get("/waitingroom/api/waitingroom", handlers.GetWaitingRoomStatus)
	r.Get("/region", handlers.GetRegion)
	r.Put("/profile/play_region", handlers.AnyNoContent)
	r.Post("/api/v1/assets/Fortnite/:versionId/:assetName", handlers.PostAssets)

	account := r.Group("/account/api")
	account.Get("/public/account", handlers.GetPublicAccounts)
	account.Get("/public/account/:accountId", handlers.GetPublicAccount)
	account.Get("/public/account/:accountId/externalAuths", handlers.GetPublicAccountExternalAuths)
	account.Get("/public/account/displayName/:displayName", handlers.GetPublicAccountByDisplayName)
	account.Get("/oauth/verify", handlers.GetTokenVerify)
	account.Post("/oauth/token", handlers.PostFortniteToken)
	account.Delete("/oauth/sessions/kill", handlers.DeleteToken)

	fortnite := r.Group("/fortnite/api")
	fortnite.Get("/receipts/v1/account/:accountId/receipts", handlers.GetFortniteReceipts)
	fortnite.Get("/v2/versioncheck/:version", handlers.GetFortniteVersion)
	fortnite.Get("/calendar/v1/timeline", handlers.GetFortniteTimeline)

	storefront := fortnite.Group("/storefront/v2")
	storefront.Use(handlers.MiddlewareFortnite)
	storefront.Get("/catalog", handlers.GetStorefrontCatalog)
	storefront.Get("/keychain", handlers.GetStorefrontKeychain)

	matchmaking := fortnite.Group("/matchmaking")
	matchmaking.Get("/session/findPlayer/:accountId", handlers.GetMatchmakingSession)

	storage := fortnite.Group("/cloudstorage")
	storage.Get("/system", handlers.GetCloudStorageFiles)
	storage.Get("/system/config", handlers.GetCloudStorageConfig)
	storage.Get("/system/:fileName", handlers.GetCloudStorageFile)

	user := storage.Group("/user")
	user.Use(handlers.MiddlewareFortnite)
	user.Get("/:accountId", handlers.GetUserStorageFiles)
	user.Get("/:accountId/:fileName", handlers.GetUserStorageFile)
	user.Put("/:accountId/:fileName", handlers.PutUserStorageFile)
	
	game := fortnite.Group("/game/v2")
	game.Get("/enabled_features", handlers.GetGameEnabledFeatures)
	game.Post("/tryPlayOnPlatform/account/:accountId", handlers.PostGamePlatform)
	game.Post("/grant_access/:accountId", handlers.PostGameAccess)
	game.Post("/creative/discovery/surface/:accountId", handlers.PostDiscovery)
	game.Post("/profileToken/verify/:accountId", handlers.AnyNoContent)

	profile := game.Group("/profile/:accountId")
	profile.Use(handlers.MiddlewareFortnite)
	profile.Post("/client/:action", handlers.PostProfileAction)
	profile.Post("/dedicated_server/:action", handlers.PostProfileAction)

	lightswitch := r.Group("/lightswitch/api")
	lightswitch.Use(handlers.MiddlewareFortnite)
	lightswitch.Get("/service/bulk/status", handlers.GetLightswitchBulkStatus)

	snow := r.Group("/snow")
	snow.Get("/cosmetics", handlers.GetPreloadedCosmetics)
	snow.Get("/image/:playlist", handlers.GetPlaylistImage)

	discord := snow.Group("/discord")
	discord.Get("/", handlers.GetDiscordOAuthURL)
	
	player := snow.Group("/player")
	player.Use(handlers.MiddlewareWeb)
	player.Get("/", handlers.GetPlayer)
	player.Get("/locker", handlers.GetPlayerLocker)

	r.Hooks().OnListen(func(ld fiber.ListenData) error {
		aid.Print("Listening on " + aid.Config.API.Host + ":" + ld.Port)
		return nil
	})
	
	r.All("*", func(c *fiber.Ctx) error { return c.Status(fiber.StatusNotFound).JSON(aid.ErrorNotFound) })
	
	err := r.Listen("0.0.0.0" + aid.Config.API.Port)
	if err != nil {
		panic(fmt.Sprintf("Failed to listen: %v", err))
	}
}
