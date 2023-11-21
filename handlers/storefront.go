package handlers

import (
	"github.com/goccy/go-json"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/fortnite"
	"github.com/ectrc/snow/person"
	"github.com/ectrc/snow/storage"
	"github.com/gofiber/fiber/v2"
)

func GetStorefrontCatalog(c *fiber.Ctx) error {
	person := c.Locals("person").(*person.Person)
	storefront := fortnite.NewCatalog()

	bundleStorefront := fortnite.NewStorefront("bundles")
	{
		bundle := fortnite.NewBundleEntry("v2:/hello_og", "OG Bundle", 300)
		bundle.Asset = "/Game/Catalog/NewDisplayAssets/DAv2_CID_A_183_M_AntiquePal_S7A9W.DAv2_CID_A_183_M_AntiquePal_S7A9W"
		bundle.AddBundleGrant(*fortnite.NewBundleItem("AthenaCharacter:CID_028_Athena_Commando_F", 1000, 500, 800))
		bundle.AddBundleGrant(*fortnite.NewBundleItem("AthenaCharacter:CID_001_Athena_Commando_F", 1000, 500, 800))
		bundle.AddMeta("AnalyticOfferGroupId", "3")
		bundle.AddMeta("SectionId", "OGBundles")
		bundle.AddMeta("TileSize", "DoubleWide")
		bundle.AddMeta("NewDisplayAssetPath", bundle.Asset)
		bundleStorefront.Add(*bundle)

		random := fortnite.NewItemEntry("v2:/random", "Random Bundle", 300)
		random.AddGrant("AthenaCharacter:CID_Random")
		random.AddMeta("AnalyticOfferGroupId", "3")
		random.AddMeta("SectionId", "OGBundles")
		random.AddMeta("TileSize", "DoubleWide")

		bundleStorefront.Add(*random)
	}
	storefront.Add(bundleStorefront)

	return c.Status(fiber.StatusOK).JSON(storefront.GenerateFortniteCatalog(person))
}

func GetStorefrontKeychain(c *fiber.Ctx) error {
	var keychain []string
	err := json.Unmarshal(*storage.Asset("keychain.json"), &keychain)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(aid.JSON{"error":err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(keychain)
}