package fortnite

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/ectrc/snow/aid"
	"github.com/google/uuid"
)

type FortniteCatalogStarterPackGrant struct {
	TemplateID string
	Quantity int
}

func NewFortniteCatalogStarterPackGrant(templateID string, quantity int) *FortniteCatalogStarterPackGrant {
	return &FortniteCatalogStarterPackGrant{
		TemplateID: templateID,
		Quantity: quantity,
	}
}

type FortniteCatalogStarterPack struct {
	ID string
	DevName string
	Grants []*FortniteCatalogStarterPackGrant
	Meta struct {
		IconSize string
		BannerOverride string
		DisplayAssetPath string
		NewDisplayAssetPath string
		OriginalOffer int
		ExtraBonus int
	}
	Price struct {
		PriceType string
		PriceToPay int
	}
	Title string
	Description string
	LongDescription string
	Priority int
	SeasonsAllowed []int
}

func NewFortniteCatalogStarterPack(price int) *FortniteCatalogStarterPack {
	return &FortniteCatalogStarterPack{
		ID: "v2:/" + aid.RandomString(32),
		Price: struct {
			PriceType string
			PriceToPay int	
		}{"RealMoney", price},
	}
}

func (f *FortniteCatalogStarterPack) GenerateFortniteCatalogStarterPackResponse() aid.JSON {
	grantsResponse := []aid.JSON{}
	for _, grant := range f.Grants {
		grantsResponse = append(grantsResponse, aid.JSON{
			"templateId": grant.TemplateID,
			"quantity": grant.Quantity,
		})
	}

	prices := []aid.JSON{}
	switch f.Price.PriceType {
	case "RealMoney":
		prices = append(prices, aid.JSON{
			"currencyType": "RealMoney",
			"currencySubType": "",
			"regularPrice": 0,
			"dynamicRegularPrice": -1,
			"finalPrice": 0,
			"saleExpiration": "9999-12-31T23:59:59.999Z",
			"basePrice": 0,
		})
	case "MtxCurrency":
		prices = append(prices, aid.JSON{
			"currencyType": "MtxCurrency",
			"currencySubType": "",
			"regularPrice": f.Price.PriceToPay,
			"dynamicRegularPrice": f.Price.PriceToPay,
			"finalPrice": f.Price.PriceToPay,
			"saleExpiration": "9999-12-31T23:59:59.999Z",
			"basePrice": f.Price.PriceToPay,
		})
	}

	return aid.JSON{
		"offerId": f.ID,
		"devName": f.DevName,
		"offerType": "StaticPrice",
		"prices": prices,
		"categories": []string{},
		"dailyLimit": -1,
		"weeklyLimit": -1,
		"monthlyLimit": -1,
		"refundable": false,
		"appStoreId": []string{
			"",
			"app-" + f.ID,
		},
		"requirements": []aid.JSON{},
		"metaInfo": []aid.JSON{
			{
				"key": "SectionId",
				"value": "LimitedTime",
			},
			{
				"key": "IconSize",
				"value": f.Meta.IconSize,
			},
			{
				"key": "BannerOverride",
				"value": f.Meta.BannerOverride,
			},
			{
				"key": "DisplayAssetPath",
				"value": f.Meta.DisplayAssetPath,
			},
			{
				"key": "NewDisplayAssetPath",
				"value": f.Meta.NewDisplayAssetPath,
			},
			{
				"key": "MtxQuantity",
				"value": f.Meta.OriginalOffer + f.Meta.ExtraBonus,
			},
			{
				"key": "MtxBonus",
				"value": f.Meta.ExtraBonus,
			},
		},
		"meta": aid.JSON{
			"IconSize": f.Meta.IconSize,
			"BannerOverride": f.Meta.BannerOverride,
			"SectionID": "LimitedTime",
			"DisplayAssetPath": f.Meta.DisplayAssetPath,
			"NewDisplayAssetPath": f.Meta.NewDisplayAssetPath,
			"MtxQuantity": f.Meta.OriginalOffer + f.Meta.ExtraBonus,
			"MtxBonus": f.Meta.ExtraBonus,
		},
		"catalogGroup": "",
		"catalogGroupPriority": 0,
		"sortPriority": f.Priority,
		"bannerOverride": f.Meta.BannerOverride,
		"title": f.Title,
		"shortDescription": "",
		"description": f.Description,
		"displayAssetPath": f.Meta.DisplayAssetPath,
		"itemGrants": []aid.JSON{},
	}
}

func (f *FortniteCatalogStarterPack) GenerateFortniteCatalogBulkOfferResponse() aid.JSON {
	return aid.JSON{
		"id": "app-" + f.ID,
		"title": f.Title,
		"description": f.Description,
		"longDescription": f.LongDescription,
		"technicalDetails": "",
		"keyImages": []aid.JSON{},
		"categories": []aid.JSON{},
		"namespace": "fn",
		"status": "ACTIVE",
		"creationDate": time.Now().Format(time.RFC3339),
		"lastModifiedDate": time.Now().Format(time.RFC3339),
		"customAttributes": aid.JSON{},
		"internalName": f.Title,
		"recurrence": "ONCE",
		"items": []aid.JSON{},
		"price": DataClient.GetLocalizedPrice("GBP", f.Price.PriceToPay),
		"currentPrice": DataClient.GetLocalizedPrice("GBP", f.Price.PriceToPay),
		"currencyCode": "GBP",
		"basePrice": DataClient.GetLocalizedPrice("USD", f.Price.PriceToPay),
		"basePriceCurrencyCode": "USD",
		"recurringPrice": 0,
		"freeDays": 0,
		"maxBillingCycles": 0,
		"seller": aid.JSON{},
		"viewableDate": time.Now().Format(time.RFC3339),
		"effectiveDate": time.Now().Format(time.RFC3339),
		"expiryDate": "9999-12-31T23:59:59.999Z",
		"vatIncluded": true,
		"isCodeRedemptionOnly": false,
		"isFeatured": false,
		"taxSkuId": "FN_Currency",
		"merchantGroup": "FN_MKT",
		"priceTier": fmt.Sprintf("%d",  DataClient.GetLocalizedPrice("USD", f.Price.PriceToPay)),
		"urlSlug": "fortnite--" + f.Title,
		"roleNamesToGrant": []aid.JSON{},
		"tags": []aid.JSON{},
		"purchaseLimit": -1,
		"ignoreOrder": false,
		"fulfillToGroup": false,
		"fraudItemType": "V-Bucks",
		"shareRevenue": false,
		"offerType": "OTHERS",
		"unsearchable": false,
		"releaseDate": time.Now().Format(time.RFC3339),
		"releaseOffer": "",
		"title4Sort": f.Title,
		"countriesBlacklist": []string{},
		"selfRefundable": false,
		"refundType": "NON_REFUNDABLE",
		"pcReleaseDate": time.Now().Format(time.RFC3339),
		"priceCalculationMode": "FIXED",
		"assembleMode": "SINGLE",
		"publisherDisplayName": "Epic Games",
		"developerDisplayName": "Epic Games",
		"visibilityType": "IS_LISTED",
		"currencyDecimals": 2,
		"allowPurchaseForPartialOwned": true,
		"shareRevenueWithUnderageAffiliates": false,
		"platformWhitelist": []string{},
		"platformBlacklist": []string{},
		"partialItemPrerequisiteCheck": false,
		"upgradeMode": "UPGRADED_WITH_PRICE_FULL",
	}
}

func (startPack *FortniteCatalogStarterPack) AddGrant(g *FortniteCatalogStarterPackGrant) {
	startPack.Grants = append(startPack.Grants, g)
}

type FortniteCatalogCurrencyOffer struct {
	ID string
	DevName string
	Price struct {
		OriginalOffer int
		ExtraBonus int
	}
	Meta struct {
		IconSize string
		CurrencyAnalyticsName string
		BannerOverride string
	}
	Title string
	Description string
	LongDescription string
	Priority int
}

func NewFortniteCatalogCurrencyOffer(original, bonus int) *FortniteCatalogCurrencyOffer {
	return &FortniteCatalogCurrencyOffer{
		ID: "v2:/"+aid.RandomString(32),
		Price: struct {
			OriginalOffer int
			ExtraBonus int
		}{original, bonus},
	}
}

func (f *FortniteCatalogCurrencyOffer) GenerateFortniteCatalogCurrencyOfferResponse() aid.JSON {
	return aid.JSON{
		"offerId": f.ID,
		"devName": f.DevName,
		"offerType": "StaticPrice",
		"prices": []aid.JSON{{
			"currencyType": "RealMoney",
			"currencySubType": "",
			"regularPrice": 0,
			"dynamicRegularPrice": -1,
			"finalPrice": 0,
			"saleExpiration": "9999-12-31T23:59:59.999Z",
			"basePrice": 0,
		}},
		"categories": []string{},
		"dailyLimit": -1,
		"weeklyLimit": -1,
		"monthlyLimit": -1,
		"refundable": false,
		"appStoreId": []string{
			"",
			"app-" + f.ID,
		},
		"requirements": []aid.JSON{},
		"metaInfo": []aid.JSON{
			{
				"key": "MtxQuantity",
				"value": f.Price.OriginalOffer + f.Price.ExtraBonus,
			},
			{
				"key": "MtxBonus",
				"value": f.Price.ExtraBonus,
			},
			{
				"key": "IconSize",
				"value": f.Meta.IconSize,
			},
			{
				"key": "BannerOverride",
				"value": f.Meta.BannerOverride,
			},
			{
				"Key": "CurrencyAnalyticsName",
				"Value": f.Meta.CurrencyAnalyticsName,
			},
		},
		"meta": aid.JSON{
			"IconSize": f.Meta.IconSize,
			"CurrencyAnalyticsName": f.Meta.CurrencyAnalyticsName,
			"BannerOverride": f.Meta.BannerOverride,
			"MtxQuantity": f.Price.OriginalOffer + f.Price.ExtraBonus,
			"MtxBonus": f.Price.ExtraBonus,
		},
		"catalogGroup": "",
		"catalogGroupPriority": 0,
		"sortPriority": f.Priority,
		"bannerOverride": f.Meta.BannerOverride,
		"title": f.Title,
		"shortDescription": "",
		"description": f.Description,
		"displayAssetPath": "/Game/Catalog/DisplayAssets/DA_" + f.Meta.CurrencyAnalyticsName + ".DA_" + f.Meta.CurrencyAnalyticsName,
		"itemGrants": []aid.JSON{},
	}
}

func (f *FortniteCatalogCurrencyOffer) GenerateFortniteCatalogBulkOfferResponse() aid.JSON{
	return aid.JSON{
		"id": "app-" + f.ID,
		"title": f.Title,
		"description": f.Description,
		"longDescription": f.LongDescription,
		"technicalDetails": "",
		"keyImages": []aid.JSON{},
		"categories": []aid.JSON{},
		"namespace": "fn",
		"status": "ACTIVE",
		"creationDate": time.Now().Format(time.RFC3339),
		"lastModifiedDate": time.Now().Format(time.RFC3339),
		"customAttributes": aid.JSON{},
		"internalName": f.Title,
		"recurrence": "ONCE",
		"items": []aid.JSON{},
		"price": DataClient.GetStorefrontCurrencyOfferPrice("GBP", f.Price.OriginalOffer + f.Price.ExtraBonus),
		"currentPrice": DataClient.GetStorefrontCurrencyOfferPrice("GBP", f.Price.OriginalOffer + f.Price.ExtraBonus),
		"currencyCode": "GBP",
		"basePrice": DataClient.GetStorefrontCurrencyOfferPrice("USD", f.Price.OriginalOffer + f.Price.ExtraBonus),
		"basePriceCurrencyCode": "USD",
		"recurringPrice": 0,
		"freeDays": 0,
		"maxBillingCycles": 0,
		"seller": aid.JSON{},
		"viewableDate": time.Now().Format(time.RFC3339),
		"effectiveDate": time.Now().Format(time.RFC3339),
		"expiryDate": "9999-12-31T23:59:59.999Z",
		"vatIncluded": true,
		"isCodeRedemptionOnly": false,
		"isFeatured": false,
		"taxSkuId": "FN_Currency",
		"merchantGroup": "FN_MKT",
		"priceTier": fmt.Sprintf("%d", DataClient.GetStorefrontCurrencyOfferPrice("USD", f.Price.OriginalOffer + f.Price.ExtraBonus)),
		"urlSlug": "fortnite--" + f.Title,
		"roleNamesToGrant": []aid.JSON{},
		"tags": []aid.JSON{},
		"purchaseLimit": -1,
		"ignoreOrder": false,
		"fulfillToGroup": false,
		"fraudItemType": "V-Bucks",
		"shareRevenue": false,
		"offerType": "OTHERS",
		"unsearchable": false,
		"releaseDate": time.Now().Format(time.RFC3339),
		"releaseOffer": "",
		"title4Sort": f.Title,
		"countriesBlacklist": []string{},
		"selfRefundable": false,
		"refundType": "NON_REFUNDABLE",
		"pcReleaseDate": time.Now().Format(time.RFC3339),
		"priceCalculationMode": "FIXED",
		"assembleMode": "SINGLE",
		"publisherDisplayName": "Epic Games",
		"developerDisplayName": "Epic Games",
		"visibilityType": "IS_LISTED",
		"currencyDecimals": 2,
		"allowPurchaseForPartialOwned": true,
		"shareRevenueWithUnderageAffiliates": false,
		"platformWhitelist": []string{},
		"platformBlacklist": []string{},
		"partialItemPrerequisiteCheck": false,
		"upgradeMode": "UPGRADED_WITH_PRICE_FULL",
	}
}

type FortniteCatalogCosmeticOffer struct {
	ID string
	Grants []*FortniteItem
	TotalPrice int
	Meta struct {
		DisplayAssetPath string
		NewDisplayAssetPath string
		SectionId string
		TileSize string
		Category string
		ProfileId string
	}
	Frontend struct {
		Title string
		Description string
		ShortDescription string
	}
	Giftable bool
	BundleInfo struct {
		IsBundle bool
		PricePercent float32
	}
}

func NewFortniteCatalogSectionOffer() *FortniteCatalogCosmeticOffer {
	return &FortniteCatalogCosmeticOffer{}
}

func (f *FortniteCatalogCosmeticOffer) GenerateID() {
	for _, item := range f.Grants {
		f.ID += item.Type.BackendValue + ":" + item.ID + ","
	}

	f.ID = "v2:/" + aid.Hash([]byte(f.ID))
}

func (f *FortniteCatalogCosmeticOffer) GenerateTotalPrice() {
	if !f.BundleInfo.IsBundle {
		f.TotalPrice = DataClient.GetStorefrontCosmeticOfferPrice(f.Grants[0].Rarity.BackendValue, f.Grants[0].Type.BackendValue)
		return
	}

	for _, item := range f.Grants {
		f.TotalPrice += DataClient.GetStorefrontCosmeticOfferPrice(item.Rarity.BackendValue, item.Rarity.BackendValue)
	}
}

func (f *FortniteCatalogCosmeticOffer) GenerateFortniteCatalogCosmeticOfferResponse() aid.JSON {
	f.GenerateTotalPrice()

	itemGrantResponse := []aid.JSON{}
	purchaseRequirementsResponse := []aid.JSON{}

	for _, item := range f.Grants {
		itemGrantResponse = append(itemGrantResponse, aid.JSON{
			"templateId": item.Type.BackendValue + ":" + item.ID,
			"quantity": 1,
		})

		purchaseRequirementsResponse = append(purchaseRequirementsResponse, aid.JSON{
			"requirementType": "DenyOnItemOwnership",
			"requiredId":	item.Type.BackendValue + ":" + item.ID,
			"minQuantity": 1,
		})
	}

	return aid.JSON{
		"devName": uuid.New().String(),
		"offerId": f.ID,
		"offerType": "StaticPrice",
		"prices": []aid.JSON{{
			"currencyType": "MtxCurrency",
			"currencySubType": "",
			"regularPrice": f.TotalPrice,
			"dynamicRegularPrice": f.TotalPrice,
			"finalPrice": f.TotalPrice,
			"basePrice": f.TotalPrice,
			"saleExpiration": "9999-12-31T23:59:59.999Z",
		}},
		"itemGrants": itemGrantResponse,
		"meta": aid.JSON{
			"TileSize": f.Meta.TileSize,
			"SectionId": f.Meta.SectionId,
			"NewDisplayAssetPath": f.Meta.NewDisplayAssetPath,
			"DisplayAssetPath": f.Meta.DisplayAssetPath,
		},
		"metaInfo": []aid.JSON{
			{
				"Key": "TileSize",
				"Value": f.Meta.TileSize,
			},
			{
				"Key": "SectionId",
				"Value": f.Meta.SectionId,
			},
			{
				"Key": "NewDisplayAssetPath",
				"Value": f.Meta.NewDisplayAssetPath,
			},
			{
				"Key": "DisplayAssetPath",
				"Value": f.Meta.DisplayAssetPath,
			},
		},
		"giftInfo": aid.JSON{
			"bIsEnabled": f.Giftable,
			"forcedGiftBoxTemplateId": "",
			"purchaseRequirements": purchaseRequirementsResponse,
			"giftRecordIds": []string{},
		},
		"purchaseRequirements": purchaseRequirementsResponse,
		"categories": []string{f.Meta.Category},
		"title": f.Frontend.Title,
		"description": f.Frontend.Description,
		"shortDescription": f.Frontend.ShortDescription,
		"displayAssetPath": f.Meta.DisplayAssetPath,
		"appStoreId": []string{},
		"fufillmentIds": []string{},
		"dailyLimit": -1,
		"weeklyLimit": -1,
		"monthlyLimit": -1,
		"sortPriority": 0,
		"catalogGroupPriority": 0,
		"filterWeight": 0,
		"refundable": true,
	}
}

type FortniteCatalogSection struct {
	Name string
	Offers []*FortniteCatalogCosmeticOffer
}

func NewFortniteCatalogSection(name string) *FortniteCatalogSection {
	return &FortniteCatalogSection{
		Name: name,
	}
}

func (f *FortniteCatalogSection) GenerateFortniteCatalogSectionResponse() aid.JSON {
	catalogEntiresResponse := []aid.JSON{}
	for _, offer := range f.Offers {
		catalogEntiresResponse = append(catalogEntiresResponse, offer.GenerateFortniteCatalogCosmeticOfferResponse())
	}

	return aid.JSON{
		"name": f.Name,
		"catalogEntries": catalogEntiresResponse,
	}
}

func (f *FortniteCatalogSection) GetGroupedOffers() map[string][]*FortniteCatalogCosmeticOffer {
	groupedOffers := map[string][]*FortniteCatalogCosmeticOffer{}

	for _, offer := range f.Offers {
		if groupedOffers[offer.Meta.Category] == nil {
			groupedOffers[offer.Meta.Category] = []*FortniteCatalogCosmeticOffer{}
		}

		groupedOffers[offer.Meta.Category] = append(groupedOffers[offer.Meta.Category], offer)
	}

	return groupedOffers
}

type FortniteCatalog struct {
	Sections []*FortniteCatalogSection
	MoneyOffers []*FortniteCatalogCurrencyOffer
	StarterPacks []*FortniteCatalogStarterPack
}

func NewFortniteCatalog() *FortniteCatalog {
	return &FortniteCatalog{
		Sections: []*FortniteCatalogSection{},
		MoneyOffers: []*FortniteCatalogCurrencyOffer{},
		StarterPacks: []*FortniteCatalogStarterPack{},
	}
}

func (f *FortniteCatalog) AddSection(section *FortniteCatalogSection) {
	f.Sections = append(f.Sections, section)
}

func (f *FortniteCatalog) AddMoneyOffer(offer *FortniteCatalogCurrencyOffer) {
	offer.Priority = -len(f.MoneyOffers)
	f.MoneyOffers = append(f.MoneyOffers, offer)
}

func (f *FortniteCatalog) AddStarterPack(pack *FortniteCatalogStarterPack) {
	pack.Priority = -len(f.StarterPacks)
	f.StarterPacks = append(f.StarterPacks, pack)
}

func (f *FortniteCatalog) GenerateFortniteCatalogResponse() aid.JSON {
	catalogSectionsResponse := []aid.JSON{}

	for _, section := range f.Sections {
		catalogSectionsResponse = append(catalogSectionsResponse, section.GenerateFortniteCatalogSectionResponse())
	}

	currencyOffersResponse := []aid.JSON{}
	for _, offer := range f.MoneyOffers {
		currencyOffersResponse = append(currencyOffersResponse, offer.GenerateFortniteCatalogCurrencyOfferResponse())
	}
	catalogSectionsResponse = append(catalogSectionsResponse, aid.JSON{
		"name": "CurrencyStorefront",
		"catalogEntries": currencyOffersResponse,
	})

	starterPacksResponse := []aid.JSON{}
	for _, pack := range f.StarterPacks {
		for _, season := range pack.SeasonsAllowed {
			if season == aid.Config.Fortnite.Season {
				starterPacksResponse = append(starterPacksResponse, pack.GenerateFortniteCatalogStarterPackResponse())
				break
			}
		}
	}
	catalogSectionsResponse = append(catalogSectionsResponse, aid.JSON{
		"name": "BRStarterKits",
		"catalogEntries": starterPacksResponse,
	})

	return aid.JSON{
		"storefronts": catalogSectionsResponse,
		"refreshIntervalHrs": 24,
		"dailyPurchaseHrs": 24,
		"expiration": "9999-12-31T23:59:59.999Z",
	}
}

func (f *FortniteCatalog) FindCosmeticOfferById(id string) *FortniteCatalogCosmeticOffer {
	for _, section := range f.Sections {
		for _, offer := range section.Offers {
			if offer.ID == id {
				return offer
			}
		}
	}

	return nil
}

func (f *FortniteCatalog) FindCurrencyOfferById(id string) *FortniteCatalogCurrencyOffer {
	for _, offer := range f.MoneyOffers {
		if offer.ID == id {
			return offer
		}
	}

	return nil
}

func (f *FortniteCatalog) FindStarterPackById(id string) *FortniteCatalogStarterPack {
	for _, pack := range f.StarterPacks {
		if pack.ID == id {
			return pack
		}
	}

	return nil
}

func NewRandomFortniteCatalog() *FortniteCatalog {
	aid.SetRandom(rand.New(rand.NewSource(int64(aid.Config.Fortnite.ShopSeed) + aid.CurrentDayUnix())))
	catalog := NewFortniteCatalog()

	daily := NewFortniteCatalogSection("BRDailyStorefront")
	for len(daily.Offers) < DataClient.GetStorefrontDailyItemCount(aid.Config.Fortnite.Season) {
		entry := newCosmeticOfferFromFortniteitem(GetRandomItemWithDisplayAssetOfNotType("AthenaCharacter"), false)
		entry.Meta.SectionId = "Daily"
		daily.Offers = append(daily.Offers, entry)
	}
	catalog.AddSection(daily)

	weekly := NewFortniteCatalogSection("BRWeeklyStorefront")
	for len(weekly.GetGroupedOffers()) < DataClient.GetStorefrontWeeklySetCount(aid.Config.Fortnite.Season) {
		set := GetRandomSet()
		for _, item := range set.Items {
			if item.DisplayAssetPath == "" || item.DisplayAssetPath2 == "" {
				continue
			}

			entry := newCosmeticOfferFromFortniteitem(item, true)
			entry.Meta.Category = set.BackendName
			entry.Meta.SectionId = "Featured"
			weekly.Offers = append(weekly.Offers, entry)
		}
	}
	catalog.AddSection(weekly)

	if aid.Config.Fortnite.EnableVBucks {
		smallCurrencyOffer := newCurrencyOfferFromName("Small Currency Pack", 1000, 0)
		smallCurrencyOffer.Meta.IconSize = "XSmall"
		smallCurrencyOffer.Meta.CurrencyAnalyticsName = "MtxPack1000"
		catalog.AddMoneyOffer(smallCurrencyOffer)

		mediumCurrencyOffer := newCurrencyOfferFromName("Medium Currency Pack", 2000, 800)
		mediumCurrencyOffer.Meta.IconSize = "Small"
		mediumCurrencyOffer.Meta.CurrencyAnalyticsName = "MtxPack2800"
		mediumCurrencyOffer.Meta.BannerOverride = "12PercentExtra"
		catalog.AddMoneyOffer(mediumCurrencyOffer)

		intermediateCurrencyOffer := newCurrencyOfferFromName("Intermediate Currency Pack", 6000, 1500)
		intermediateCurrencyOffer.Meta.IconSize = "Medium"
		intermediateCurrencyOffer.Meta.CurrencyAnalyticsName = "MtxPack7500"
		intermediateCurrencyOffer.Meta.BannerOverride = "25PercentExtra"
		catalog.AddMoneyOffer(intermediateCurrencyOffer)

		jumboCurrencyOffer := newCurrencyOfferFromName("Jumbo Currency Pack", 10000, 3500)
		jumboCurrencyOffer.Meta.IconSize = "XLarge"
		jumboCurrencyOffer.Meta.CurrencyAnalyticsName = "MtxPack13500"
		jumboCurrencyOffer.Meta.BannerOverride = "35PercentExtra"
		catalog.AddMoneyOffer(jumboCurrencyOffer)

		rogueAgentStarterPack := newStarterPackOfferFromName("The Rogue Agent Pack", 499, []*FortniteCatalogStarterPackGrant{
			NewFortniteCatalogStarterPackGrant("AthenaCharacter:CID_090_Athena_Commando_M_Tactical", 1),
			NewFortniteCatalogStarterPackGrant("AthenaBackpack:BID_030_TacticalRogue", 1),
			NewFortniteCatalogStarterPackGrant("Currency:MtxPurchased", 600),
		}...)
		rogueAgentStarterPack.Meta.DisplayAssetPath = "/Game/Catalog/DisplayAssets/DA_Featured_CID_090_Athena_Commando_M_Tactical.DA_Featured_CID_090_Athena_Commando_M_Tactical"
		rogueAgentStarterPack.SeasonsAllowed = []int{4}
		catalog.AddStarterPack(rogueAgentStarterPack)

		wingmanStarterPack := newStarterPackOfferFromName("The Wingman Pack", 499, []*FortniteCatalogStarterPackGrant{
			NewFortniteCatalogStarterPackGrant("AthenaCharacter:CID_139_Athena_Commando_M_FighterPilot", 1),
			NewFortniteCatalogStarterPackGrant("AthenaBackpack:BID_056_FighterPilot", 1),
			NewFortniteCatalogStarterPackGrant("Currency:MtxPurchased", 600),
		}...)
		wingmanStarterPack.Meta.DisplayAssetPath = "/Game/Catalog/DisplayAssets/DA_Featured_CID_139_Athena_Commando_M_FighterPilot.DA_Featured_CID_139_Athena_Commando_M_FighterPilot"
		wingmanStarterPack.SeasonsAllowed = []int{4, 5}
		catalog.AddStarterPack(wingmanStarterPack)

		aceStarterPack := newStarterPackOfferFromName("The Ace Pack", 499, []*FortniteCatalogStarterPackGrant{
			NewFortniteCatalogStarterPackGrant("AthenaCharacter:CID_195_Athena_Commando_F_Bling", 1),
			NewFortniteCatalogStarterPackGrant("AthenaBackpack:BID_101_BlingFemale", 1),
			NewFortniteCatalogStarterPackGrant("Currency:MtxPurchased", 600),
		}...)
		aceStarterPack.Meta.DisplayAssetPath = "/Game/Catalog/DisplayAssets/DA_Featured_CID_195_Athena_Commando_F_Bling.DA_Featured_CID_195_Athena_Commando_F_Bling"
		aceStarterPack.SeasonsAllowed = []int{5, 6}
		catalog.AddStarterPack(aceStarterPack)

		summitStarterPack := newStarterPackOfferFromName("The Summit Striker Pack", 499, []*FortniteCatalogStarterPackGrant{
			NewFortniteCatalogStarterPackGrant("AthenaCharacter:CID_253_Athena_Commando_M_MilitaryFashion2", 1),
			NewFortniteCatalogStarterPackGrant("AthenaBackpack:BID_134_MilitaryFashion", 1),
			NewFortniteCatalogStarterPackGrant("Currency:MtxPurchased", 600),
		}...)
		summitStarterPack.Meta.DisplayAssetPath = "/Game/Catalog/DisplayAssets/DA_Featured_CID_253_Athena_Commando_M_MilitaryFashion2.DA_Featured_CID_253_Athena_Commando_M_MilitaryFashion2"
		summitStarterPack.SeasonsAllowed = []int{6, 7}
		catalog.AddStarterPack(summitStarterPack)

		cobaltStarterPack := newStarterPackOfferFromName("The Cobalt Pack", 499, []*FortniteCatalogStarterPackGrant{
			NewFortniteCatalogStarterPackGrant("AthenaCharacter:CID_327_Athena_Commando_M_BlueMystery", 1),
			NewFortniteCatalogStarterPackGrant("AthenaBackpack:BID_203_BlueMystery", 1),
			NewFortniteCatalogStarterPackGrant("Currency:MtxPurchased", 600),
		}...)
		cobaltStarterPack.Meta.DisplayAssetPath = "/Game/Catalog/DisplayAssets/DA_Featured_CID_327_Athena_Commando_M_BlueMystery.DA_Featured_CID_327_Athena_Commando_M_BlueMystery"
		cobaltStarterPack.SeasonsAllowed = []int{7}
		catalog.AddStarterPack(cobaltStarterPack)

		lagunaStarterPack := newStarterPackOfferFromName("The Laguna Pack", 499, []*FortniteCatalogStarterPackGrant{
			NewFortniteCatalogStarterPackGrant("AthenaCharacter:CID_367_Athena_Commando_F_Tropical", 1),
			NewFortniteCatalogStarterPackGrant("AthenaBackpack:BID_231_TropicalFemale", 1),
			NewFortniteCatalogStarterPackGrant("AthenaItemWrap:Wrap_033_TropicalGirl", 1),
			NewFortniteCatalogStarterPackGrant("Currency:MtxPurchased", 600),
		}...)
		lagunaStarterPack.Meta.DisplayAssetPath = "/Game/Catalog/DisplayAssets/DA_Featured_CID_367_Athena_Commando_F_Tropical.DA_Featured_CID_367_Athena_Commando_F_Tropical"
		lagunaStarterPack.SeasonsAllowed = []int{8}
		catalog.AddStarterPack(lagunaStarterPack)

		wildeStarterPack := newStarterPackOfferFromName("The Wilde Pack", 499, []*FortniteCatalogStarterPackGrant{
			NewFortniteCatalogStarterPackGrant("AthenaCharacter:CID_420_Athena_Commando_F_WhiteTiger", 1),
			NewFortniteCatalogStarterPackGrant("AthenaBackpack:BID_277_WhiteTiger", 1),
			NewFortniteCatalogStarterPackGrant("Currency:MtxPurchased", 600),
		}...)
		wildeStarterPack.Meta.DisplayAssetPath = "/Game/Catalog/DisplayAssets/DA_Featured_CID_420_Athena_Commando_F_WhiteTiger.DA_Featured_CID_420_Athena_Commando_F_WhiteTiger"
		wildeStarterPack.SeasonsAllowed = []int{9}
		catalog.AddStarterPack(wildeStarterPack)

		redStrikePack := newStarterPackOfferFromName("The Red Strike Pack", 499, []*FortniteCatalogStarterPackGrant{
			NewFortniteCatalogStarterPackGrant("AthenaCharacter:CID_384_Athena_Commando_M_StreetAssassin", 1),
			NewFortniteCatalogStarterPackGrant("AthenaBackpack:BID_247_StreetAssassin", 1),
			NewFortniteCatalogStarterPackGrant("Currency:MtxPurchased", 600),
		}...)
		redStrikePack.Meta.DisplayAssetPath = "/Game/Catalog/DisplayAssets/DA_Featured_CID_384_Athena_Commando_M_StreetAssasin.DA_Featured_CID_384_Athena_Commando_M_StreetAssasin"
		redStrikePack.SeasonsAllowed = []int{10}
		catalog.AddStarterPack(redStrikePack)

		// Below is an example of a custom starter pack
		// Uncomment to use.
		// snowCustomPack := newStarterPackOfferFromName("Snow Gift", 0, []*FortniteCatalogStarterPackGrant{
		// 	NewFortniteCatalogStarterPackGrant("AthenaCharacter:CID_384_Athena_Commando_M_StreetAssassin", 1),
		// }...)
		// snowCustomPack.Meta.DisplayAssetPath = "/Game/Catalog/DisplayAssets/DA_Featured_CID_TBD_Athena_Commando_M_RaptorArcticCamo_Bundle.DA_Featured_CID_TBD_Athena_Commando_M_RaptorArcticCamo_Bundle"
		// snowCustomPack.SeasonsAllowed = []int{1,2,3,4,5,6,7,8,9,10}
		// snowCustomPack.Meta.OriginalOffer = 1000
		// snowCustomPack.Meta.ExtraBonus = 500
		// snowCustomPack.Description = ""
		// snowCustomPack.LongDescription = "Thank you for using Snow! Here's a special offer for you!"
		// catalog.AddStarterPack(snowCustomPack)
	}
	return catalog
}

func newCosmeticOfferFromFortniteitem(fortniteItem *FortniteItem, addAssets bool) *FortniteCatalogCosmeticOffer {
	displayAsset := regexp.MustCompile(`[^/]+$`).FindString(fortniteItem.DisplayAssetPath)

	entry := NewFortniteCatalogSectionOffer()
	entry.Meta.TileSize = "Small"
	if fortniteItem.Type.BackendValue == "AthenaCharacter" {
		entry.Meta.TileSize = "Normal"
	}
	if addAssets {
		entry.Meta.NewDisplayAssetPath = "/Game/Catalog/NewDisplayAssets/" + fortniteItem.DisplayAssetPath2 + "." + fortniteItem.DisplayAssetPath2
		if displayAsset != "" {
			entry.Meta.DisplayAssetPath = "/Game/Catalog/DisplayAssets/" + displayAsset + "." + displayAsset
		}
	}
	entry.Meta.ProfileId = "athena"
	entry.Giftable = true
	entry.Grants = append(entry.Grants, fortniteItem)
	entry.GenerateTotalPrice()
	entry.GenerateID()

	return entry
}

func newCurrencyOfferFromName(name string, original, bonus int) *FortniteCatalogCurrencyOffer {
	formattedPrice := aid.FormatNumber(original + bonus)
	offer := NewFortniteCatalogCurrencyOffer(original, bonus)
	offer.Meta.IconSize = "Small"
	offer.Meta.CurrencyAnalyticsName = name
	offer.DevName = name
	offer.Title = formattedPrice + " V-Bucks"
	offer.Description = "Buy " + formattedPrice + " Fortnite V-Bucks, the in-game currency that can be spent in Fortnite Battle Royale and Creative modes. You can purchase new customization items like Outfits, Gliders, Pickaxes, Emotes, Wraps and the latest season's Battle Pass! Gliders and Contrails may not be used in Save the World mode."
	offer.LongDescription = "Buy " + formattedPrice + " Fortnite V-Bucks, the in-game currency that can be spent in Fortnite Battle Royale and Creative modes. You can purchase new customization items like Outfits, Gliders, Pickaxes, Emotes, Wraps and the latest season's Battle Pass! Gliders and Contrails may not be used in Save the World mode.\n\nAll V-Bucks purchased on the Epic Games Store are not redeemable or usable on Nintendo Switch™."

	return offer
}

func newStarterPackOfferFromName(name string, totalPrice int, grants ...*FortniteCatalogStarterPackGrant) *FortniteCatalogStarterPack {
	mainString := "Jump into Fortnite Battle Royale with the " + strings.ReplaceAll(name, "The ", "") + ". Includes:\n\n- 600 V-Bucks"

	for _, grant := range grants {
		fortniteItem := DataClient.FortniteItems[strings.Split(grant.TemplateID, ":")[1]]
		if fortniteItem != nil {
			mainString += "\n- " + fortniteItem.Name + " " + fortniteItem.Type.DisplayValue + " - Battle Royale Only"
		}
	}

	offer := NewFortniteCatalogStarterPack(totalPrice)
	offer.DevName = name + "StarterPack"
	offer.Title = name
	offer.Description = mainString
	offer.LongDescription = mainString + "\n\nV-Bucks are an in-game currency that can be spent in both the Battle Royale PvP mode and the Save the World PvE campaign. In Battle Royale, you can use V-bucks to purchase new customization items like outfits, emotes, pickaxes, gliders, and more! In Save the World you can purchase Llama Pinata card packs that contain weapon, trap and gadget schematics as well as new Heroes and more! \n\nNote: Items do not transfer between the Battle Royale mode and the Save the World campaign."
	offer.Meta.OriginalOffer = 500
	offer.Meta.ExtraBonus = 100

	for _, grant := range grants {
		offer.AddGrant(grant)
	}

	return offer
}