package main

import (
	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/handlers"
	"github.com/ectrc/snow/person"
	"github.com/ectrc/snow/storage"

	"github.com/gofiber/fiber/v2"
)
func init() {
	aid.LoadConfig()

	var device storage.Storage
	switch aid.Config.Database.Type {
	case "postgres":
		postgresStorage := storage.NewPostgresStorage()

		if aid.Config.Database.DropAllTables {
			aid.Print("Dropping all tables")
			postgresStorage.DropTables()
		}

		postgresStorage.Migrate(&storage.DB_Person{}, "Persons")
		postgresStorage.Migrate(&storage.DB_Profile{}, "Profiles")
		postgresStorage.Migrate(&storage.DB_Item{}, "Items")
		postgresStorage.Migrate(&storage.DB_Gift{}, "Gifts")
		postgresStorage.Migrate(&storage.DB_Quest{}, "Quests")
		postgresStorage.Migrate(&storage.DB_Loot{}, "Loot")
		postgresStorage.Migrate(&storage.DB_VariantChannel{}, "Variants")
		postgresStorage.Migrate(&storage.DB_PAttribute{}, "Attributes")

		device = postgresStorage
	}

	storage.Repo = storage.NewStorage(device)
}

func init() {
	if aid.Config.Database.DropAllTables {
		person.NewFortnitePerson("ac", "1")
	}

	aid.PrintTime("Loading all persons from database", func() {
		for _, person := range person.AllFromDatabase() {
			aid.Print("Loaded person: " + person.DisplayName)
		}
	})

	aid.PrintTime("Loading all persons from cache", func() {
		for _, person := range person.AllFromCache() {
			aid.Print("Loaded person: " + person.DisplayName)
		}
	})
}

func main() {
	r := fiber.New()

	r.Use(aid.FiberLogger())
	r.Use(aid.FiberLimiter())
	r.Use(aid.FiberCors())

	r.Get("/content/api/pages/fortnite-game", handlers.GetContentPages)

	account := r.Group("/account/api")
	account.Get("/public/account/:accountId", handlers.GetPublicAccount)
	account.Get("/public/account/:accountId/externalAuths", handlers.GetPublicAccountExternalAuths)
	account.Get("/oauth/verify", handlers.GetOAuthVerify)
	account.Post("/oauth/token", handlers.PostOAuthToken)
	account.Delete("/oauth/sessions/kill", handlers.DeleteOAuthSessions)

	fortnite := r.Group("/fortnite/api")
	fortnite.Get("/receipts/v1/account/:accountId/receipts", handlers.GetFortniteReceipts)
	fortnite.Get("/v2/versioncheck/*", handlers.GetFortniteVersion)
	fortnite.Get("/calendar/v1/timeline", handlers.GetFortniteTimeline)

	storefront := fortnite.Group("/storefront/v2")
	storefront.Get("/catalog", handlers.GetStorefrontCatalog)
	storefront.Get("/keychain", handlers.GetStorefrontKeychain)

	matchmaking := fortnite.Group("/matchmaking")
	matchmaking.Get("/session/findPlayer/:accountId", handlers.GetMatchmakingSession)

	storage := fortnite.Group("/cloudstorage")
	storage.Get("/system", handlers.GetCloudStorageFiles)
	storage.Get("/system/config", handlers.GetCloudStorageConfig)
	storage.Get("/system/:fileName", handlers.GetCloudStorageFile)
	storage.Get("/user/:accountId", handlers.GetUserStorageFiles)
	storage.Get("/user/:accountId/:fileName", handlers.GetUserStorageFile)
	storage.Put("/user/:accountId/:fileName", handlers.PutUserStorageFile)
	
	game := fortnite.Group("/game/v2")
	game.Get("/enabled_features", handlers.GetGameEnabledFeatures)
	game.Post("/tryPlayOnPlatform/account/:accountId", handlers.PostGamePlatform)
	game.Post("/grant_access/:accountId", handlers.PostGameAccess)

	profile := game.Group("/profile/:accountId")
	profile.Post("/client/:action", handlers.PostProfileAction)

	lightswitch := r.Group("/lightswitch/api")
	lightswitch.Get("/service/bulk/status", handlers.GetLightswitchBulkStatus)

	r.All("*", func(c *fiber.Ctx) error { return c.Status(fiber.StatusNotFound).JSON(aid.ErrorNotFound) })
	r.Listen(aid.Config.API.Host + aid.Config.API.Port)
}