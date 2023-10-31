package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ectrc/snow/config"
	"github.com/ectrc/snow/person"
	"github.com/ectrc/snow/storage"
)

const (
	DROP_TABLES = true
)

func init() {
	config := config.Get()

	var device storage.Storage

	switch config.Database.Type {
	case "postgres":
		postgresStorage := storage.NewPostgresStorage()

		if DROP_TABLES {
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
	if DROP_TABLES {
		user := person.NewPerson()
		snapshot := user.Snapshot()

		quest := person.NewQuest("Quest:Quest_1", "ChallengeBundle:Daily_1", "ChallengeBundleSchedule:Paid_1")
		{
			quest.AddObjective("quest_objective_eliminateplayers", 0)
			quest.AddObjective("quest_objective_top1", 0)
			quest.AddObjective("quest_objective_place_top10", 0)
			
			quest.UpdateObjectiveCount("quest_objective_eliminateplayers", 10)
			quest.UpdateObjectiveCount("quest_objective_place_top10", -3)
			
			quest.RemoveObjective("quest_objective_top1")
		}
		user.AthenaProfile.Quests.AddQuest(quest)

		giftBox := person.NewGift("GiftBox:GB_Default", 1, user.ID, "Hello, Bully!")
		{
			giftBox.AddLoot(person.NewItem("AthenaCharacter:CID_002_Athena_Commando_F_Default", 1))
		}
		user.CommonCoreProfile.Gifts.AddGift(giftBox)

		currency := person.NewItem("Currency:MtxPurchased", 100)
		user.CommonCoreProfile.Items.AddItem(currency)

		user.FindChanges(*snapshot)
		user.Save()
		printjson(user.Snapshot())
	}

	go storage.Cache.CacheKiller()
}

func main() {
	persons := person.AllFromDatabase()

	for _, person := range persons {
		fmt.Println(person)
	}

	wait()
}

func wait() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func printjson(v interface{}) {
	json1, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(json1))
}