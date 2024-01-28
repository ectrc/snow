package handlers

import (
	"github.com/ectrc/snow/aid"
	"github.com/gofiber/fiber/v2"
)

func RedirectSocket(c *fiber.Ctx) error {
	return c.Redirect("/socket")
}

func AnyNoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func PostGamePlatform(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("true")
}

func GetGameEnabledFeatures(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON([]string{})
}

func PostGameAccess(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("true")
}

func GetFortniteReceipts(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON([]string{})
}

func GetMatchmakingSession(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).Send([]byte{})
}

func GetFortniteVersion(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(aid.JSON{
		"type": "NO_UPDATE",
	})
}

func GetWaitingRoomStatus(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func GetRegion(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(aid.JSON{
		"continent": aid.JSON{
			"code": "EU",
		},
		"country": aid.JSON{
			"iso_code": "GB",
		},
		"subdivisions": []aid.JSON{},
	})
}