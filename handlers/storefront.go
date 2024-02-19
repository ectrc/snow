package handlers

import (
	"strings"

	"github.com/goccy/go-json"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/fortnite"
	"github.com/ectrc/snow/storage"
	"github.com/gofiber/fiber/v2"
)

func GetStorefrontCatalog(c *fiber.Ctx) error {
	shop := fortnite.NewRandomFortniteCatalog()
	return c.Status(200).JSON(shop.GenerateFortniteCatalogResponse())
}

func GetStorefrontKeychain(c *fiber.Ctx) error {
	var keychain []string
	err := json.Unmarshal(*storage.Asset("keychain.json"), &keychain)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(aid.JSON{"error":err.Error()})
	}

	return c.Status(200).JSON(keychain)
}

func GetStorefrontCatalogBulkOffers(c *fiber.Ctx) error {
	shop := fortnite.NewRandomFortniteCatalog()

	appStoreIdBytes := c.Request().URI().QueryArgs().PeekMulti("id")
	appStoreIds := make([]string, len(appStoreIdBytes))
	for i, id := range appStoreIdBytes {
		appStoreIds[i] = string(id)
	}

	response := aid.JSON{}
	for _, id := range appStoreIds {
		offer := shop.FindCurrencyOfferById(strings.ReplaceAll(id, "app-", ""))
		if offer == nil {
			continue
		}

		response[id] = offer.GenerateFortniteCatalogBulkOfferResponse()
	}

	for _, id := range appStoreIds {
		offer := shop.FindStarterPackById(strings.ReplaceAll(id, "app-", ""))
		if offer == nil {
			continue
		}

		response[id] = offer.GenerateFortniteCatalogBulkOfferResponse()
	}

	return c.Status(200).JSON(response)
}