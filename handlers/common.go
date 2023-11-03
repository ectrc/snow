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