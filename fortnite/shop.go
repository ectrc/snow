package fortnite

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"github.com/ectrc/snow/aid"
	"github.com/google/uuid"
)


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
}

func NewFortniteCatalog() *FortniteCatalog {
	return &FortniteCatalog{
		Sections: []*FortniteCatalogSection{},
		MoneyOffers: []*FortniteCatalogCurrencyOffer{},
	}
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

func NewRandomFortniteCatalog() *FortniteCatalog {
	aid.SetRandom(rand.New(rand.NewSource(int64(aid.Config.Fortnite.ShopSeed) + aid.CurrentDayUnix())))
	catalog := NewFortniteCatalog()

	daily := NewFortniteCatalogSection("BRDailyStorefront")
	for len(daily.Offers) < DataClient.GetStorefrontDailyItemCount(aid.Config.Fortnite.Season) {
		entry := newCosmeticOfferFromFortniteitem(GetRandomItemWithDisplayAssetOfNotType("AthenaCharacter"), false)
		entry.Meta.SectionId = "Daily"
		daily.Offers = append(daily.Offers, entry)
	}
	catalog.Sections = append(catalog.Sections, daily)

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
	catalog.Sections = append(catalog.Sections, weekly)

	if aid.Config.Fortnite.EnableVBucks {
		smallCurrencyOffer := newCurrencyOfferFromName("Small Currency Pack", 1000, 0)
		smallCurrencyOffer.Meta.IconSize = "XSmall"
		smallCurrencyOffer.Meta.CurrencyAnalyticsName = "MtxPack1000"
		smallCurrencyOffer.Priority = -len(catalog.MoneyOffers)
		catalog.MoneyOffers = append(catalog.MoneyOffers, smallCurrencyOffer)

		mediumCurrencyOffer := newCurrencyOfferFromName("Medium Currency Pack", 2000, 800)
		mediumCurrencyOffer.Meta.IconSize = "Small"
		mediumCurrencyOffer.Meta.CurrencyAnalyticsName = "MtxPack2800"
		mediumCurrencyOffer.Meta.BannerOverride = "12PercentExtra"
		mediumCurrencyOffer.Priority = -len(catalog.MoneyOffers)
		catalog.MoneyOffers = append(catalog.MoneyOffers, mediumCurrencyOffer)

		intermediateCurrencyOffer := newCurrencyOfferFromName("Intermediate Currency Pack", 4000, 1000)
		intermediateCurrencyOffer.Meta.IconSize = "Medium"
		intermediateCurrencyOffer.Meta.CurrencyAnalyticsName = "MtxPack5000"
		intermediateCurrencyOffer.Meta.BannerOverride = "25PercentExtra"
		intermediateCurrencyOffer.Priority = -len(catalog.MoneyOffers)
		catalog.MoneyOffers = append(catalog.MoneyOffers, intermediateCurrencyOffer)

		jumboCurrencyOffer := newCurrencyOfferFromName("Jumbo Currency Pack", 10000, 3500)
		jumboCurrencyOffer.Meta.IconSize = "XLarge"
		jumboCurrencyOffer.Meta.CurrencyAnalyticsName = "MtxPack13500"
		jumboCurrencyOffer.Meta.BannerOverride = "35PercentExtra"
		jumboCurrencyOffer.Priority = -len(catalog.MoneyOffers)
		catalog.MoneyOffers = append(catalog.MoneyOffers, jumboCurrencyOffer)
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