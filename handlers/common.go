package handlers

import (
	"github.com/ectrc/snow/aid"
	"github.com/gofiber/fiber/v2"
)

func AnyNoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func PostTryPlayOnPlatform(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("true")
}

func GetEnabledFeatures(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON([]string{})
}

func PostGrantAccess(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("true")
}

func GetAccountReceipts(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON([]string{})
}

func GetSessionFindPlayer(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON([]string{})
}

func GetVersionCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(aid.JSON{
		"type": "NO_UPDATE",
	})
}

func GetContentPages(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(aid.JSON{
		"battlepassaboutmessages": aid.JSON{
			"news": aid.JSON{
				"messages": []aid.JSON{},
			},
			"lastModified": "0000-00-00T00:00:00.000Z",
		},
		"subgameselectdata": aid.JSON{
			"saveTheWorldUnowned": aid.JSON{
				"message": aid.JSON{
					"title": "Co-op PvE",
					"body": "Cooperative PvE storm-fighting adventure!",
					"spotlight": false,
					"hidden": true,
					"messagetype": "normal",
				},
			},
			"battleRoyale": aid.JSON{
				"message": aid.JSON{
					"title": "100 Player PvP",
					"body": "100 Player PvP Battle Royale.\n\nPvE progress does not affect Battle Royale.",
					"spotlight": false,
					"hidden": true,
					"messagetype": "normal",
				},
			},
			"creative": aid.JSON{
				"message": aid.JSON{
					"title": "New Featured Islands!",
					"body": "Your Island. Your Friends. Your Rules.\n\nDiscover new ways to play Fortnite, play community made games with friends and build your dream island.",
					"spotlight": false,
					"hidden": true,
					"messagetype": "normal",
				},
			},
			"lastModified": "0000-00-00T00:00:00.000Z",
		},
		"lastModified": "0000-00-00T00:00:00.000Z",
	})
}