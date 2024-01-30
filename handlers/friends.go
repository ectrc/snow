package handlers

import (
	"fmt"
	"time"

	"github.com/ectrc/snow/aid"
	p "github.com/ectrc/snow/person"
	"github.com/ectrc/snow/socket"
	"github.com/gofiber/fiber/v2"
)

func GetFriendList(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)
	result := []aid.JSON{}

	person.Relationships.Range(func(key string, value *p.Relationship) bool {
		switch value.Direction {
		case p.RelationshipInboundDirection:
			result = append(result, value.GenerateFortniteFriendEntry(p.GenerateTypeTowardsPerson))
		case p.RelationshipOutboundDirection:
			result = append(result, value.GenerateFortniteFriendEntry(p.GenerateTypeFromPerson))
		}
		return true
	})

	return c.Status(200).JSON(result)
}

func PostCreateFriend(c *fiber.Ctx) error {
	relationship, err := c.Locals("person").(*p.Person).CreateRelationship(c.Params("wanted"))
	if err != nil {
		aid.Print(fmt.Sprintf("Error creating relationship: %s", err.Error()))
		return c.Status(400).JSON(aid.ErrorBadRequest(err.Error()))
	}
	
	from, found := socket.GetJabberSocketByPersonID(relationship.From.ID)
	if found {
		from.JabberSendMessageToPerson(aid.JSON{
			"type": "com.epicgames.friends.core.apiobjects.Friend",
			"timestamp": time.Now().Format(time.RFC3339),
			"payload": relationship.GenerateFortniteFriendEntry(p.GenerateTypeFromPerson),
		})
	}

	towards, found := socket.GetJabberSocketByPersonID(relationship.Towards.ID)
	if found {
		towards.JabberSendMessageToPerson(aid.JSON{
			"type": "com.epicgames.friends.core.apiobjects.Friend",
			"timestamp": time.Now().Format(time.RFC3339),
			"payload": relationship.GenerateFortniteFriendEntry(p.GenerateTypeTowardsPerson),
		})
	}

	return c.SendStatus(204)
}

func DeleteFriend(c *fiber.Ctx) error {
	return c.SendStatus(204)
}

func GetFriendListSummary(c *fiber.Ctx) error {
	return c.Status(200).JSON([]aid.JSON{})
}

func GetPersonSearch(c *fiber.Ctx) error {
	return c.Status(200).JSON([]aid.JSON{})
}