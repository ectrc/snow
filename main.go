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
		postgresStorage.Migrate(&storage.DB_Loadout{}, "Loadouts")
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
	storage.Cache = storage.NewPersonsCacheMutex()
}

func init() {
	if aid.Config.Database.DropAllTables {
		person.NewFortnitePerson("ac", "ket")
	}

	aid.PrintTime("Loading all persons from database", func() {
		for _, person := range person.AllFromDatabase() {
			aid.Print("Loaded person: " + person.DisplayName)
		}
	})

	go storage.Cache.CacheKiller()
}

func main() {
	r := fiber.New()

	r.Use(aid.FiberLogger())
	r.Use(aid.FiberLimiter())
	r.Use(aid.FiberCors())

	account := r.Group("/account/api")
	account.Get("/public/account/:accountId", handlers.GetPublicAccount)
	account.Get("/public/account/:accountId/externalAuths", handlers.GetPublicAccountExternalAuths)
	account.Post("/oauth/token", handlers.PostOAuthToken)
	account.Delete("/oauth/sessions/kill", handlers.DeleteOAuthSessions)

	fortnite := r.Group("/fortnite/api")
	fortnite.Get("/receipts/v1/account/:accountId/receipts", handlers.GetAccountReceipts)
	fortnite.Get("/versioncheck/*", handlers.GetVersionCheck)

	matchmaking := fortnite.Group("/matchmaking")
	matchmaking.Get("/session/findPlayer/:accountId", handlers.GetSessionFindPlayer)

	storage := fortnite.Group("/cloudstorage")
	storage.Get("/system", handlers.GetCloudStorageFiles)
	storage.Get("/system/config", handlers.GetCloudStorageConfig)
	storage.Get("/system/:fileName", handlers.GetCloudStorageFile)
	
	game := fortnite.Group("/game/v2")
	game.Post("/tryPlayOnPlatform/account/:accountId", handlers.PostTryPlayOnPlatform)
	game.Post("/grant_access/:accountId", handlers.PostGrantAccess)
	game.Get("/enabled_features", handlers.GetEnabledFeatures)

	profile := game.Group("/profile/:accountId")
	profile.Post("/client/:action", handlers.PostProfileAction)

	lightswitch := r.Group("/lightswitch/api")
	lightswitch.Get("/service/bulk/status", handlers.GetLightswitchBulkStatus)

	r.All("*", func(c *fiber.Ctx) error { return c.Status(fiber.StatusNotFound).JSON(aid.ErrorNotFound) })
	r.Listen(aid.Config.API.Host + aid.Config.API.Port)
}