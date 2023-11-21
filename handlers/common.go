package handlers

import (
	"strconv"

	"github.com/ectrc/snow/aid"
	"github.com/gofiber/fiber/v2"
)

func AnyNoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func PostGamePlatform(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("true")
}

func GetGameEnabledFeatures(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON([]string{})
}

func PostGameAccess(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("true")
}

func GetFortniteReceipts(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON([]string{})
}

func GetMatchmakingSession(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).Send([]byte{})
}

func GetFortniteVersion(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(aid.JSON{
		"type": "NO_UPDATE",
	})
}

func GetWaitingRoomStatus(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func GetRegion(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(aid.JSON{
		"continent": aid.JSON{
			"code": "EU",
		},
		"country": aid.JSON{
			"iso_code": "GB",
		},
		"subdivisions": []aid.JSON{},
	})
}

func GetContentPages(c *fiber.Ctx) error {
	seasonString := strconv.Itoa(aid.Config.Fortnite.Season)
	
	return c.Status(fiber.StatusOK).JSON(aid.JSON{
		"battlepassaboutmessages": aid.JSON{
			"news": aid.JSON{
				"messages": []aid.JSON{},
			},
			"lastModified": "0000-00-00T00:00:00.000Z",
		},
		"subgameselectdata": aid.JSON{
			"saveTheWorldUnowned": aid.JSON{
				"message": aid.JSON{
					"title": "Co-op PvE",
					"body": "Cooperative PvE storm-fighting adventure!",
					"spotlight": false,
					"hidden": true,
					"messagetype": "normal",
				},
			},
			"battleRoyale": aid.JSON{
				"message": aid.JSON{
					"title": "100 Player PvP",
					"body": "100 Player PvP Battle Royale.\n\nPvE progress does not affect Battle Royale.",
					"spotlight": false,
					"hidden": true,
					"messagetype": "normal",
				},
			},
			"creative": aid.JSON{
				"message": aid.JSON{
					"title": "New Featured Islands!",
					"body": "Your Island. Your Friends. Your Rules.\n\nDiscover new ways to play Fortnite, play community made games with friends and build your dream island.",
					"spotlight": false,
					"hidden": true,
					"messagetype": "normal",
				},
			},
			"lastModified": "0000-00-00T00:00:00.000Z",
		},
		"dynamicbackgrounds": aid.JSON{
			"backgrounds": aid.JSON{"backgrounds": []aid.JSON{
				{
					"key": "lobby",
					"stage": "season" + seasonString,
				},
				{
					"key": "vault",
					"stage": "season" + seasonString,
				},
			}},
			"lastModified": "0000-00-00T00:00:00.000Z",
		},
		"shopSections": aid.JSON{
			"sectionList": aid.JSON{
				"sections": []aid.JSON{
          {
            "bSortOffersByOwnership": false,
            "bShowIneligibleOffersIfGiftable": false,
            "bEnableToastNotification": true,
            "background":  aid.JSON{
              "stage": "default",
              "_type": "DynamicBackground",
              "key": "vault",
            },
            "_type": "ShopSection",
            "landingPriority": 0,
            "bHidden": false,
            "sectionId": "Featured",
            "bShowTimer": true,
            "sectionDisplayName": "Featured",
            "bShowIneligibleOffers": true,
          },
          {
            "bSortOffersByOwnership": false,
            "bShowIneligibleOffersIfGiftable": false,
            "bEnableToastNotification": true,
            "background":  aid.JSON{
              "stage": "default",
              "_type": "DynamicBackground",
              "key": "vault",
            },
            "_type": "ShopSection",
            "landingPriority": 1,
            "bHidden": false,
            "sectionId": "Daily",
            "bShowTimer": true,
            "sectionDisplayName": "Daily",
            "bShowIneligibleOffers": true,
          },
          {
            "bSortOffersByOwnership": false,
            "bShowIneligibleOffersIfGiftable": false,
            "bEnableToastNotification": false,
            "background":  aid.JSON{
              "stage": "default",
              "_type": "DynamicBackground",
              "key": "vault",
            },
            "_type": "ShopSection",
            "landingPriority": 2,
            "bHidden": false,
            "sectionId": "Battlepass",
            "bShowTimer": false,
            "sectionDisplayName": "Battle Pass",
            "bShowIneligibleOffers": false,
          },
          {
            "bSortOffersByOwnership": false,
            "bShowIneligibleOffersIfGiftable": false,
            "bEnableToastNotification": false,
            "background":  aid.JSON{
              "stage": "default",
              "_type": "DynamicBackground",
              "key": "vault",
            },
            "_type": "ShopSection",
            "landingPriority": 3,
            "bHidden": false,
            "sectionId": "SnowSection",
            "bShowTimer": false,
            "sectionDisplayName": "Snow Specials",
            "bShowIneligibleOffers": false,
          },
          {
            "bSortOffersByOwnership": false,
            "bShowIneligibleOffersIfGiftable": false,
            "bEnableToastNotification": false,
            "background":  aid.JSON{
              "stage": "default",
              "_type": "DynamicBackground",
              "key": "vault",
            },
            "_type": "ShopSection",
            "landingPriority": 3,
            "bHidden": false,
            "sectionId": "OGBundles",
            "bShowTimer": false,
            "sectionDisplayName": "OG Bundles",
            "bShowIneligibleOffers": false,
          },
        },
			},
			"lastModified": "0000-00-00T00:00:00.000Z",
		},
		"lastModified": "0000-00-00T00:00:00.000Z",
	})
}