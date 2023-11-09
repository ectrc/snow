package handlers

import (
	"github.com/goccy/go-json"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/storage"
	"github.com/gofiber/fiber/v2"
)

func GetStorefrontCatalog(c *fiber.Ctx) error {
	storefront := aid.JSON{
    "refreshIntervalHrs": 24,
    "dailyPurchaseHrs": 24,
    "expiration": aid.TimeEndOfDay(),
    "storefronts": []aid.JSON{},
	}

	return c.Status(fiber.StatusOK).JSON(storefront)
}

func GetStorefrontKeychain(c *fiber.Ctx) error {
	var keychain []string
	err := json.Unmarshal(*storage.Asset("keychain.json"), &keychain)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(aid.JSON{"error":err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(keychain)
}