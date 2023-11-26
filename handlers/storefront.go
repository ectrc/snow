package handlers

import (
	"github.com/goccy/go-json"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/storage"
	"github.com/gofiber/fiber/v2"
)

func GetStorefrontCatalog(c *fiber.Ctx) error {
// 	person := c.Locals("person").(*person.Person)
// 	storefront := fortnite.NewCatalog()

// 	daily := fortnite.NewStorefront("BRDailyStorefront")
// 	weekly := fortnite.NewStorefront("BRWeeklyStorefront")
	
// 	for len(weekly.CatalogEntries) < 8 {
// 		set := fortnite.Cosmetics.GetRandomSet()

// 		for _, cosmetic := range set.Items {
// 			if cosmetic.Type.BackendValue == "AthenaBackpack" {
// 				continue
// 			}

// 			entry := fortnite.NewCatalogEntry().Section("Featured").DisplayAsset(cosmetic.DisplayAssetPath).SetPrice(fortnite.GetPriceForRarity(cosmetic.Rarity.BackendValue))
// 			entry.AddGrant(cosmetic.Type.BackendValue + ":" + cosmetic.ID)

// 			if cosmetic.Backpack != "" {
// 				entry.AddGrant("AthenaBackpack:" + cosmetic.Backpack)
// 			}

// 			if cosmetic.Type.BackendValue != "AthenaCharacter" {
// 				entry.TileSize("Small")
// 			}

// 			if cosmetic.Type.BackendValue == "AthenaCharacter" {
// 				entry.TileSize("Normal")
// 				entry.Priority = -99999
// 			}

// 			entry.Panel = set.Name

// 			weekly.Add(*entry)
// 		}
// 	}

// 	storefront.Add(daily)
// 	storefront.Add(weekly)

	// return c.Status(fiber.StatusOK).JSON(storefront.GenerateFortniteCatalog(person))

	var x aid.JSON
	err := json.Unmarshal(*storage.Asset("hide_a.json"), &x)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(aid.JSON{"error":err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(x)
}

func GetStorefrontKeychain(c *fiber.Ctx) error {
	var keychain []string
	err := json.Unmarshal(*storage.Asset("keychain.json"), &keychain)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(aid.JSON{"error":err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(keychain)
}