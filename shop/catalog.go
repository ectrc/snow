package shop

import "github.com/ectrc/snow/aid"

type StorefrontCatalog struct {
	Sections []*StorefrontCatalogSection
}

func NewStorefrontCatalog() *StorefrontCatalog {
	return &StorefrontCatalog{
		Sections: make([]*StorefrontCatalogSection, 0),
	}
}

func (c *StorefrontCatalog) AddSection(section *StorefrontCatalogSection) {
	c.Sections = append(c.Sections, section)
}
func (c *StorefrontCatalog) AddSections(sections ...*StorefrontCatalogSection) {
	c.Sections = append(c.Sections, sections...)
}

func (c *StorefrontCatalog) GetOfferByID(offerID string) (interface{}, StorefrontCatalogOfferEnum) {
	for _, section := range c.Sections {
		found, type_ := section.GetOfferByID(offerID)
		if found != nil {
			return found, type_
		}
	}

	return nil, -1
}

func (c *StorefrontCatalog) GenerateFortniteCatalogResponse() aid.JSON {
	sectionsResponse := []aid.JSON{}

	for _, section := range c.Sections {
		sectionsResponse = append(sectionsResponse, section.GenerateFortniteCatalogSectionResponse())
	}

	return aid.JSON{
		"storefronts": sectionsResponse,
		"refreshIntervalHrs": 24,
		"dailyPurchaseHrs": 24,
		"expiration": "9999-12-31T23:59:59.999Z",
	}
}