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

	person.IncomingRelationships.Range(func(key string, value *p.Relationship[p.RelationshipInboundDirection]) bool {
		result = append(result, value.GenerateFortniteFriendEntry())
		return true
	})

	person.OutgoingRelationships.Range(func(key string, value *p.Relationship[p.RelationshipOutboundDirection]) bool {
		result = append(result, value.GenerateFortniteFriendEntry())
		return true
	})

	return c.Status(200).JSON(result)
}

func PostCreateFriend(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)
	wanted := c.Params("wanted")

	direction, err := person.CreateRelationship(wanted)
	if err != nil {
		aid.Print(fmt.Sprintf("Error creating relationship: %s", err.Error()))
		return c.Status(400).JSON(aid.ErrorBadRequest(err.Error()))
	}

	socket.JabberSockets.Range(func(key string, value *socket.Socket[socket.JabberData]) bool {
		aid.Print(fmt.Sprintf("Checking socket: %s", key, value))
		return true
	})

	personSocket, found := socket.GetJabberSocketByPersonID(wanted)
	aid.Print(fmt.Sprintf("Found socket: %t", found), fmt.Sprintf("Direction: %s", direction))
	if found {
		payload := aid.JSON{}
		switch direction {
		case "INBOUND":
			payload = personSocket.Person.IncomingRelationships.MustGet(wanted).GenerateFortniteFriendEntry()
		case "OUTBOUND":
			payload = personSocket.Person.OutgoingRelationships.MustGet(wanted).GenerateFortniteFriendEntry()
		}

		personSocket.JabberSendMessageToPerson(aid.JSON{
			"type": "com.epicgames.friends.core.apiobjects.Friend",
			"timestamp": time.Now().Format(time.RFC3339),
			"payload": payload,
		})
	}

	friendSocket, found := socket.GetJabberSocketByPersonID(wanted)
	aid.Print(fmt.Sprintf("Found friend socket: %t", found), fmt.Sprintf("Direction: %s", direction))
	if found {
		payload := aid.JSON{}
		switch direction {
		case "INBOUND":
			payload = friendSocket.Person.OutgoingRelationships.MustGet(person.ID).GenerateFortniteFriendEntry()
		case "OUTBOUND":
			payload = friendSocket.Person.IncomingRelationships.MustGet(person.ID).GenerateFortniteFriendEntry()
		}

		friendSocket.JabberSendMessageToPerson(aid.JSON{
			"type": "com.epicgames.friends.core.apiobjects.Friend",
			"timestamp": time.Now().Format(time.RFC3339),
			"payload": payload,
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