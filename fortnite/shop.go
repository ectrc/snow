package fortnite

import (
	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/person"
)

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

	for _, entry := range s.CatalogEntries {
		json["catalogEntries"] = append(json["catalogEntries"].([]aid.JSON), entry.GenerateResponse(p))
	}

	return json
}

type Entry struct {
	Price int
	ID string
	Name string
	Title string
	Description string
	Type string
	Meta []aid.JSON
	Panel string
	Priority int
	Asset string
	Grants []string
	IsBundle bool
	BundleMeta BundleMeta
}

func NewItemEntry(id string, name string, price int) *Entry {
	return &Entry{
		Price: price,
		ID: id,
		Name: name,
		Type: "StaticPrice",
	}
}

func NewBundleEntry(id string, name string, price int) *Entry {
	return &Entry{
		Price: price,
		ID: id,
		Name: name,
		Type: "DynamicBundle",
		IsBundle: true,
		BundleMeta: BundleMeta{
			FloorPrice: price,
			RegularBasePrice: price,
			DiscountedBasePrice: price,
		},
	}
}

type BundleMeta struct {
	FloorPrice int
	RegularBasePrice int
	DiscountedBasePrice int
	DisplayType string // "AmountOff" or "PercentOff"
	BundleItems []BundleItem
}

type BundleItem struct {
	TemplateID string
	RegularPrice int
	DiscountedPrice int
	AlreadyOwnedPriceReduction int
}

func NewBundleItem(templateId string, regularPrice int, discountedPrice int, alreadyOwnedPriceReduction int) *BundleItem {
	return &BundleItem{
		TemplateID: templateId,
		RegularPrice: regularPrice,
		DiscountedPrice: discountedPrice,
		AlreadyOwnedPriceReduction: alreadyOwnedPriceReduction,
	}
}

func (e *Entry) AddGrant(templateId string) *Entry {
	e.Grants = append(e.Grants, templateId)
	return e
}

func (e *Entry) AddBundleGrant(B BundleItem) *Entry {
	e.BundleMeta.BundleItems = append(e.BundleMeta.BundleItems, B)
	return e
}

func (e *Entry) AddMeta(key string, value interface{}) *Entry {
	e.Meta = append(e.Meta, aid.JSON{
		"Key": key,
		"Value": value,
	})
	return e
}

func (e *Entry) GenerateResponse(p *person.Person) aid.JSON {
	json := aid.JSON{
		"offerId": e.ID,
		"devName": e.Name,
		"offerType": e.Type,
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
		"catalogGroupPriority": e.Priority,
		"dailyLimit": -1,
		"weeklyLimit": -1,
		"monthlyLimit": -1,
		"fufillmentIds": []string{},
		"filterWeight": 0,
		"appStoreId": []string{},
		"refundable": false,
		"itemGrants": []aid.JSON{},
		"metaInfo": e.Meta,
		"meta": aid.JSON{},
		"displayAssetPath": e.Asset,
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

	if e.IsBundle {
		json["dynamicBundleInfo"] = aid.JSON{
			"discountedBasePrice": e.BundleMeta.DiscountedBasePrice,
			"regularBasePrice": e.BundleMeta.RegularBasePrice,
			"floorPrice": e.BundleMeta.FloorPrice,
			"currencyType": "MtxCurrency",
			"currencySubType": "Currency",
			"displayType": "AmountOff",
			"bundleItems": []aid.JSON{},
		}

		for _, bundleItem := range e.BundleMeta.BundleItems {
			json["prices"] = []aid.JSON{}

			json["dynamicBundleInfo"].(aid.JSON)["bundleItems"] = append(json["dynamicBundleInfo"].(aid.JSON)["bundleItems"].([]aid.JSON), aid.JSON{
				"regularPrice": bundleItem.RegularPrice,
				"discountedPrice": bundleItem.DiscountedPrice,
				"alreadyOwnedPriceReduction": bundleItem.AlreadyOwnedPriceReduction,
				"item": aid.JSON{
					"templateId": bundleItem.TemplateID,
					"quantity": 1,
				},
			})

			grants = append(grants, aid.JSON{
				"templateId": bundleItem.TemplateID,
				"quantity": 1,
			})

			if item := p.AthenaProfile.Items.GetItemByTemplateID(bundleItem.TemplateID); item != nil {
				requirements = append(requirements, aid.JSON{
					"requirementType": "DenyOnItemOwnership",
					"requiredId": item.ID,
					"minQuantity": 1,
				})
			}
		}
	}

	json["itemGrants"] = grants
	json["requirements"] = requirements
	json["metaInfo"] = meta

	return json
}