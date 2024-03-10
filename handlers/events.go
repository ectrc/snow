// structs from https://github.com/FabianFG/Fortnite-Api/

package handlers

import (
	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/fortnite"
	"github.com/ectrc/snow/person"
	"github.com/gofiber/fiber/v2"
)

var (
	hypeTokens = map[int]string{
		0: "ARENA_S8_Division1",
		25: "ARENA_S8_Division2",
		75: "ARENA_S8_Division3",
		125: "ARENA_S8_Division4",
		175: "ARENA_S8_Division5",
		225: "ARENA_S8_Division6",
		300: "ARENA_S8_Division7",
	}
)

func GetEvents(c *fiber.Ctx) error {
	person := c.Locals("person").(*person.Person)

	events := []aid.JSON{}
	templates := []aid.JSON{}
	tokens := []string{}

	for _, event := range fortnite.ArenaEvents {
		events = append(events, event.GenerateFortniteEvent())

		for _, window := range event.Windows {
			templates = append(templates, window.Template.GenerateFortniteEventTemplate())
		}
	}

	for limit, token := range hypeTokens {
		if person.CurrentSeasonStats.Hype >= limit {
			tokens = []string{token}
		}
	}

	return c.Status(200).JSON(aid.JSON{
		"player": aid.JSON{
			"gameId": "Fortnite",
			"accountId": person.ID,
			"tokens": tokens,
			"teams": aid.JSON{},
			"pendingPayouts": []string{},
			"pendingPenalties": aid.JSON{},
			"persistentScores": aid.JSON{
				"Hype": person.CurrentSeasonStats.Hype,
			},
			"groupIdentity": aid.JSON{},
		},
		"events": events,
		"templates": templates,
	})
}

func GetEventsBulkHistory(c *fiber.Ctx) error {
	return c.Status(200).JSON([]aid.JSON{})
}