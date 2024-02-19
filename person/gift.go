package person

import (
	"time"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/storage"
	"github.com/google/uuid"
)

type Gift struct {
	ID string
	ProfileID  string
	TemplateID string
	Quantity int
	FromID string
	GiftedAt int64
	Message string
	Loot []*Item
}

func NewGift(templateID string, quantity int, fromID string, message string) *Gift {
	return &Gift{
		ID: uuid.New().String(),
		TemplateID: templateID,
		Quantity: quantity,
		FromID: fromID,
		GiftedAt: time.Now().Unix(),
		Message: message,
		Loot: []*Item{},
	}
}

func FromDatabaseGift(gift *storage.DB_Gift) *Gift {
	loot := []*Item{}

	for _, item := range gift.Loot {
		loot = append(loot, FromDatabaseGiftLoot(&item))
	}

	return &Gift{
		ID: gift.ID,
		ProfileID: gift.ProfileID,
		TemplateID: gift.TemplateID,
		Quantity: gift.Quantity,
		FromID: gift.FromID,
		GiftedAt: gift.GiftedAt,
		Message: gift.Message,
		Loot: loot,
	}
}

func (g *Gift) GenerateFortniteGiftEntry() aid.JSON {
	json := aid.JSON{
		"templateId": g.TemplateID,
		"attributes": aid.JSON{
			"params": aid.JSON{},
			"lootList": []aid.JSON{},
			"fromAccountId": g.FromID,
			"giftedOn": time.Unix(g.GiftedAt, 0).Format(time.RFC3339),
		},
		"quantity": 1,
	}

	for _, loot := range g.Loot {
		json["attributes"].(aid.JSON)["lootList"] = append(json["attributes"].(aid.JSON)["lootList"].([]aid.JSON), aid.JSON{
			"itemGuid": loot.ID,
			"itemType": loot.TemplateID,
			"itemProfile": loot.ProfileType,
			"quantity": loot.Quantity,
		})
	}

	if g.Message != "" {
		json["attributes"].(aid.JSON)["params"].(aid.JSON)["userMessage"] = g.Message
	}

	return json
}

func (g *Gift) AddLoot(loot *Item) {
	g.Loot = append(g.Loot, loot)
	//storage.Repo.SaveGiftLoot(g.ID, loot)
}

func (g *Gift) FillLoot(loot []*Item) {
	g.Loot = loot
}

func (g *Gift) Delete() {
	for _, item := range g.Loot {
		item.DeleteLoot()
	}

	storage.Repo.DeleteGift(g.ID)
}

func (g *Gift) ToDatabase(profileId string) *storage.DB_Gift {
	profileLoot := []storage.DB_GiftLoot{}

	for _, item := range g.Loot {
		profileLoot = append(profileLoot, *item.ToGiftLootDatabase(g.ID))
	}

	return &storage.DB_Gift{
		ID: g.ID,
		ProfileID: profileId,
		TemplateID: g.TemplateID,
		Quantity: g.Quantity,
		FromID: g.FromID,
		GiftedAt: g.GiftedAt,
		Message: g.Message,
		Loot: profileLoot,
	}
}

func (g *Gift) Save() {
	if g.ProfileID == "" {
		return
	}

	storage.Repo.SaveGift(g.ToDatabase(g.ProfileID))
	for _, item := range g.Loot {
		item.SaveLoot(g.ID)
	}
}

func (g *Gift) Snapshot() GiftSnapshot {
	loot := []Item{}

	for _, item := range g.Loot {
		loot = append(loot, *item)
	}

	return GiftSnapshot{
		ID: g.ID,
		TemplateID: g.TemplateID,
		Quantity: g.Quantity,
		FromID: g.FromID,
		GiftedAt: g.GiftedAt,
		Message: g.Message,
		Loot: loot,
	}
}