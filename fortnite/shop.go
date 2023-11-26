package fortnite

import (
	"sort"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/person"
	"github.com/google/uuid"
)

var (
	Rarities = map[string]int{
		"EFortRarity::Legendary": 2000,
		"EFortRarity::Epic": 1500,
		"EFortRarity::Rare": 1200,
		"EFortRarity::Uncommon": 800,
		"EFortRarity::Common": 500,
	}
)

func GetPriceForRarity(rarity string) int {
	return Rarities[rarity]
}

type Catalog struct {
	RefreshIntervalHrs int          `json:"refreshIntervalHrs"`
	DailyPurchaseHrs int          `json:"dailyPurchaseHrs"`
	Expiration string       `json:"expiration"`
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

	names := []string{}
	for _, entry := range s.CatalogEntries {
		grantStrings := entry.Grants
		sort.Strings(grantStrings)
		for _, grant := range grantStrings {
			entry.Name += grant + "-"
		}

		names = append(names, entry.Name)
	}

	aid.PrintJSON(names)
	sort.Strings(names)

	for _, devname := range names {
		for _, entry := range s.CatalogEntries {
			grantStrings := entry.Grants
			sort.Strings(grantStrings)
			for _, grant := range grantStrings {
				entry.Name += grant + "-"
			}

			if devname == entry.Name {
				json["catalogEntries"] = append(json["catalogEntries"].([]aid.JSON), entry.GenerateResponse(p))
			}
		}
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
	Title string
	ShortDescription string
}

func NewCatalogEntry(meta ...aid.JSON) *Entry {
	return &Entry{
		ID: uuid.New().String(),
		Meta: meta,
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

func (e *Entry) TileSize(size string) *Entry {
	e.Meta = append(e.Meta, aid.JSON{
		"Key": "TileSize",
		"Value": size,
	})
	return e
}

func (e *Entry) PanelType(panel string) *Entry {
	e.Panel = panel
	return e
}

func (e *Entry) Section(sectionId string) *Entry {
	e.Meta = append(e.Meta, aid.JSON{
		"Key": "SectionId",
		"Value": sectionId,
	})
	return e
}

func (e *Entry) DisplayAsset(asset string) *Entry {
	e.DisplayAssetPath = asset
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
	sort.Strings(grantStrings)
	for _, grant := range grantStrings {
		e.Name += grant + "-"
	}

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
		"dailyLimit": -1,
		"weeklyLimit": -1,
		"monthlyLimit": -1,
		"fufillmentIds": []string{},
		"filterWeight": e.Priority,
		"appStoreId": []string{},
		"refundable": false,
		"itemGrants": []aid.JSON{},
		"metaInfo": e.Meta,
		"meta": aid.JSON{},
		"title": e.Title,
		"shortDescription": e.ShortDescription,
	}

	if e.DisplayAssetPath != "" {
		json["displayAssetPath"] = "/" + e.DisplayAssetPath
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