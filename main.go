package main

import (
	"github.com/ectrc/snow/aid"
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
		{
			user.CommonCoreProfile.Attributes.AddAttribute(person.NewAttribute("xp", 1030))
			user.CommonCoreProfile.Attributes.AddAttribute(person.NewAttribute("level", 100))
			user.CommonCoreProfile.Attributes.AddAttribute(person.NewAttribute("quest_manager", aid.JSON{}))

			user.CommonCoreProfile.Items.AddItem(person.NewItem("Currency:MtxPurchased", 100))
			user.CommonCoreProfile.Items.AddItem(person.NewItem("Token:CampaignAccess", 1))
		
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
				giftBox.AddLoot(person.NewItemWithType("AthenaCharacter:CID_002_Athena_Commando_F_Default", 1, "athena"))
			}
			user.CommonCoreProfile.Gifts.AddGift(giftBox)
		}
		user.Save()

		snapshot := user.CommonCoreProfile.Snapshot()
		{
			vbucks := user.CommonCoreProfile.Items.GetItemByTemplateID("Currency:MtxPurchased")
			vbucks.Quantity = 200
			vbucks.Favorite = true
			
			user.CommonCoreProfile.Items.DeleteItem(user.CommonCoreProfile.Items.GetItemByTemplateID("Token:CampaignAccess").ID)
			user.CommonCoreProfile.Items.AddItem(person.NewItem("Token:ReceiveMtxCurrency", 1))
		}	
		user.CommonCoreProfile.Diff(snapshot)
		user.Save()
	}

	go storage.Cache.CacheKiller()
}

func main() {
	users := person.AllFromDatabase()

	for _, user := range users {
		aid.PrintJSON(user.Snapshot())
	}

	// aid.WaitForExit()
}