package shop

import "github.com/ectrc/snow/aid"

type ShopGrantProfileType string
const ShopGrantProfileTypeAthena ShopGrantProfileType = "athena"
const ShopGrantProfileTypeCommonCore ShopGrantProfileType = "common_core"

type StorefrontCatalogOfferGrant struct {
	TemplateID string
	Quantity int
	ProfileType ShopGrantProfileType
}

type StorefrontCatalogOfferPriceType string
const StorefrontCatalogOfferPriceTypeMtxCurrency StorefrontCatalogOfferPriceType = "MtxCurrency"
const StorefrontCatalogOfferPriceTypeRealMoney StorefrontCatalogOfferPriceType = "RealMoney"

type StorefrontCatalogOfferPriceSaleType string
const StorefrontCatalogOfferPriceSaleTypeNone StorefrontCatalogOfferPriceSaleType = ""
const StorefrontCatalogOfferPriceSaleTypeAmountOff StorefrontCatalogOfferPriceSaleType = "AmountOff"
const StorefrontCatalogOfferPriceSaleTypeStrikethrough StorefrontCatalogOfferPriceSaleType = "Strikethrough"

type StorefrontCatalogOfferPriceMtxCurrency struct {
	PriceType StorefrontCatalogOfferPriceType
	SaleType StorefrontCatalogOfferPriceSaleType
	OriginalPrice int
	FinalPrice int
}

type StorefrontCatalogOfferPriceRealMoney struct {
	PriceType StorefrontCatalogOfferPriceType
	SaleType StorefrontCatalogOfferPriceSaleType
	BasePrice float64
	LocalPrice float64
}

type OfferDisplay struct {
	Title string
	Description string
	ShortDescription string
	LongDescription string
}

type StorefrontCatalogOfferEnum int
const StorefrontCatalogOfferEnumItem StorefrontCatalogOfferEnum = 0
const StorefrontCatalogOfferEnumCurrency StorefrontCatalogOfferEnum = 1
const StorefrontCatalogOfferEnumStarterKit StorefrontCatalogOfferEnum = 2
const StorefrontCatalogOfferEnumBattlePass StorefrontCatalogOfferEnum = 3

type StorefrontCatalogOfferGeneric interface {
	StorefrontCatalogOfferTypeItem | StorefrontCatalogOfferTypeCurrency | StorefrontCatalogOfferTypeStarterKit | StorefrontCatalogOfferTypeBattlePass
}

type StorefrontCatalogOffer[T StorefrontCatalogOfferGeneric] interface {
	GetOffer() *T
	GetOfferID() string
	GetOfferType() StorefrontCatalogOfferEnum
	GetRewards() []*StorefrontCatalogOfferGrant
	GenerateFortniteCatalogOfferResponse() aid.JSON
	GenerateFortniteBulkOffersResponse() aid.JSON
}

var storefrontCatalogOfferPriceMultiplier = map[string]float64{
	"USD": 1.2503128911,
	"GBP": 1.0,
}

type StorefrontCatalogSection struct {
	SectionType StorefrontCatalogOfferEnum
	Name string
	Offers []interface{} // *StorefrontCatalogOfferTypeItem | *StorefrontCatalogOfferTypeCurrency | *StorefrontCatalogOfferTypeStarterKit | *StorefrontCatalogOfferTypeBattlePass
}