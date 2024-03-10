package shop

import (
	"fmt"

	"github.com/ectrc/snow/aid"
)

type StorefrontCatalogOfferMetaTypeBattlePass struct {
	TileSize string
	SectionID string
	DisplayAssetPath string
	NewDisplayAssetPath string
	Priority int
	OnlyOnce bool
}

type StorefrontCatalogOfferTypeBattlePass struct {
	OfferID string
	OfferType StorefrontCatalogOfferEnum
	Rewards []*StorefrontCatalogOfferGrant
	Price *StorefrontCatalogOfferPriceMtxCurrency
	Diplay *OfferDisplay
	Categories []string
	Meta *StorefrontCatalogOfferMetaTypeBattlePass
}

func NewBattlePassCatalogOffer() *StorefrontCatalogOfferTypeBattlePass {
	return &StorefrontCatalogOfferTypeBattlePass{
		OfferID: aid.RandomString(32),
		OfferType: StorefrontCatalogOfferEnumBattlePass,
		Rewards: make([]*StorefrontCatalogOfferGrant, 0),
		Price: &StorefrontCatalogOfferPriceMtxCurrency{},
		Diplay: &OfferDisplay{},
		Categories: make([]string, 0),
		Meta: &StorefrontCatalogOfferMetaTypeBattlePass{},
	}
}

func (o *StorefrontCatalogOfferTypeBattlePass) GetOffer() *StorefrontCatalogOfferTypeBattlePass {
	return o
}

func (o *StorefrontCatalogOfferTypeBattlePass) GetOfferID() string {
	return o.OfferID
}

func (o *StorefrontCatalogOfferTypeBattlePass) GetOfferType() StorefrontCatalogOfferEnum {
	return o.OfferType
}

func (o *StorefrontCatalogOfferTypeBattlePass) GetOfferPrice() *StorefrontCatalogOfferPriceMtxCurrency {
	return o.Price
}

func (o *StorefrontCatalogOfferTypeBattlePass) GetRewards() []*StorefrontCatalogOfferGrant {
	return o.Rewards
}

func (o *StorefrontCatalogOfferTypeBattlePass) GenerateFortniteCatalogOfferResponse() aid.JSON {
	return aid.JSON{
		"offerId": o.OfferID,
		"offerType": "StaticPrice",
		"devName": fmt.Sprintf("[BOOK] %s", o.Diplay.ShortDescription),
		"itemGrants": []string{},
		"requirements": aid.Ternary[[]aid.JSON](o.Meta.OnlyOnce, []aid.JSON{{
			"requirementType": "DenyOnFulfillment",
			"requiredId": o.GetOfferID(),
			"minQuantity": 1,
		}}, []aid.JSON{}),
		"fulfillmentIds": []string{o.OfferID},
		"categories": o.Categories,
		"metaInfo": []aid.JSON{
			{
				"Key": "TileSize",
				"Value": o.Meta.TileSize,
			},
			{
				"Key": "SectionId",
				"Value": o.Meta.SectionID,
			},
			{
				"Key": "NewDisplayAssetPath",
				"Value": o.Meta.NewDisplayAssetPath,
			},
			{
				"Key": "DisplayAssetPath",
				"Value": o.Meta.DisplayAssetPath,
			},
		},
		"meta": aid.JSON{
			"TileSize": o.Meta.TileSize,
			"SectionId": o.Meta.SectionID,
			"DisplayAssetPath": o.Meta.DisplayAssetPath,
			"NewDisplayAssetPath": o.Meta.NewDisplayAssetPath,
		},
		"giftInfo": aid.JSON{
			"bIsEnabled": false,
			"forcedGiftBoxTemplateId": "",
			"purchaseRequirements": []string{},
			"giftRecordIds": []string{},
		},
		"prices": []aid.JSON{{
			"currencyType": "MtxCurrency",
			"currencySubType": "Currency",
			"regularPrice": o.Price.OriginalPrice,
			"dynamicRegularPrice": -1,
			"finalPrice": o.Price.FinalPrice,
			"basePrice": o.Price.OriginalPrice,
			"saleType": o.Price.SaleType,
			"saleExpiration": "9999-12-31T23:59:59.999Z",
		}},
		"displayAssetPath": o.Meta.DisplayAssetPath,
		"refundable": false,
		"title": o.Diplay.Title,
		"description": o.Diplay.Description,
		"shortDescription": o.Diplay.ShortDescription,
		"appStoreId": []string{},
		"dailyLimit": -1,
		"weeklyLimit": -1,
		"monthlyLimit": -1,
		"sortPriority": o.Meta.Priority,
		"catalogGroupPriority": 0,
		"filterWeight": 0,
	}
}

func (o *StorefrontCatalogOfferTypeBattlePass) GenerateFortniteBulkOffersResponse() aid.JSON {
	return aid.JSON{}
}