package handlers

import (
	"github.com/ectrc/snow/aid"
	"github.com/gofiber/fiber/v2"
)

func GetFriendList(c *fiber.Ctx) error {
	return c.Status(200).JSON([]aid.JSON{})
}

func PostCreateFriend(c *fiber.Ctx) error {
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