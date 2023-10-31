package person

import (
	"time"

	"github.com/ectrc/snow/storage"
	"github.com/google/uuid"
)

type Gift struct {
	ID         string
	TemplateID string
	Quantity   int
	FromID     string
	GiftedAt   int64
	Message		 string
	Loot			 []*Item
}

func NewGift(templateID string, quantity int, fromID string, message string) *Gift {
	return &Gift{
		ID:         uuid.New().String(),
		TemplateID: templateID,
		Quantity:   quantity,
		FromID:     fromID,
		GiftedAt:   time.Now().Unix(),
		Message:		message,
		Loot:				[]*Item{},
	}
}

func FromDatabaseGift(gift *storage.DB_Gift) *Gift {
	loot := []*Item{}

	for _, item := range gift.Loot {
		loot = append(loot, FromDatabaseLoot(&item))
	}

	return &Gift{
		ID:         gift.ID,
		TemplateID: gift.TemplateID,
		Quantity:   gift.Quantity,
		FromID:     gift.FromID,
		GiftedAt:   gift.GiftedAt,
		Message:		gift.Message,
		Loot:				loot,
	}
}

func (g *Gift) AddLoot(loot *Item) {
	g.Loot = append(g.Loot, loot)
	//storage.Repo.SaveGiftLoot(g.ID, loot)
}

func (g *Gift) FillLoot(loot []*Item) {
	g.Loot = loot
}

func (g *Gift) Delete() {
	g.Quantity = 0
	g.Loot = []*Item{}
	//storage.Repo.DeleteGift(g.ID)
}

func (g *Gift) ToDatabase(profileId string) *storage.DB_Gift {
	profileLoot := []storage.DB_Loot{}

	for _, item := range g.Loot {
		profileLoot = append(profileLoot, *item.ToLootDatabase(g.ID))
	}

	return &storage.DB_Gift{
		ProfileID: profileId,
		ID:         g.ID,
		TemplateID: g.TemplateID,
		Quantity:   g.Quantity,
		FromID:     g.FromID,
		GiftedAt:   g.GiftedAt,
		Message:		g.Message,
		Loot:				profileLoot,
	}
}

func (g *Gift) Save() {
	//storage.Repo.SaveGift(g.ToDatabase())
}

func (g *Gift) Snapshot() GiftSnapshot {
	loot := []Item{}

	for _, item := range g.Loot {
		loot = append(loot, *item)
	}

	return GiftSnapshot{
		ID:         g.ID,
		TemplateID: g.TemplateID,
		Quantity:   g.Quantity,
		FromID:     g.FromID,
		GiftedAt:   g.GiftedAt,
		Message:		g.Message,
		Loot:				loot,
	}
}