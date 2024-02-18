package handlers

import (
	"github.com/goccy/go-json"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/fortnite"
	"github.com/ectrc/snow/storage"
	"github.com/gofiber/fiber/v2"
)

func GetStorefrontCatalog(c *fiber.Ctx) error {
	shop := fortnite.NewRandomFortniteCatalog()
	return c.Status(200).JSON(shop.GenerateFortniteCatalog())
}

func GetStorefrontKeychain(c *fiber.Ctx) error {
	var keychain []string
	err := json.Unmarshal(*storage.Asset("keychain.json"), &keychain)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(aid.JSON{"error":err.Error()})
	}

	return c.Status(200).JSON(keychain)
}