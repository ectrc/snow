package fortnite

import (
	"strings"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/person"
	"github.com/google/uuid"
)

var (
	Rarities = map[string]map[string]int{
		"EFortRarity::Legendary": {
			"AthenaCharacter": 2000,
			"AthenaBackpack": 1500,
			"AthenaPickaxe": 1500,
			"AthenaGlider": 1800,
			"AthenaDance": 500,
			"AthenaItemWrap": 800,
		},
		"EFortRarity::Epic": {
			"AthenaCharacter": 1500,
			"AthenaBackpack": 1200,
			"AthenaPickaxe": 1200,
			"AthenaGlider": 1500,
			"AthenaDance": 800,
			"AthenaItemWrap": 800,
		},
		"EFortRarity::Rare": {
			"AthenaCharacter": 1200,
			"AthenaBackpack": 800,
			"AthenaPickaxe": 800,
			"AthenaGlider": 800,
			"AthenaDance": 500,
			"AthenaItemWrap": 600,
		},
		"EFortRarity::Uncommon": {
			"AthenaCharacter": 800,
			"AthenaBackpack": 200,
			"AthenaPickaxe": 500,
			"AthenaGlider": 500,
			"AthenaDance": 200,
			"AthenaItemWrap": 300,
		},
		"EFortRarity::Common": {
			"AthenaCharacter": 500,
			"AthenaBackpack": 200,
			"AthenaPickaxe": 500,
			"AthenaGlider": 500,
			"AthenaDance": 200,
			"AthenaItemWrap": 300,
		},
	}
	StaticCatalog = NewCatalog() 
)

func GetPriceForRarity(rarity string, backendType string) int {
	return Rarities[rarity][backendType]
}

type Catalog struct {
	RefreshIntervalHrs int `json:"refreshIntervalHrs"`
	DailyPurchaseHrs int `json:"dailyPurchaseHrs"`
	Expiration string `json:"expiration"`
	Storefronts []Storefront `json:"storefronts"`
}

func NewCatalog() *Catalog {
	return &Catalog{
		RefreshIntervalHrs: 24,
		DailyPurchaseHrs: 24,
		Expiration: aid.TimeEndOfDay(),
		Storefronts: []Storefront{},
	}
}

func (c *Catalog) Add(storefront *Storefront) {
	c.Storefronts = append(c.Storefronts, *storefront)
}

func (c *Catalog) GenerateFortniteCatalog(p *person.Person) aid.JSON {
	json := aid.JSON{
		"refreshIntervalHrs": c.RefreshIntervalHrs,
		"dailyPurchaseHrs": c.DailyPurchaseHrs,
		"expiration": c.Expiration,
		"storefronts": []aid.JSON{},
	}

	for _, storefront := range c.Storefronts {
		json["storefronts"] = append(json["storefronts"].([]aid.JSON), storefront.GenerateResponse(p))
	}

	return json
}

func (c *Catalog) CheckIfOfferIsDuplicate(entry Entry) bool {
	for _, storefront := range c.Storefronts {
		for _, catalogEntry := range storefront.CatalogEntries {
			if catalogEntry.Grants[0] == entry.Grants[0] {
				return true
			}
		}
	}

	return false
}

func (c *Catalog) GetOfferById(id string) *Entry {
	for _, storefront := range c.Storefronts {
		for _, catalogEntry := range storefront.CatalogEntries {
			if catalogEntry.ID == id {
				return &catalogEntry
			}
		}
	}

	return nil
}

type Storefront struct {
	Name string `json:"name"`
	CatalogEntries []Entry `json:"catalogEntries"`
}

func NewStorefront(name string) *Storefront {
	return &Storefront{
		Name: name,
		CatalogEntries: []Entry{},
	}
}

func (s *Storefront) Add(entry Entry) {
	s.CatalogEntries = append(s.CatalogEntries, entry)
}

func (s *Storefront) GenerateResponse(p *person.Person) aid.JSON {
	json := aid.JSON{
		"name": s.Name,
		"catalogEntries": []aid.JSON{},
	}

	for _, entry := range s.CatalogEntries {
		json["catalogEntries"] = append(json["catalogEntries"].([]aid.JSON), entry.GenerateResponse(p))
	}

	return json
}

type Entry struct {
	ID string
	Name string
	Price int
	Meta []aid.JSON
	Panel string
	Priority int
	Grants []string
	DisplayAssetPath string
	NewDisplayAssetPath string
	Title string
	ShortDescription string
	ProfileType string
}

func NewCatalogEntry(profile string, meta ...aid.JSON) *Entry {
	return &Entry{
		ID: uuid.New().String(),
		Meta: meta,
		ProfileType: profile,
	}
}

func (e *Entry) AddGrant(templateId string) *Entry {
	e.Grants = append(e.Grants, templateId)
	return e
}

func (e *Entry) AddMeta(key string, value interface{}) *Entry {
	e.Meta = append(e.Meta, aid.JSON{
		"Key": key,
		"Value": value,
	})
	return e
}

func (e *Entry) SetTileSize(size string) *Entry {
	e.Meta = append(e.Meta, aid.JSON{
		"Key": "TileSize",
		"Value": size,
	})
	return e
}

func (e *Entry) SetPanel(panel string) *Entry {
	e.Panel = panel
	return e
}

func (e *Entry) SetSection(sectionId string) *Entry {
	for _, m := range e.Meta {
		if m["Key"] == "SectionId" {
			m["Value"] = sectionId
			return e
		}
	}

	e.Meta = append(e.Meta, aid.JSON{
		"Key": "SectionId",
		"Value": sectionId,
	})
	return e
}

func (e *Entry) SetDisplayAsset(asset string) *Entry {
	displayAsset := "DAv2_Featured_" + asset
	e.DisplayAssetPath = "/Game/Catalog/DisplayAssets/" + displayAsset + "." + displayAsset
	return e
}

func (e *Entry) SetNewDisplayAsset(asset string) *Entry {
	e.NewDisplayAssetPath = "/Game/Catalog/NewDisplayAssets/" + asset + "." + asset
	return e
}

func (e *Entry) SetDisplayAssetPath(path string) *Entry {
	paths := strings.Split(path, "/")
	id := paths[len(paths)-1]

	e.DisplayAssetPath = "/Game/Catalog/DisplayAssets/" + id + "." + id
	return e
}

func (e *Entry) SetNewDisplayAssetPath(path string) *Entry {
	e.NewDisplayAssetPath = path
	return e
}

func (e *Entry) SetTitle(title string) *Entry {
	e.Title = title
	return e
}

func (e *Entry) SetShortDescription(description string) *Entry {
	e.ShortDescription = description
	return e
}

func (e *Entry) SetPrice(price int) *Entry {
	e.Price = price
	return e
}

func (e *Entry) GenerateResponse(p *person.Person) aid.JSON {
	grantStrings := e.Grants
	for _, grant := range grantStrings {
		e.Name += grant + "-"
	}

	if e.NewDisplayAssetPath == "" && len(e.Grants) != 0 {
		safeTemplateId := strings.ReplaceAll(strings.Split(e.Grants[0], ":")[1], "Athena_Commando_", "")
		newDisplayAsset := "DAv2_" + safeTemplateId
		e.NewDisplayAssetPath = "/Game/Catalog/NewDisplayAssets/" + newDisplayAsset + "." + newDisplayAsset
	}
	e.AddMeta("NewDisplayAssetPath", e.NewDisplayAssetPath)

	if e.DisplayAssetPath == "" && len(e.Grants) != 0 {	
		displayAsset := "DA_Featured_" + strings.Split(e.Grants[0], ":")[1]
		e.DisplayAssetPath = "/Game/Catalog/DisplayAssets/" + displayAsset + "." + displayAsset
	}
	e.AddMeta("DisplayAssetPath", e.DisplayAssetPath)

	json := aid.JSON{
		"offerId": e.ID,
		"devName": e.Name,
		"offerType": "StaticPrice",
		"prices": []aid.JSON{
			{
				"currencyType": "MtxCurrency",
				"currencySubType": "Currency",
				"regularPrice": e.Price,
				"dynamicRegularPrice": e.Price,
				"finalPrice": e.Price,
				"basePrice": e.Price,
				"saleExpiration": aid.TimeEndOfDay(),
			},
		},
		"categories": []string{},
		"catalogGroupPriority": 0,
		"sortPriority": e.Priority,
		"dailyLimit": -1,
		"weeklyLimit": -1,
		"monthlyLimit": -1,
		"fufillmentIds": []string{},
		"filterWeight": 0.0,
		"appStoreId": []string{},
		"refundable": false,
		"itemGrants": []aid.JSON{},
		"metaInfo": e.Meta,
		"meta": aid.JSON{},
		"displayAssetPath": e.DisplayAssetPath,
		"title": e.Title,
		"shortDescription": e.ShortDescription,
	}
	grants := []aid.JSON{}
	requirements := []aid.JSON{}
	meta := []aid.JSON{}

	for _, templateId := range e.Grants {
		grants = append(grants, aid.JSON{
			"templateId": templateId,
			"quantity": 1,
		})

		if item := p.AthenaProfile.Items.GetItemByTemplateID(templateId); item != nil {
			requirements = append(requirements, aid.JSON{
				"requirementType": "DenyOnItemOwnership",
				"requiredId": item.ID,
				"minQuantity": 1,
			})
		}
	}

	for _, m := range e.Meta {
		meta = append(meta, m)
		json["meta"].(aid.JSON)[m["Key"].(string)] = m["Value"]
	}

	if e.Panel != "" {
		json["categories"] = []string{e.Panel}
	}

	json["itemGrants"] = grants
	json["requirements"] = requirements
	json["metaInfo"] = meta

	return json
}

func GenerateRandomStorefront() {
	storefront := NewCatalog()

	daily := NewStorefront("BRDailyStorefront")
	weekly := NewStorefront("BRWeeklyStorefront")

	for i := 0; i < 4; i++ {
		if aid.Config.Fortnite.Season < 14 {
			break
		}

		item := Cosmetics.GetRandomItemByType("AthenaCharacter")
		entry := NewCatalogEntry("athena")
		entry.SetSection("Daily")

		if item.DisplayAssetPath2 == "" {
			i--
			continue
		}
		entry.SetNewDisplayAsset(item.DisplayAssetPath2)

		if item.DisplayAssetPath != "" {
			entry.SetDisplayAssetPath(item.DisplayAssetPath)
		}
		entry.SetPrice(GetPriceForRarity(item.Rarity.BackendValue, item.Type.BackendValue))
		entry.AddGrant(item.Type.BackendValue + ":" + item.ID)
		entry.SetTileSize("Normal")
		entry.Priority = 1

		if item.Backpack != "" {
			entry.AddGrant("AthenaBackpack:" + item.Backpack)
		}

		if storefront.CheckIfOfferIsDuplicate(*entry) {
			continue
		}

		daily.Add(*entry)
	}
	
	for i := 0; i < 6; i++ {
		item := Cosmetics.GetRandomItemByNotType("AthenaCharacter")
		entry := NewCatalogEntry("athena")
		entry.SetSection("Daily")

		if item.DisplayAssetPath2 == "" {
			i--
			continue
		}
		entry.SetNewDisplayAsset(item.DisplayAssetPath2)
		
		if item.DisplayAssetPath != "" {
			entry.SetDisplayAssetPath(item.DisplayAssetPath)
		}
		entry.SetPrice(GetPriceForRarity(item.Rarity.BackendValue, item.Type.BackendValue))
		entry.AddGrant(item.Type.BackendValue + ":" + item.ID)
		entry.SetTileSize("Small")

		if storefront.CheckIfOfferIsDuplicate(*entry) {
			continue
		}

		daily.Add(*entry)
	}

	minimumItems := 8
	if aid.Config.Fortnite.Season < 11 {
		minimumItems = 3
	}

	minimumSets := 4
	if aid.Config.Fortnite.Season < 11 {
		minimumSets = 3
	}

	setsAdded := 0
	for len(weekly.CatalogEntries) < minimumItems || setsAdded < minimumSets {
		set := Cosmetics.GetRandomSet()
		
		itemsAdded := 0
		itemsToAdd := []*Entry{}
		for _, item := range set.Items {
			entry := NewCatalogEntry("athena")
			entry.SetSection("Featured")
			entry.SetPanel(set.BackendName)

			if item.DisplayAssetPath2 == "" {
				continue
			}
			entry.SetNewDisplayAsset(item.DisplayAssetPath2)

			if item.Type.BackendValue == "AthenaCharacter" {
				entry.SetTileSize("Normal")
				if aid.Config.Fortnite.Season < 14 {
					itemsAdded += 1
				} else {
					itemsAdded += 2
				}
				entry.Priority = 1
			} else {
				entry.SetTileSize("Small")
				itemsAdded += 1
			}

			if item.DisplayAssetPath != "" {
				entry.SetDisplayAssetPath(item.DisplayAssetPath)
			}
			entry.SetPrice(GetPriceForRarity(item.Rarity.BackendValue, item.Type.BackendValue))
			entry.AddGrant(item.Type.BackendValue + ":" + item.ID)

			itemsToAdd = append(itemsToAdd, entry)
		}

		if itemsAdded % 2 != 0 {
			itemsToAdd = itemsToAdd[:len(itemsToAdd)-1]
		}

		for _, entry := range itemsToAdd {
			if storefront.CheckIfOfferIsDuplicate(*entry) {
				continue
			}

			weekly.Add(*entry)
		}

		setsAdded++
	}

	storefront.Add(daily)
	storefront.Add(weekly)
	
	StaticCatalog = storefront
	aid.Print("(snow) generated random storefront")
}