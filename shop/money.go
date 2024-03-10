package shop

import (
	"fmt"
	"strings"

	"github.com/ectrc/snow/aid"
)

type StorefrontCatalogOfferMetaTypeCurrency struct {
	IconSize string
	DisplayAssetPath string
	NewDisplayAssetPath string
	BannerOverride string
	CurrencyAnalyticsName string
	OriginalOffer int
	ExtraBonus int
	FeaturedImageURL string
	Priority int
}

type StorefrontCatalogOfferTypeCurrency struct {
	OfferType StorefrontCatalogOfferEnum
	Rewards []*StorefrontCatalogOfferGrant
	Price *StorefrontCatalogOfferPriceRealMoney
	Diplay *OfferDisplay
	Categories []string
	Meta *StorefrontCatalogOfferMetaTypeCurrency
}

func NewCurrencyCatalogOffer() *StorefrontCatalogOfferTypeCurrency {
	return &StorefrontCatalogOfferTypeCurrency{
		OfferType: StorefrontCatalogOfferEnumItem,
		Rewards: make([]*StorefrontCatalogOfferGrant, 0),
		Price: &StorefrontCatalogOfferPriceRealMoney{},
		Diplay: &OfferDisplay{},
		Categories: make([]string, 0),
		Meta: &StorefrontCatalogOfferMetaTypeCurrency{},
	}
}

func (o *StorefrontCatalogOfferTypeCurrency) GetOffer() *StorefrontCatalogOfferTypeCurrency {
	return o
}

func (o *StorefrontCatalogOfferTypeCurrency) GetOfferID() string {
	return fmt.Sprintf("money://%s", strings.ReplaceAll(o.Diplay.Title, " ", ""))
}

func (o *StorefrontCatalogOfferTypeCurrency) GetOfferType() StorefrontCatalogOfferEnum {
	return o.OfferType
}

func (o *StorefrontCatalogOfferTypeCurrency) GetOfferPrice() *StorefrontCatalogOfferPriceRealMoney {
	return o.Price
}

func (o *StorefrontCatalogOfferTypeCurrency) GetRewards() []*StorefrontCatalogOfferGrant {
	return o.Rewards
}

func (o *StorefrontCatalogOfferTypeCurrency) GenerateFortniteCatalogOfferResponse() aid.JSON {
	return aid.JSON{
		"offerId": o.GetOfferID(),
		"offerType": "StaticPrice",
		"devName": fmt.Sprintf("[CURRENCY] %s", o.Diplay.Title),
		"itemGrants": []aid.JSON{},
		"requirements": []aid.JSON{},
		"fulfillmentIds": []string{o.GetOfferID()},
		"categories": o.Categories,
		"metaInfo": []aid.JSON{
			{
				"key": "IconSize",
				"value": o.Meta.IconSize,
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
			{
				"Key": "CurrencyAnalyticsName",
				"Value": o.Meta.CurrencyAnalyticsName,
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
			"IconSize": o.Meta.IconSize,
			"DisplayAssetPath": o.Meta.DisplayAssetPath,
			"NewDisplayAssetPath": o.Meta.NewDisplayAssetPath,
			"BannerOverride": o.Meta.BannerOverride,
			"CurrencyAnalyticsName": o.Meta.CurrencyAnalyticsName,
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
		"bannerOverride": o.Meta.BannerOverride,
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
		"sortPriority": o.Meta.Priority,
		"catalogGroupPriority": 0,
		"filterWeight": 0,
	}
}

func (o *StorefrontCatalogOfferTypeCurrency) GenerateFortniteBulkOffersResponse() aid.JSON {
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