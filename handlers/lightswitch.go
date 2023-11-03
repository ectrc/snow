package handlers

import (
	"github.com/ectrc/snow/aid"
	"github.com/gofiber/fiber/v2"
)

func GetLightswitchBulkStatus(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON([]aid.JSON{{
		"serviceInstanceId": "fortnite",
		"status": "UP",
		"banned": false,
	}})
}