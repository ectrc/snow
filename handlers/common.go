package handlers

import (
	"github.com/ectrc/snow/aid"
	p "github.com/ectrc/snow/person"
	"github.com/ectrc/snow/storage"
	"github.com/gofiber/fiber/v2"
)

func RedirectSocket(c *fiber.Ctx) error {
	return c.Redirect("/socket")
}

func AnyNoContent(c *fiber.Ctx) error {
	return c.SendStatus(204)
}

func PostGamePlatform(c *fiber.Ctx) error {
	return c.Status(200).SendString("true")
}

func GetGameEnabledFeatures(c *fiber.Ctx) error {
	return c.Status(200).JSON([]string{})
}

func PostGameAccess(c *fiber.Ctx) error {
	return c.Status(200).SendString("true")
}

func GetFortniteReceipts(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)
	receipts := []aid.JSON{}

	person.Receipts.RangeReceipts(func(key string, value *p.Receipt) bool {
		if value.State == "OK" {
			return true
		}
		
		receipts = append(receipts, value.GenerateFortniteReceiptEntry())
		return true
	})

	return c.Status(200).JSON(receipts)
}

func GetMatchmakingSession(c *fiber.Ctx) error {
	return c.Status(200).Send([]byte{})
}

func GetFortniteVersion(c *fiber.Ctx) error {
	return c.Status(200).JSON(aid.JSON{
		"type": "NO_UPDATE",
	})
}

func GetWaitingRoomStatus(c *fiber.Ctx) error {
	return c.SendStatus(204)
}

func GetAffiliate(c *fiber.Ctx) error {
	slugger := p.FindByDisplay(c.Params("slug"))
	if slugger == nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("Invalid affiliate slug"))
	}

	return c.Status(200).JSON(aid.JSON{
		"id": slugger.ID,
		"displayName": slugger.DisplayName,
		"slug": slugger.DisplayName,
		"status": "ACTIVE",
		"verified": false,
	})
}

func GetRegion(c *fiber.Ctx) error {
	return c.Status(200).JSON(aid.JSON{
		"continent": aid.JSON{
			"code": "EU",
		},
		"country": aid.JSON{
			"iso_code": "GB",
		},
		"subdivisions": []aid.JSON{},
	})
}

func SendJSONResponseFromAsset(c *fiber.Ctx, asset string) error {
	bytes := storage.Asset(asset)
	if bytes == nil {
		return c.Status(404).JSON(aid.JSON{})
	}

	stringBytes := string(*bytes)
	c.Set("Content-Type", "application/json")
	return c.Status(200).SendString(stringBytes)
}