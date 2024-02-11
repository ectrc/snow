package handlers

import (
	"github.com/ectrc/snow/aid"
	p "github.com/ectrc/snow/person"
	"github.com/gofiber/fiber/v2"
)

func GetUserParties(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)

	response := aid.JSON{
		"current": []aid.JSON{},
		"invites": []aid.JSON{},
		"pending": []aid.JSON{},
		"pings": []aid.JSON{},
	}

	person.Parties.Range(func(key string, party *p.Party) bool {
		response["current"] = append(response["current"].([]aid.JSON), party.GenerateFortniteParty())
		return true
	})

	return c.Status(200).JSON(response)
}