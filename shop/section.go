package shop

import "github.com/ectrc/snow/aid"

func NewStorefrontCatalogSection(name string, type_ StorefrontCatalogOfferEnum) *StorefrontCatalogSection {
	return &StorefrontCatalogSection{
		Name: name,
		SectionType: type_,
		Offers: make([]interface{}, 0),
	}
}

func (s *StorefrontCatalogSection) GenerateFortniteCatalogSectionResponse() aid.JSON {
	catalogEntiresResponse := []aid.JSON{}
	
	for _, entry := range s.Offers {
		switch s.SectionType {
		case StorefrontCatalogOfferEnumItem:
			s := entry.(*StorefrontCatalogOfferTypeItem)
			catalogEntiresResponse = append(catalogEntiresResponse, s.GenerateFortniteCatalogOfferResponse())
		case StorefrontCatalogOfferEnumCurrency:
			s := entry.(*StorefrontCatalogOfferTypeCurrency)
			catalogEntiresResponse = append(catalogEntiresResponse, s.GenerateFortniteCatalogOfferResponse())
		case StorefrontCatalogOfferEnumStarterKit:
			s := entry.(*StorefrontCatalogOfferTypeStarterKit)
			catalogEntiresResponse = append(catalogEntiresResponse, s.GenerateFortniteCatalogOfferResponse())
		case StorefrontCatalogOfferEnumBattlePass:
			s := entry.(*StorefrontCatalogOfferTypeBattlePass)
			catalogEntiresResponse = append(catalogEntiresResponse, s.GenerateFortniteCatalogOfferResponse())
		}
	}

	return aid.JSON{
		"name": s.Name,
		"catalogEntries": catalogEntiresResponse,
	}
}

func (s *StorefrontCatalogSection) GetGroupedOffersLength() int {
	if s.SectionType != StorefrontCatalogOfferEnumItem {
		return len(s.Offers)
	}

	newOffers := []*StorefrontCatalogOfferTypeItem{}
	for _, offer := range s.Offers {
		newOffers = append(newOffers, offer.(*StorefrontCatalogOfferTypeItem))
	}

	groupedOffers := map[string][]*StorefrontCatalogOfferTypeItem{}
	for _, offer := range newOffers {
		if _, ok := groupedOffers[offer.Categories[0]]; !ok {
			groupedOffers[offer.Categories[0]] = []*StorefrontCatalogOfferTypeItem{}
		}

		groupedOffers[offer.Categories[0]] = append(groupedOffers[offer.Categories[0]], offer)
	}

	return len(groupedOffers)
}

func (s *StorefrontCatalogSection) AddOffer(offer interface{}) {
	s.Offers = append(s.Offers, offer)
}

func (s *StorefrontCatalogSection) GetOfferByID(offerID string) (interface{}, StorefrontCatalogOfferEnum) {
	for _, offer := range s.Offers {
		switch s.SectionType {
		case StorefrontCatalogOfferEnumItem:
			o := offer.(*StorefrontCatalogOfferTypeItem)
			if o.GetOfferID() == offerID {
				return o, StorefrontCatalogOfferEnumItem
			}
		case StorefrontCatalogOfferEnumCurrency:
			o := offer.(*StorefrontCatalogOfferTypeCurrency)
			if o.GetOfferID() == offerID {
				return o, StorefrontCatalogOfferEnumCurrency
			}
		case StorefrontCatalogOfferEnumStarterKit:
			o := offer.(*StorefrontCatalogOfferTypeStarterKit)
			if o.GetOfferID() == offerID {
				return o, StorefrontCatalogOfferEnumStarterKit
			}
		case StorefrontCatalogOfferEnumBattlePass:
			o := offer.(*StorefrontCatalogOfferTypeBattlePass)
			if o.GetOfferID() == offerID {
				return o, StorefrontCatalogOfferEnumBattlePass
			}
		}
	}

	return nil, -1
}