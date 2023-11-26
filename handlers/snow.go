package handlers

import (
	"github.com/ectrc/snow/fortnite"
	"github.com/gofiber/fiber/v2"
)

func GetPrelaodedCosmetics(c *fiber.Ctx) error {
	return c.JSON(fortnite.Cosmetics)
}