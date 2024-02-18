package fortnite

import (
	"math/rand"
	"regexp"

	"github.com/ectrc/snow/aid"
	"github.com/google/uuid"
)

var (
	priceLookup = map[string]map[string]int{
		"EFortRarity::Legendary": {
			"AthenaCharacter": 2000,
			"AthenaBackpack":  1500,
			"AthenaPickaxe":   1500,
			"AthenaGlider":    1800,
			"AthenaDance":     500,
			"AthenaItemWrap":  800,
		},
		"EFortRarity::Epic": {
			"AthenaCharacter": 1500,
			"AthenaBackpack":  1200,
			"AthenaPickaxe":   1200,
			"AthenaGlider":    1500,
			"AthenaDance":     800,
			"AthenaItemWrap":  800,
		},
		"EFortRarity::Rare": {
			"AthenaCharacter": 1200,
			"AthenaBackpack":  800,
			"AthenaPickaxe":   800,
			"AthenaGlider":    800,
			"AthenaDance":     500,
			"AthenaItemWrap":  600,
		},
		"EFortRarity::Uncommon": {
			"AthenaCharacter": 800,
			"AthenaBackpack":  200,
			"AthenaPickaxe":   500,
			"AthenaGlider":    500,
			"AthenaDance":     200,
			"AthenaItemWrap":  300,
		},
		"EFortRarity::Common": {
			"AthenaCharacter": 500,
			"AthenaBackpack":  200,
			"AthenaPickaxe":   500,
			"AthenaGlider":    500,
			"AthenaDance":     200,
			"AthenaItemWrap":  300,
		},
	}

	dailyItemLookup = []struct {
		Season int
		Items int
	}{
		{2, 4},
		{4, 6},
		{13, 10},
	}

	weeklySetLookup = []struct {
		Season int
		Sets int
	}{
		{2, 2},
		{4, 3},
		{11, 4},
		{13, 3},
	}
)

func price(rarity, type_ string) int {
	return priceLookup[rarity][type_]
}

func dailyItems(season int) int {
	currentValue := 4

	for _, item := range dailyItemLookup {
		if item.Season > season {
			continue
		}

		currentValue = item.Items
	}

	return currentValue
}

func weeklySets(season int) int {
	currentValue := 2

	for _, set := range weeklySetLookup {
		if set.Season > season {
			continue
		}

		currentValue = set.Sets
	}

	return currentValue
}

type FortniteCatalogSectionOffer struct {
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

func NewFortniteCatalogSectionOffer() *FortniteCatalogSectionOffer {
	return &FortniteCatalogSectionOffer{}
}

func (f *FortniteCatalogSectionOffer) GenerateID() {
	for _, item := range f.Grants {
		f.ID += item.Type.BackendValue + ":" + item.ID + ","
	}

	f.ID = "v2:/" + aid.Hash([]byte(f.ID))
}

func (f *FortniteCatalogSectionOffer) GenerateTotalPrice() {
	if !f.BundleInfo.IsBundle {
		f.TotalPrice = price(f.Grants[0].Rarity.BackendValue, f.Grants[0].Type.BackendValue)
		return
	}

	for _, item := range f.Grants {
		f.TotalPrice += price(item.Rarity.BackendValue, item.Rarity.BackendValue)
	}
}

func (f *FortniteCatalogSectionOffer) GenerateFortniteCatalogSectionOffer() aid.JSON {
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
	Offers []*FortniteCatalogSectionOffer
}

func NewFortniteCatalogSection(name string) *FortniteCatalogSection {
	return &FortniteCatalogSection{
		Name: name,
	}
}

func (f *FortniteCatalogSection) GenerateFortniteCatalogSection() aid.JSON {
	catalogEntiresResponse := []aid.JSON{}
	for _, offer := range f.Offers {
		catalogEntiresResponse = append(catalogEntiresResponse, offer.GenerateFortniteCatalogSectionOffer())
	}

	return aid.JSON{
		"name": f.Name,
		"catalogEntries": catalogEntiresResponse,
	}
}

func (f *FortniteCatalogSection) GetGroupedOffers() map[string][]*FortniteCatalogSectionOffer {
	groupedOffers := map[string][]*FortniteCatalogSectionOffer{}

	for _, offer := range f.Offers {
		if groupedOffers[offer.Meta.Category] == nil {
			groupedOffers[offer.Meta.Category] = []*FortniteCatalogSectionOffer{}
		}

		groupedOffers[offer.Meta.Category] = append(groupedOffers[offer.Meta.Category], offer)
	}

	return groupedOffers
}

type FortniteCatalog struct {
	Sections []*FortniteCatalogSection
}

func NewFortniteCatalog() *FortniteCatalog {
	return &FortniteCatalog{}
}

func (f *FortniteCatalog) GenerateFortniteCatalog() aid.JSON {
	catalogSectionsResponse := []aid.JSON{}
	for _, section := range f.Sections {
		catalogSectionsResponse = append(catalogSectionsResponse, section.GenerateFortniteCatalogSection())
	}

	return aid.JSON{
		"storefronts": catalogSectionsResponse,
		"refreshIntervalHrs": 24,
		"dailyPurchaseHrs": 24,
		"expiration": "9999-12-31T23:59:59.999Z",
	}
}

func NewRandomFortniteCatalog() *FortniteCatalog {
	aid.SetRandom(rand.New(rand.NewSource(int64(aid.Config.Fortnite.ShopSeed) + aid.CurrentDayUnix())))
	catalog := NewFortniteCatalog()

	daily := NewFortniteCatalogSection("BRDailyStorefront")
	for len(daily.Offers) < dailyItems(aid.Config.Fortnite.Season) {
		entry := newEntryFromFortniteItem(GetRandomItemWithDisplayAssetOfNotType("AthenaCharacter"), false)
		entry.Meta.SectionId = "Daily"
		daily.Offers = append(daily.Offers, entry)
	}
	catalog.Sections = append(catalog.Sections, daily)

	weekly := NewFortniteCatalogSection("BRWeeklyStorefront")
	for len(weekly.GetGroupedOffers()) < weeklySets(aid.Config.Fortnite.Season) {
		set := GetRandomSet()
		for _, item := range set.Items {
			if item.DisplayAssetPath == "" || item.DisplayAssetPath2 == "" {
				continue
			}

			entry := newEntryFromFortniteItem(item, true)
			entry.Meta.Category = set.BackendName
			entry.Meta.SectionId = "Featured"
			weekly.Offers = append(weekly.Offers, entry)
		}
	}
	catalog.Sections = append(catalog.Sections, weekly)

	return catalog
}

func newEntryFromFortniteItem(fortniteItem *FortniteItem, addAssets bool) *FortniteCatalogSectionOffer {
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

func GetOfferByOfferId(id string) *FortniteCatalogSectionOffer {
	catalog := NewRandomFortniteCatalog()

	for _, section := range catalog.Sections {
		for _, offer := range section.Offers {
			if offer.ID == id {
				return offer
			}
		}
	}

	return nil
}