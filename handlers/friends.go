package handlers

import (
	"time"

	"github.com/ectrc/snow/aid"
	p "github.com/ectrc/snow/person"
	"github.com/ectrc/snow/storage"
	"github.com/gofiber/fiber/v2"
)

func GetFriendList(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)
	result := map[string]aid.JSON{}

	for _, partial := range storage.Repo.GetFriendsForPerson(person.ID) {
		friend := person.GetFriend(partial.ID)
		if friend == nil {
			continue
		}

		result[partial.ID] = friend.GenerateFriendResponse()
	}

	response := []aid.JSON{}
	for _, friend := range result {
		response = append(response, friend)
	}

	return c.Status(200).JSON(response)
}

func PostCreateFriend(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)
	friendId := c.Params("wanted")

	existing := person.GetFriend(friendId)
	if existing != nil && (existing.Direction == "BOTH" || existing.Direction == "OUTGOING") {
		return c.Status(400).JSON(aid.ErrorBadRequest("already active friend request"))
	}

	person.AddFriend(friendId)
	
	socket := FindSocketForPerson(person)
	socket.PresenceWriteJSON(aid.JSON{
		"payload": person.GetFriend(friendId).GenerateFriendResponse(),
		"type": "com.epicgames.friends.core.apiobjects.Friend",
		"timestamp": time.Now().Format(time.RFC3339),
	})

	friendSocket := FindSocketForPerson(p.Find(friendId))
	friendSocket.PresenceWriteJSON(aid.JSON{
		"payload": friendSocket.Person.GetFriend(person.ID).GenerateFriendResponse(),
		"type": "com.epicgames.friends.core.apiobjects.Friend",
		"timestamp": time.Now().Format(time.RFC3339),
	})

	return c.SendStatus(204)
}

func DeleteFriend(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)
	wanted := c.Params("wanted")

	existing := person.GetFriend(wanted)
	if existing == nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("not friends"))
	}

	person.RemoveFriend(wanted)
	existing.Person.RemoveFriend(person.ID)

	return c.SendStatus(204)
}

func GetFriendListSummary(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)

	all := map[string]*p.Friend{}
	for _, partial := range storage.Repo.GetFriendsForPerson(person.ID) {
		friend := person.GetFriend(partial.ID)
		if friend == nil {
			continue
		}

		all[partial.ID] = friend
	}

	result := aid.JSON{
		"friends": []aid.JSON{},
		"incoming": []aid.JSON{},
		"outgoing": []aid.JSON{},
		"settings": aid.JSON{
			"acceptInvites": "public",
		},
	}

	for _, friend := range all {
		switch friend.Status {
		case p.FriendStatusAccepted:
			result["friends"] = append(result["friends"].([]aid.JSON), friend.GenerateSummaryResponse())
		case p.FriendStatusPending:
			switch friend.Direction {
			case p.FriendDirectionIncoming:
				result["incoming"] = append(result["incoming"].([]aid.JSON), friend.GenerateSummaryResponse())
			case p.FriendDirectionOutgoing:
				result["outgoing"] = append(result["outgoing"].([]aid.JSON), friend.GenerateSummaryResponse())
			}
		}
	}

	return c.Status(200).JSON(result)
}

func GetPersonSearch(c *fiber.Ctx) error {
	query := c.Query("prefix")

	matches := storage.Repo.GetPersonsByPartialDisplayFromDB(query)
	if matches == nil {
		return c.Status(200).JSON([]aid.JSON{})
	}

	result := []aid.JSON{}
	for i, match := range matches {
		result = append(result, aid.JSON{
			"accountId": match.ID,
			"epicMutuals": 0,
			"sortPosition": i,
			"matchType": "prefix",
			"matches": []aid.JSON{{
				"value": match.DisplayName,
				"matchType": "prefix",
			}},
		})
	}

	return c.Status(200).JSON(result)
}