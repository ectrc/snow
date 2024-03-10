package handlers

import (
	"strings"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/person"
	p "github.com/ectrc/snow/person"
	"github.com/ectrc/snow/shop"
	"github.com/ectrc/snow/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetHtmlPurchasePage(c *fiber.Ctx) error {
	c.Set("X-UEL", "DEFAULT")
	c.Set("X-Download-Options", "noopen")
	c.Set("X-DNS-Prefetch-Control", "off")
	c.Set("x-epic-correlation-id", uuid.New().String())
	c.Set("X-Frame-Options", "SAMEORIGIN")

	var cookies struct {
		Token string `cookie:"EPIC_BEARER_TOKEN"`
	}

	if err := c.CookieParser(&cookies); err != nil {
		return c.SendStatus(401)
	}

	if cookies.Token == "" {
		return c.SendStatus(401)
	}

	person, err := aid.GetSnowFromToken(cookies.Token)
	if err != nil {
		return c.SendStatus(401)
	}
	c.Locals("person", person)

	fileBytes := storage.Asset("purchase.html")
	if fileBytes == nil {
		return c.SendStatus(404)
	}

	c.Set("content-type", "text/html")
	return c.SendString(string(*fileBytes))
}

func GetPurchaseAsset(c *fiber.Ctx) error {
	asset := c.Query("asset")

	type_ := strings.Split(asset, ".")
	fileBytes := storage.Asset(asset)
	if fileBytes == nil {
		return c.SendStatus(404)
	}
	
	c.Set("content-type", "text/" + type_[1])
	return c.SendString(string(*fileBytes))
}

func GetPurchaseOffer(c *fiber.Ctx) error {
	player := c.Locals("person").(*person.Person)
	offerId := c.Query("offerId")
	if offerId == "" {
		return c.SendStatus(400)
	}

	store := shop.GetShop()
	offerRaw, type_ := store.GetOfferByID(offerId)
	if offerRaw == nil {
		return c.SendStatus(404)
	}

	response := aid.JSON{
		"user": aid.JSON{
			"displayName": player.DisplayName,
		},	
	}

	switch type_ {
	case shop.StorefrontCatalogOfferEnumCurrency:
		offer := offerRaw.(*shop.StorefrontCatalogOfferTypeCurrency)
		response["offer"] = aid.JSON{
			"id": offer.GetOfferID(),
			"price": aid.FormatPrice(int(offer.Price.LocalPrice)),
			"name": offer.Diplay.Title,
			"imageUrl": offer.Meta.FeaturedImageURL,
			"type": "currency",
		}
	case shop.StorefrontCatalogOfferEnumStarterKit:
		offer := offerRaw.(*shop.StorefrontCatalogOfferTypeStarterKit)
		response["offer"] = aid.JSON{
			"id": offer.GetOfferID(),
			"price": aid.FormatPrice(int(offer.Price.LocalPrice)),
			"name": offer.Diplay.Title,
			"imageUrl": offer.Meta.FeaturedImageURL,
			"type": "starterpack",
		}
	default:
		break
	}

	return c.Status(200).JSON(response)
}

func PostPurchaseOffer(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)
	
	var body struct {
		OfferId string `json:"offerId" binding:"required"`
		Type string `json:"type" binding:"required"` // "currency" or "starterpack"
	}

	aid.PrintJSON(body)

	if err := c.BodyParser(&body); err != nil {
		return c.SendStatus(400)
	}

	lookup := map[string]func(*fiber.Ctx, *p.Person, string) error{
		"currency": purchaseCurrency,
		"starterpack": purchaseStarterPack,
	}

	if handler, ok := lookup[body.Type]; ok {
		return handler(c, person, body.OfferId)
	}

	return c.SendStatus(400)
}

func purchaseCurrency(c *fiber.Ctx, person *p.Person, offerId string) error {
	offerRaw, type_ := shop.GetShop().GetOfferByID(offerId)
	if offerRaw == nil {
		return c.Status(404).JSON(aid.ErrorNotFound)
	}
	if type_ != shop.StorefrontCatalogOfferEnumCurrency {
		return c.Status(400).JSON(aid.ErrorBadRequest)
	}
	offer := offerRaw.(*shop.StorefrontCatalogOfferTypeCurrency)

	receipt := p.NewReceipt(offerId, int(offer.Price.BasePrice))
	for _, grant := range offer.Rewards {
		item := p.NewItem(grant.TemplateID, grant.Quantity)
		item.ProfileType = string(grant.ProfileType)
		receipt.AddLoot(item)
	}
	person.Receipts.AddReceipt(receipt).Save()
 
	return c.Status(200).JSON(aid.JSON{
		"receipt": receipt.GenerateUnrealReceiptEntry(),
	})
}

func purchaseStarterPack(c *fiber.Ctx, person *p.Person, offerId string) error {
	offerRaw, type_ := shop.GetShop().GetOfferByID(offerId)
	if offerRaw == nil {
		return c.Status(404).JSON(aid.ErrorNotFound)
	}
	if type_ != shop.StorefrontCatalogOfferEnumStarterKit {
		return c.Status(400).JSON(aid.ErrorBadRequest)
	}
	offer := offerRaw.(*shop.StorefrontCatalogOfferTypeStarterKit)

	receipt := p.NewReceipt(offerId, int(offer.Price.BasePrice))
	for _, grant := range offer.Rewards {
		item := p.NewItem(grant.TemplateID, grant.Quantity)
		item.ProfileType = string(grant.ProfileType)
		receipt.AddLoot(item)
	}
	person.Receipts.AddReceipt(receipt).Save()
 
	return c.Status(200).JSON(aid.JSON{
		"receipt": receipt.GenerateUnrealReceiptEntry(),
	})
}