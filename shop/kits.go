package shop

import (
	"fmt"
	"strings"

	"github.com/ectrc/snow/aid"
)

type StorefrontCatalogOfferMetaTypeStarterKit struct {
	TileSize string
	DisplayAssetPath string
	NewDisplayAssetPath string
	OriginalOffer int
	ExtraBonus int
	FeaturedImageURL string
	Priority int
	ReleaseSeason int
}

type StorefrontCatalogOfferTypeStarterKit struct {
	OfferType StorefrontCatalogOfferEnum
	Rewards []*StorefrontCatalogOfferGrant
	Price *StorefrontCatalogOfferPriceRealMoney
	Diplay *OfferDisplay
	Categories []string
	Meta *StorefrontCatalogOfferMetaTypeStarterKit
}

func NewStarterKitCatalogOffer() *StorefrontCatalogOfferTypeStarterKit {
	return &StorefrontCatalogOfferTypeStarterKit{
		OfferType: StorefrontCatalogOfferEnumStarterKit,
		Rewards: make([]*StorefrontCatalogOfferGrant, 0),
		Price: &StorefrontCatalogOfferPriceRealMoney{},
		Diplay: &OfferDisplay{},
		Categories: make([]string, 0),
		Meta: &StorefrontCatalogOfferMetaTypeStarterKit{},
	}
}

func (o *StorefrontCatalogOfferTypeStarterKit) GetOffer() *StorefrontCatalogOfferTypeStarterKit {
	return o
}

func (o *StorefrontCatalogOfferTypeStarterKit) GetOfferID() string {
	return fmt.Sprintf("kit://%s", strings.ReplaceAll(o.Diplay.Title, " ", ""))
}

func (o *StorefrontCatalogOfferTypeStarterKit) GetOfferType() StorefrontCatalogOfferEnum {
	return o.OfferType
}

func (o *StorefrontCatalogOfferTypeStarterKit) GetOfferPrice() *StorefrontCatalogOfferPriceRealMoney {
	return o.Price
}

func (o *StorefrontCatalogOfferTypeStarterKit) GetRewards() []*StorefrontCatalogOfferGrant {
	return o.Rewards
}

func (o *StorefrontCatalogOfferTypeStarterKit) GenerateFortniteCatalogOfferResponse() aid.JSON {
	return aid.JSON{
		"offerId": o.GetOfferID(),
		"offerType": "StaticPrice",
		"devName": fmt.Sprintf("[STARTER KIT] %s", o.Diplay.Title),
		"itemGrants": []aid.JSON{},
		"requirements": []aid.JSON{{
			"requirementType": "DenyOnFulfillment",
			"requiredId": o.GetOfferID(),
			"minQuantity": 1,
		}},
		"categories": o.Categories,
		"metaInfo": []aid.JSON{
			{
				"Key": "TileSize",
				"Value": o.Meta.TileSize,
			},
			{
				"Key": "NewDisplayAssetPath",
				"Value": o.Meta.NewDisplayAssetPath,
			},
			{
				"Key": "DisplayAssetPath",
				"Value": o.Meta.DisplayAssetPath,
			},
			{
				"key": "MtxQuantity",
				"value": o.Meta.OriginalOffer,
			},
			{
				"key": "MtxBonus",
				"value": o.Meta.ExtraBonus,
			},
		},
		"meta": aid.JSON{
			"TileSize": o.Meta.TileSize,
			"DisplayAssetPath": o.Meta.DisplayAssetPath,
			"NewDisplayAssetPath": o.Meta.NewDisplayAssetPath,
			"MtxQuantity": o.Meta.OriginalOffer,
			"MtxBonus": o.Meta.ExtraBonus,
		},
		"giftInfo": aid.JSON{
			"bIsEnabled": false,
			"forcedGiftBoxTemplateId": "",
			"purchaseRequirements": []aid.JSON{},
			"giftRecordIds": []string{},
		},
		"prices": []aid.JSON{{
			"currencyType": "RealMoney",
			"currencySubType": "",
			"regularPrice": -1,
			"dynamicRegularPrice": -1,
			"finalPrice": -1,
			"basePrice": -1,
			"saleExpiration": "9999-12-31T23:59:59.999Z",
		}},
		"displayAssetPath": o.Meta.DisplayAssetPath,
		"refundable": false,
		"title": o.Diplay.Title,
		"description": o.Diplay.Description,
		"shortDescription": o.Diplay.ShortDescription,
		"appStoreId": []string{
			"",
			o.GetOfferID(),
		},
		"dailyLimit": -1,
		"weeklyLimit": -1,
		"monthlyLimit": -1,
		"sortPriority": 0,
		"catalogGroupPriority": 0,
		"filterWeight": 0,
	}
}

func (o *StorefrontCatalogOfferTypeStarterKit) GenerateFortniteBulkOffersResponse() aid.JSON {
	return aid.JSON{
		"id": o.GetOfferID(),
		"title": o.Diplay.Title,
		"shortDescription": o.Diplay.ShortDescription,
		"longDescription": o.Diplay.LongDescription,
		"creationDate": "0000-00-00T00:00:00.000Z",
		"price": o.Price.LocalPrice,
		"currentPrice": o.Price.LocalPrice,
		"currencyCode": "USD",
		"basePrice": o.Price.BasePrice,
		"basePriceCurrencyCode": "GBP",
	}
}