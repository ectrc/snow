package handlers

import (
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
		return c.Status(400).JSON(aid.ErrorBadRequest(err.Error()))
	}
	
	from, found := socket.JabberSockets.Get(relationship.From.ID)
	if found {
		from.JabberSendMessageToPerson(aid.JSON{
			"type": "com.epicgames.friends.core.apiobjects.Friend",
			"timestamp": time.Now().Format(time.RFC3339),
			"payload": relationship.GenerateFortniteFriendEntry(p.GenerateTypeFromPerson),
		})
		from.JabberNotifyFriends()
	}

	towards, found := socket.JabberSockets.Get(relationship.Towards.ID)
	if found {
		towards.JabberSendMessageToPerson(aid.JSON{
			"type": "com.epicgames.friends.core.apiobjects.Friend",
			"timestamp": time.Now().Format(time.RFC3339),
			"payload": relationship.GenerateFortniteFriendEntry(p.GenerateTypeTowardsPerson),
		})
		towards.JabberNotifyFriends()
	}

	return c.SendStatus(204)
}

func DeleteFriend(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)
	
	relationship, found := person.Relationships.Get(c.Params("wanted"))
	if !found {
		return c.Status(404).JSON(aid.ErrorNotFound)
	}

	from, found := socket.JabberSockets.Get(relationship.From.ID)
	if found {
		from.JabberSendMessageToPerson(aid.JSON{
			"type": "com.epicgames.friends.core.apiobjects.FriendRemoval",
			"timestamp": time.Now().Format(time.RFC3339),
			"payload": relationship.GenerateFortniteFriendRemovalEntry(p.GenerateTypeFromPerson),
		})
		from.JabberNotifyFriends()
	}

	towards, found := socket.JabberSockets.Get(relationship.Towards.ID)
	if found {
		towards.JabberSendMessageToPerson(aid.JSON{
			"type": "com.epicgames.friends.core.apiobjects.FriendRemoval",
			"timestamp": time.Now().Format(time.RFC3339),
			"payload": relationship.GenerateFortniteFriendRemovalEntry(p.GenerateTypeTowardsPerson),
		})
		towards.JabberNotifyFriends()
	}

	relationship.Delete()
	person.Relationships.Delete(c.Params("friend"))

	return c.SendStatus(204)
}

func GetFriendListSummary(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)

	summary := aid.JSON{
		"friends": []aid.JSON{},
		"blocklist": []aid.JSON{},
		"incoming": []aid.JSON{},
		"outgoing": []aid.JSON{},
		"suggested": []aid.JSON{},
	}

	person.Relationships.Range(func(key string, value *p.Relationship) bool {
		switch value.Direction {
		case p.RelationshipInboundDirection:
			res := value.GenerateFortniteSummaryEntry(p.GenerateTypeTowardsPerson)
			if value.Status == "ACCEPTED" {
				summary["friends"] = append(summary["friends"].([]aid.JSON), res)
				break
			}

			summary["incoming"] = append(summary["incoming"].([]aid.JSON), res)
		case p.RelationshipOutboundDirection:
			res := value.GenerateFortniteSummaryEntry(p.GenerateTypeFromPerson)
			if value.Status == "ACCEPTED" {
				summary["friends"] = append(summary["friends"].([]aid.JSON), res)
				break
			}
			
			summary["outgoing"] = append(summary["outgoing"].([]aid.JSON), res)
		}
			
		return true
	})

	return c.Status(200).JSON(summary)
}

func GetPersonSearch(c *fiber.Ctx) error {
	return c.Status(200).JSON([]aid.JSON{})
}