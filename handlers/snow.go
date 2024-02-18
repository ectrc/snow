package handlers

import (
	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/fortnite"
	p "github.com/ectrc/snow/person"
	"github.com/ectrc/snow/storage"
	"github.com/gofiber/fiber/v2"
)

func MiddlewareOnlyDebug(c *fiber.Ctx) error {
	if aid.Config.API.Debug {
		return c.Next()
	}

	return c.SendStatus(403)
}

func GetSnowPreloadedCosmetics(c *fiber.Ctx) error {
	return c.JSON(fortnite.External)
}

func GetPlayerLocker(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)

	items := make([]p.Item, 0)
	person.AthenaProfile.Items.RangeItems(func(key string, value *p.Item) bool {
		items = append(items, *value)
		return true
	})

	return c.JSON(items)
}

func GetPlayer(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)

	return c.JSON(aid.JSON{
		"id": person.ID,
		"displayName": person.DisplayName,
		"discord": person.Discord,
	})
}

func GetSnowCachedPlayers(c *fiber.Ctx) error {
	persons := p.AllFromCache()
	players := make([]p.PersonSnapshot, len(persons))

	for i, person := range persons {
		players[i] = *person.Snapshot()
	}

	return c.JSON(players)
}

func GetSnowConfig(c *fiber.Ctx) error {
	return c.JSON(aid.JSON{
		"basic": aid.Config,		
		"amazon": storage.Repo.Amazon,
	})
}

func GetSnowParties(c *fiber.Ctx) error {
	parties := []aid.JSON{}

	p.Parties.Range(func(key string, value *p.Party) bool {
		parties = append(parties, value.GenerateFortniteParty())
		return true
	})

	return c.JSON(parties)
}

func GetSnowShop(c *fiber.Ctx) error {
	shop := fortnite.NewRandomFortniteCatalog()
	return c.JSON(shop.GenerateFortniteCatalog())
}