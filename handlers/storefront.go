package handlers

import (
	"encoding/json"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/storage"
	"github.com/gofiber/fiber/v2"
)

func GetStorefrontCatalog(c *fiber.Ctx) error {
	t := aid.JSON{
		"refreshIntervalHrs": 24,
		"dailyPurchaseHrs": 24,
		"expiration": aid.TimeEndOfDay(),
		"storefronts": []aid.JSON{},
	}

	str := `{
          "name": "CurrencyStorefront",
          "catalogEntries": [
              {
                "devName": "[VIRTUAL]1 x Isabelle for 1200 MtxCurrency",
                "offerId": "v2:/f3d84c3ded015ae12a0c8ae3cc60d771a45df0d90f0af5e1cfbd454fa3083c94",
                "fulfillmentIds": [],
                "dailyLimit": -1,
                "weeklyLimit": -1,
                "monthlyLimit": -1,
                "categories": [
                    "Panel 03"
                ],
                "prices": [
                    {
                        "currencyType": "MtxCurrency",
                        "currencySubType": "",
                        "regularPrice": 1200,
                        "dynamicRegularPrice": 1200,
                        "finalPrice": 1200,
                        "saleExpiration": "9999-12-31T23:59:59.999Z",
                        "basePrice": 1200
                    }
                ],
                "meta": {
                    "NewDisplayAssetPath": "/Game/Catalog/DisplayAssets/DA_BR_Season8_BattlePass.DA_BR_Season8_BattlePass",
                    "offertag": "",
                    "SectionId": "Featured",
                    "TileSize": "Normal",
                    "AnalyticOfferGroupId": "3",
                    "ViolatorTag": "",
                    "ViolatorIntensity": "High",
                    "FirstSeen": ""
                },
                "matchFilter": "",
                "filterWeight": 0.0,
                "appStoreId": [],
                "requirements": [
                    {
                        "requirementType": "DenyOnItemOwnership",
                        "requiredId": "AthenaCharacter:CID_033_Athena_Commando_F_Medieval",
                        "minQuantity": 1
                    }
                ],
                "offerType": "StaticPrice",
                "giftInfo": {
                    "bIsEnabled": true,
                    "forcedGiftBoxTemplateId": "",
                    "purchaseRequirements": [],
                    "giftRecordIds": []
                },
                "refundable": true,
                "metaInfo": [
                    {
                        "key": "NewDisplayAssetPath",
                        "value": "/Game/Catalog/DisplayAssets/DA_BR_Season8_BattlePass.DA_BR_Season8_BattlePass"
                    },
                    {
                        "key": "offertag",
                        "value": ""
                    },
                    {
                        "key": "SectionId",
                        "value": "Featured"
                    },
                    {
                        "key": "TileSize",
                        "value": "Normal"
                    },
                    {
                        "key": "AnalyticOfferGroupId",
                        "value": "3"
                    },
                    {
                        "key": "ViolatorTag",
                        "value": ""
                    },
                    {
                        "key": "ViolatorIntensity",
                        "value": "High"
                    },
                    {
                        "key": "FirstSeen",
                        "value": ""
                    }
                ],
                "displayAssetPath": "/Game/Catalog/DisplayAssets/DA_BR_Season8_BattlePass.DA_BR_Season8_BattlePass",
                "itemGrants": [
                    {
                        "templateId": "AthenaCharacter:CID_033_Athena_Commando_F_Medieval",
                        "quantity": 1
                    }
                ],
                "additionalGrants": [],
                "sortPriority": -2,
                "catalogGroupPriority": 0
            }
          
          ]
      }`

	var catalog aid.JSON
	err := json.Unmarshal([]byte(str), &catalog)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(aid.JSON{"error":err.Error()})
	}

	t["storefronts"] = append(t["storefronts"].([]aid.JSON), catalog)

	return c.Status(fiber.StatusOK).JSON(t)
}

func GetStorefrontKeychain(c *fiber.Ctx) error {
	var keychain []string
	err := json.Unmarshal(*storage.Asset("keychain.json"), &keychain)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(aid.JSON{"error":err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(keychain)
}