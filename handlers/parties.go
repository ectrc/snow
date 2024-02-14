package handlers

import (
	"github.com/ectrc/snow/aid"
	p "github.com/ectrc/snow/person"
	"github.com/ectrc/snow/socket"
	"github.com/gofiber/fiber/v2"
)

func GetPartiesForUser(c *fiber.Ctx) error {
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

	person.Invites.Range(func(key string, invite *p.PartyInvite) bool {
		response["invites"] = append(response["invites"].([]aid.JSON), invite.GenerateFortnitePartyInvite())
		return true
	})

	return c.Status(200).JSON(response)
}

func GetPartyUserPrivacy(c *fiber.Ctx) error {
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

func GetPartyNotifications(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)
	return c.Status(200).JSON(aid.JSON{
		"pings": 0,
		"invites": person.Invites.Len(),
	})
}

func GetPartyForMember(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)

	party, ok := p.Parties.Get(c.Params("partyId"))
	if !ok {
		return c.Status(400).JSON(aid.ErrorBadRequest("Party Not Found"))
	}

	aid.Print(person.DisplayName, " is getting party ", party.ID)

	return c.Status(200).JSON(party.GenerateFortniteParty())
}

func GetPartyPingsFromFriend(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)

	friend := p.Find(c.Params("friendId"))
	if friend == nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("Friend Not Found"))
	}

	pings := []aid.JSON{}
	person.Invites.Range(func(key string, ping *p.PartyInvite) bool {
		if ping.Inviter.ID == friend.ID {
			pings = append(pings, ping.Party.GenerateFortniteParty())
		}
		return true
	})

	return c.Status(200).JSON(pings)
}

func PostPartyCreate(c *fiber.Ctx) error {
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
	
	party.AddMember(person)
	party.UpdateMemberMeta(person, body.JoinInformation.Meta)
	party.UpdateMemberConnection(person, body.JoinInformation.Connection)

	member := party.GetMember(person)
	party.PromoteMember(member)
	socket.EmitPartyMemberJoined(party, member)

	return c.Status(200).JSON(party.GenerateFortniteParty())
}

func PatchPartyUpdateState(c *fiber.Ctx) error {
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
	socket.EmitPartyMetaUpdated(party, body.Meta.Update, body.Meta.Delete, body.Meta.Update)

	return c.SendStatus(204)
}

func PatchPartyUpdateMemberState(c *fiber.Ctx) error {
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
	socket.EmitPartyMemberMetaUpdated(party, member, body.Update, body.Delete)

	return c.SendStatus(204)
}

func DeletePartyMember(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)

	party, ok := person.Parties.Get(c.Params("partyId"))
	if !ok {
		return c.Status(400).JSON(aid.ErrorBadRequest("Party Not Found"))
	}

	member := party.GetMember(person)
	if member == nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("Not in Party"))
	}

	socket.EmitPartyMemberLeft(party, member)
	party.RemoveMember(person)
	
	if party.Captain.Person.ID == person.ID && len(party.Members) > 0 {
		party.PromoteMember(party.GetFirstMember())
		socket.EmitPartyNewCaptain(party)
	}

	return c.SendStatus(204)
}

func PostPartyInvite(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)

	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("Invalid Request"))
	}

	towards := p.Find(c.Params("accountId"))
	if towards == nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("Person Not Found"))
	}

	party, ok := person.Parties.Get(c.Params("partyId"))
	if !ok {
		return c.Status(400).JSON(aid.ErrorBadRequest("Party Not Found"))
	}
	
	invite := p.NewPartyInvite(party, person, towards, body)
	party.AddInvite(invite)
	towards.Invites.Set(party.ID, invite)

	socket.EmitPartyInvite(invite)

	return c.SendStatus(204)
}

func PostPartyJoin(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)
	if person.Parties.Len() != 0 {
		return c.Status(400).JSON(aid.ErrorBadRequest("Already in a party"))
	}

	party, ok := p.Parties.Get(c.Params("partyId"))
	if !ok {
		return c.Status(400).JSON(aid.ErrorBadRequest("Party Not Found"))
	}

	var body struct {
		Meta map[string]interface{} `json:"meta"`
		Connection aid.JSON `json:"connection"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("Invalid Request"))
	}

	party.AddMember(person)
	party.UpdateMemberMeta(person, body.Meta)
	party.UpdateMemberConnection(person, body.Connection)
	
	member := party.GetMember(person)
	socket.EmitPartyMemberJoined(party, member)
	socket.EmitPartyMemberMetaUpdated(party, party.GetMember(person), body.Meta, []string{})
	socket.EmitPartyMetaUpdated(party, party.Meta, []string{}, map[string]interface{}{})

	return c.Status(200).JSON(aid.JSON{
		"party_id": party.ID,
		"status": "JOINED",
	})
}

func PostPartyPromoteMember(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)
	party, ok := person.Parties.Get(c.Params("partyId"))
	if !ok {
		return c.Status(400).JSON(aid.ErrorBadRequest("Party Not Found"))
	}

	member := party.GetMember(p.Find(c.Params("accountId")))
	if member == nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("Member Not Found"))
	}

	if party.Captain.Person.ID != person.ID {
		return c.Status(400).JSON(aid.ErrorBadRequest("Not Captain"))
	}

	party.PromoteMember(member)
	socket.EmitPartyNewCaptain(party)

	return c.SendStatus(204)
}