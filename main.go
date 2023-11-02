package main

import (
	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/person"
	"github.com/ectrc/snow/storage"
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
	go storage.Cache.CacheKiller()
}

func main() {
	aid.PrintTime("Fetching Persons", func() {
		users := person.AllFromDatabase()
		aid.Print("Found", len(users), "users")
		for _, user := range users {
			aid.Print(user.ID)
		}
	})

	// aid.WaitForExit()
}