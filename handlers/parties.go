package handlers

import (
	"time"

	"github.com/ectrc/snow/aid"
	p "github.com/ectrc/snow/person"
	"github.com/ectrc/snow/socket"
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

func GetUserPartyPrivacy(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)

	recieveIntents := person.CommonCoreProfile.Attributes.GetAttributeByKey("party.recieveIntents")
	if recieveIntents == nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("No Privacy Found"))
	}

	recieveInvites := person.CommonCoreProfile.Attributes.GetAttributeByKey("party.recieveInvites")
	if recieveIntents == nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("No Privacy Found"))
	}
	
	return c.Status(200).JSON(aid.JSON{
		"recieveIntents": aid.JSONParse(recieveIntents.ValueJSON),
		"recieveInvites": aid.JSONParse(recieveInvites.ValueJSON),
	})
}

func GetUserPartyNotifications(c *fiber.Ctx) error {
	return c.Status(200).JSON(aid.JSON{
		"pings": 0,
		"invites": 0,
	})
}

func GetPartyForMember(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)

	party, ok := person.Parties.Get(c.Params("partyId"))
	if !ok {
		return c.Status(400).JSON(aid.ErrorBadRequest("Party Not Found"))
	}

	return c.Status(200).JSON(party.GenerateFortniteParty())
}

func PostCreateParty(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)
	
	person.Parties.Range(func(key string, party *p.Party) bool {
		party.RemoveMember(person)
		return true
	})
	
	var body struct {
		Config map[string]interface{} `json:"config"`
		Meta map[string]interface{} `json:"meta"`
		JoinInformation struct {
			Meta map[string]interface{} `json:"meta"`
			Connection aid.JSON `json:"connection"`
		} `json:"join_info"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("Invalid Request"))
	}

	party := p.NewParty()
	party.UpdateMeta(body.Meta)
	party.UpdateConfig(body.Config)
	
	party.AddMember(person, "CAPTAIN")
	party.UpdateMemberMeta(person, body.JoinInformation.Meta)

	body.JoinInformation.Connection["connected_at"] = time.Now().Format(time.RFC3339)
	body.JoinInformation.Connection["updated_at"] = time.Now().Format(time.RFC3339)
	party.UpdateMemberConnections(person, body.JoinInformation.Connection)

	member := party.GetMember(person)
	if member == nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("Failed to create party"))
	}

	s, ok := socket.JabberSockets.Get(person.ID)
	if !ok {
		return c.Status(400).JSON(aid.ErrorBadRequest("No socket connection found"))
	}

	s.JabberSendMessageToPerson(aid.JSON{
		"ns": "Fortnite",
		"party_id": party.ID,
		"account_id": person.ID,
		"account_dn": person.DisplayName,
		"connection": body.JoinInformation.Connection,
		"member_state_updated": member.Meta,
		"updated_at": member.UpdatedAt.Format(time.RFC3339),
		"joined_at": member.JoinedAt.Format(time.RFC3339),
		"sent": time.Now().Format(time.RFC3339),
		"revision": 0,
		"type": "com.epicgames.social.party.notification.v0.MEMBER_JOINED",
	})

	return c.Status(200).JSON(party.GenerateFortniteParty())
}

func PatchUpdateParty(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)

	var body struct {
		Config map[string]interface{} `json:"config"`
		Meta struct {
			Update map[string]interface{} `json:"update"`
			Delete []string `json:"delete"`
		} `json:"meta"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("Invalid Request"))
	}

	party, ok := person.Parties.Get(c.Params("partyId"))
	if !ok {
		return c.Status(400).JSON(aid.ErrorBadRequest("Party Not Found"))
	}

	member := party.GetMember(person)
	if member == nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("Not in Party"))
	}

	if member.Role != "CAPTAIN" {
		return c.Status(400).JSON(aid.ErrorBadRequest("Not Captain"))
	}

	party.UpdateConfig(body.Config)
	party.UpdateMeta(body.Meta.Update)
	party.DeleteMeta(body.Meta.Delete)

	return c.Status(200).JSON(party.GenerateFortniteParty())
}

func PatchUpdatePartyMemberMeta(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)

	var body struct {
		Update map[string]interface{} `json:"update"`
		Delete []string `json:"delete"`
	}
	
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("Invalid Request"))
	}

	party, ok := person.Parties.Get(c.Params("partyId"))
	if !ok {
		return c.Status(400).JSON(aid.ErrorBadRequest("Party Not Found"))
	}

	member := party.GetMember(person)
	if member == nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("Not in Party"))
	}

	if c.Params("accountId") != person.ID {
		return c.Status(400).JSON(aid.ErrorBadRequest("Not owner of person"))
	}

	party.UpdateMemberMeta(person, body.Update)
	party.DeleteMemberMeta(person, body.Delete)

	return c.Status(200).JSON(party.GenerateFortniteParty())
}