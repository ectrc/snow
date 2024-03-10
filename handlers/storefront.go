package handlers

import (
	"github.com/goccy/go-json"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/shop"
	"github.com/ectrc/snow/storage"
	"github.com/gofiber/fiber/v2"
)

func GetStorefrontCatalog(c *fiber.Ctx) error {
	shop := shop.GetShop()
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
	store := shop.GetShop()

	appStoreIdBytes := c.Request().URI().QueryArgs().PeekMulti("id")
	appStoreIds := make([]string, len(appStoreIdBytes))
	for i, id := range appStoreIdBytes {
		appStoreIds[i] = string(id)
	}

	response := aid.JSON{}
	for _, id := range appStoreIds {
		offerRaw, type_ := store.GetOfferByID(id)
		if offerRaw == nil {
			continue
		}

		switch type_ {
		case shop.StorefrontCatalogOfferEnumCurrency:
			offer := offerRaw.(*shop.StorefrontCatalogOfferTypeCurrency)
			response[id] = offer.GenerateFortniteBulkOffersResponse()
		case shop.StorefrontCatalogOfferEnumStarterKit:
			offer := offerRaw.(*shop.StorefrontCatalogOfferTypeStarterKit)
			response[id] = offer.GenerateFortniteBulkOffersResponse()
		default:
			break
		}
	}

	return c.Status(200).JSON(response)
}