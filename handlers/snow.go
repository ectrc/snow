package handlers

import (
	"strings"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/fortnite"
	p "github.com/ectrc/snow/person"
	"github.com/ectrc/snow/storage"
	"github.com/gofiber/fiber/v2"
)

func GetPreloadedCosmetics(c *fiber.Ctx) error {
	return c.JSON(fortnite.Cosmetics)
}

func GetPlaylistImage(c *fiber.Ctx) error {
	playlist := c.Params("playlist")
	if playlist == "" {
		return c.SendStatus(404)
	}
	playlist = strings.Split(playlist, ".")[0]

	image, ok := fortnite.PlaylistImages[playlist]
	if !ok {
		return c.SendStatus(404)
	}
	
	c.Set("Content-Type", "image/png")
	return c.Send(image)
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

func GetCachedPlayers(c *fiber.Ctx) error {
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