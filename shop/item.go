package shop

import (
	"fmt"

	"github.com/ectrc/snow/aid"
)

type StorefrontCatalogOfferMetaTypeItem struct {
	TileSize string
	SectionID string
	DisplayAssetPath string
	NewDisplayAssetPath string
	BannerOverride string
	Giftable bool
	Refundable bool
}

type StorefrontCatalogOfferTypeItem struct {
	OfferID string
	OfferType StorefrontCatalogOfferEnum
	Rewards []*StorefrontCatalogOfferGrant
	Price *StorefrontCatalogOfferPriceMtxCurrency
	Diplay *OfferDisplay
	Categories []string
	Meta *StorefrontCatalogOfferMetaTypeItem
}

func NewItemCatalogOffer() *StorefrontCatalogOfferTypeItem {
	return &StorefrontCatalogOfferTypeItem{
		OfferID: aid.RandomString(32),
		OfferType: StorefrontCatalogOfferEnumItem,
		Rewards: make([]*StorefrontCatalogOfferGrant, 0),
		Price: &StorefrontCatalogOfferPriceMtxCurrency{},
		Diplay: &OfferDisplay{},
		Categories: make([]string, 0),
		Meta: &StorefrontCatalogOfferMetaTypeItem{},
	}
}

func (o *StorefrontCatalogOfferTypeItem) GetOffer() *StorefrontCatalogOfferTypeItem {
	return o
}

func (o *StorefrontCatalogOfferTypeItem) GetOfferID() string {
	return o.OfferID
}

func (o *StorefrontCatalogOfferTypeItem) GetOfferType() StorefrontCatalogOfferEnum {
	return o.OfferType
}

func (o *StorefrontCatalogOfferTypeItem) GetOfferPrice() *StorefrontCatalogOfferPriceMtxCurrency {
	return o.Price
}

func (o *StorefrontCatalogOfferTypeItem) GetRewards() []*StorefrontCatalogOfferGrant {
	return o.Rewards
}

func (o *StorefrontCatalogOfferTypeItem) GenerateFortniteCatalogOfferResponse() aid.JSON {
	itemGrantResponse := []aid.JSON{}
	purchaseRequirementsResponse := []aid.JSON{}
	developerNameResponse := "[ITEM]"

	for _, reward := range o.Rewards {
		itemGrantResponse = append(itemGrantResponse, aid.JSON{
			"templateId": reward.TemplateID,
			"quantity": reward.Quantity,
		})
	
		purchaseRequirementsResponse = append(purchaseRequirementsResponse, aid.JSON{
			"requirementType": "DenyOnItemOwnership",
			"requiredId":	reward.TemplateID,
			"minQuantity": 1,
		})

		developerNameResponse += fmt.Sprintf(" %dx %s", reward.Quantity, reward.TemplateID)
	}

	return aid.JSON{
		"offerId": o.OfferID,
		"offerType": "StaticPrice",
		"devName": fmt.Sprintf("%s for %d MtxCurrency", developerNameResponse, o.Price.OriginalPrice),
		"itemGrants": itemGrantResponse,
		"requirements": purchaseRequirementsResponse,
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
			{
				"Key": "BannerOverride",
				"Value": o.Meta.BannerOverride,
			},
		},
		"meta": aid.JSON{
			"TileSize": o.Meta.TileSize,
			"SectionId": o.Meta.SectionID,
			"DisplayAssetPath": o.Meta.DisplayAssetPath,
			"NewDisplayAssetPath": o.Meta.NewDisplayAssetPath,
			"BannerOverride": o.Meta.BannerOverride,
		},
		"giftInfo": aid.JSON{
			"bIsEnabled": o.Meta.Giftable,
			"forcedGiftBoxTemplateId": "",
			"purchaseRequirements": purchaseRequirementsResponse,
			"giftRecordIds": []string{},
		},
		"prices": []aid.JSON{{
			"currencyType": "MtxCurrency",
			"currencySubType": "Currency",
			"regularPrice": o.Price.OriginalPrice,
			"dynamicRegularPrice": -1,
			"finalPrice": o.Price.FinalPrice,
			"basePrice": o.Price.OriginalPrice,
			"saleExpiration": "9999-12-31T23:59:59.999Z",
		}},
		"bannerOverride": o.Meta.BannerOverride,
		"displayAssetPath": o.Meta.DisplayAssetPath,
		"refundable": o.Meta.Refundable,
		"title": o.Diplay.Title,
		"description": o.Diplay.Description,
		"shortDescription": o.Diplay.ShortDescription,
		"appStoreId": []string{},
		"fulfillmentIds": []string{},
		"dailyLimit": -1,
		"weeklyLimit": -1,
		"monthlyLimit": -1,
		"sortPriority": 0,
		"catalogGroupPriority": 0,
		"filterWeight": 0,
	}
}

func (o *StorefrontCatalogOfferTypeItem) GenerateFortniteBulkOffersResponse() aid.JSON {
	return aid.JSON{}
}