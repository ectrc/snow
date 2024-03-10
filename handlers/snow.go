package handlers

import (
	"time"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/fortnite"
	p "github.com/ectrc/snow/person"
	"github.com/ectrc/snow/shop"
	"github.com/ectrc/snow/socket"
	"github.com/gofiber/fiber/v2"
)

func MiddlewareOnlyDebug(c *fiber.Ctx) error {
	if aid.Config.API.Debug {
		return c.Next()
	}

	return c.SendStatus(403)
}

func GetSnowPreloadedCosmetics(c *fiber.Ctx) error {
	return c.JSON(fortnite.DataClient)
}

func GetSnowCachedPlayers(c *fiber.Ctx) error {
	persons := p.AllFromCache()
	players := make([]p.PersonSnapshot, len(persons))

	for i, person := range persons {
		players[i] = *person.Snapshot()
	}

	return c.Status(200).JSON(players)
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
	shop := shop.GetShop()
	return c.JSON(shop.GenerateFortniteCatalogResponse())
}

func PostSnowLog(c *fiber.Ctx) error {
	var body struct {
		JSON aid.JSON `json:"json"`
		URL	string `json:"url"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(err.Error())	
	}

	aid.PrintJSON(body.JSON)
	return c.JSON(body)
}

func GetPlayer(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)
	return c.Status(200).JSON(aid.JSON{
		"snapshot": person.Snapshot(),
		"season": aid.JSON{
			"level": fortnite.DataClient.SnowSeason.GetSeasonLevel(person.CurrentSeasonStats),
			"xp": fortnite.DataClient.SnowSeason.GetRelativeSeasonXP(person.CurrentSeasonStats),
			"bookLevel": fortnite.DataClient.SnowSeason.GetBookLevel(person.CurrentSeasonStats),
			"bookXp": fortnite.DataClient.SnowSeason.GetRelativeBookXP(person.CurrentSeasonStats),
		},
	})
}

func GetPlayerOkay(c *fiber.Ctx) error {
	return c.Status(200).SendString("okay")
}

func PostPlayerCreateCode(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)
	code := person.ID + "=" + time.Now().Format("2006-01-02T15:04:05.999Z")
	encrypted, sig := aid.KeyPair.EncryptAndSignB64([]byte(code))	
	return c.Status(200).SendString(encrypted + "." + sig)
}

func GetLauncherStatus(c *fiber.Ctx) error {
	return c.Status(200).JSON(aid.JSON{
		"CurrentSeason": aid.Config.Fortnite.Season,
		"CurrentBuild": aid.Config.Fortnite.Build,
		"PlayersOnline": aid.FormatNumber(socket.JabberSockets.Len()),
	})
}