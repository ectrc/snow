package shop

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/fortnite"
)

func GetShop() *StorefrontCatalog {
	aid.SetRandom(rand.New(rand.NewSource(int64(aid.Config.Fortnite.ShopSeed) + aid.CurrentDayUnix())))
	shop := NewStorefrontCatalog()

	dailySection := NewStorefrontCatalogSection("BRDailyStorefront", StorefrontCatalogOfferEnumItem)
	weeklySection := NewStorefrontCatalogSection("BRWeeklyStorefront", StorefrontCatalogOfferEnumItem)
	moneySection := NewStorefrontCatalogSection("CurrencyStorefront", StorefrontCatalogOfferEnumCurrency)
	kitSection := NewStorefrontCatalogSection("BRStarterKits", StorefrontCatalogOfferEnumStarterKit)
	bookSection := NewStorefrontCatalogSection(fmt.Sprintf("BRSeason%d", aid.Config.Fortnite.Season), StorefrontCatalogOfferEnumBattlePass)
	shop.AddSections(bookSection, dailySection, weeklySection, moneySection, kitSection)

	bookDefaultOffer := newBookOffer(aid.Ternary[string](fortnite.DataClient.SnowSeason.DefaultOfferID != "", fortnite.DataClient.SnowSeason.DefaultOfferID, "book://"+ aid.Hash([]byte(aid.RandomString(32)))), 950, 0, &StorefrontCatalogOfferGrant{TemplateID: "Snow:BattlePass", Quantity: 1, ProfileType: "athena"})
	bookDefaultOffer.Diplay.Title = "Battle Pass"
	bookDefaultOffer.Diplay.ShortDescription = "Claim your Season 8 Battle Pass!"
	bookDefaultOffer.Diplay.Description = "Fortnite Season 8\n\nInstantly get these items <Bold>valued at over 3,500 V-Bucks</>.\n  • <ItemName>Blackheart</> Progressive Outfit\n  • <ItemName>Hybrid</> Progressive Outfit\n  • <Bold>50% Bonus</> Season Match XP\n  • <Bold>10% Bonus</> Season Friend Match XP\n  • <Bold>Extra Weekly Challenges</>\n\nPlay to level up your Battle Pass, unlocking <Bold>over 100 rewards</> (typically takes 75 to 150 hours of play).\n  • <ItemName>Sidewinder</> and <Bold>4 more Outfits</>\n  • <Bold>1,300 V-Bucks</>\n  • 7 Emotes\n  • 6 Wraps\n  • 2 Pets\n  • 5 Harvesting Tools\n  • 4 Gliders\n  • 4 Back Blings\n  • 5 Contrails\n  • 14 Sprays\n  • 3 Music Tracks\n  • 1 Toy\n  • 20 Loading Screens\n  • and so much more!\nWant it all faster? You can use V-Bucks to buy tiers any time!"
	bookDefaultOffer.Meta.DisplayAssetPath = fmt.Sprintf("/Game/Catalog/DisplayAssets/DA_BR_Season%d_BattlePass.DA_BR_Season%d_BattlePass", aid.Config.Fortnite.Season, aid.Config.Fortnite.Season)
	bookDefaultOffer.Meta.Priority = 1
	bookSection.AddOffer(bookDefaultOffer)

	bookBundleOffer := newBookOffer(aid.Ternary[string](fortnite.DataClient.SnowSeason.BundleOfferID != "", fortnite.DataClient.SnowSeason.BundleOfferID, "book://"+ aid.Hash([]byte(aid.RandomString(32)))), 4700, 1850, []*StorefrontCatalogOfferGrant{
		{TemplateID: "Snow:BattlePass", Quantity: 1, ProfileType: "athena"},
		{TemplateID: "AccountResource:AthenaBattleStar", Quantity: 250, ProfileType: "athena"},
	}...)
	bookBundleOffer.Diplay.Title = "Battle Bundle"
	bookBundleOffer.Diplay.ShortDescription = "Claim your Season 8 Battle Pass + 25 Tiers!"
	bookBundleOffer.Diplay.Description = "Fortnite Season 8\n\nInstantly get these items <Bold>valued at over 10,000 V-Bucks</>.\n  • <ItemName>Blackheart</> Progressive Outfit\n  • <ItemName>Hybrid</> Progressive Outfit\n  • <ItemName>Sidewinder</> Outfit\n  • <ItemName>Tropical Camo</> Wrap\n  • <ItemName>Woodsy</> Pet\n  • <ItemName>Sky Serpents</> Glider\n  • <ItemName>Cobra</> Back Bling\n  • <ItemName>Flying Standard</> Contrail\n  • 300 V-Bucks\n  • 1 Music Track\n  • <Bold>70% Bonus</> Season Match XP\n  • <Bold>20% Bonus</> Season Friend Match XP\n  • <Bold>Extra Weekly Challenges</>\n  • and more!\n\nPlay to level up your Battle Pass, unlocking <Bold>over 75 rewards</> (typically takes 75 to 150 hours of play).\n  • <Bold>4 more Outfits</>\n  • <Bold>1,000 V-Bucks</>\n  • 6 Emotes\n  • 5 Wraps\n  • 3 Gliders\n  • 3 Back Blings\n  • 4 Harvesting Tools\n  • 4 Contrails\n  • 1 Pet\n  • 12 Sprays\n  • 2 Music Tracks\n  • and so much more!\nWant it all faster? You can use V-Bucks to buy tiers any time!"
	bookBundleOffer.Meta.DisplayAssetPath = fmt.Sprintf("/Game/Catalog/DisplayAssets/DA_BR_Season%d_BattlePassWithLevels.DA_BR_Season%d_BattlePassWithLevels", aid.Config.Fortnite.Season, aid.Config.Fortnite.Season)
	bookBundleOffer.Meta.Priority = 0
	bookSection.AddOffer(bookBundleOffer)

	bookLevelOffer := newBookOffer(aid.Ternary[string](fortnite.DataClient.SnowSeason.TierOfferID != "", fortnite.DataClient.SnowSeason.TierOfferID, "book://"+ aid.Hash([]byte(aid.RandomString(32)))), 150, 150, &StorefrontCatalogOfferGrant{TemplateID: "AccountResource:AthenaBattleStar", Quantity: 10, ProfileType: "athena"})
	bookSection.AddOffer(bookLevelOffer)

	for len(dailySection.Offers) <= fortnite.DataClient.GetStorefrontDailyItemCount(aid.Config.Fortnite.Season) {
		offer := newItemOffer(fortnite.GetRandomItemWithDisplayAssetOfNotType("AthenaCharacter"), true, true)
		offer.Meta.SectionID = "Daily"
		dailySection.AddOffer(offer)
	}
	
	for weeklySection.GetGroupedOffersLength() < fortnite.DataClient.GetStorefrontWeeklySetCount(aid.Config.Fortnite.Season) {
		set := fortnite.GetRandomSet()
		for _, item := range set.Items {
			offer := newItemOffer(item, true, true)
			offer.Meta.SectionID = "Featured"
			offer.Categories = append(offer.Categories, set.BackendName)
			weeklySection.AddOffer(offer)
		}
	}

	xp := NewItemCatalogOffer()
	xp.OfferID = "item://AthenaSeasonalXP"
	xp.Meta.TileSize = "Small"
	xp.Meta.Giftable = false
	xp.Meta.Refundable = false
	xp.Meta.DisplayAssetPath = "/Game/Catalog/DisplayAssets/DA_FoundersPack_4_5.DA_FoundersPack_4_5"
	xp.Rewards = append(xp.Rewards, &StorefrontCatalogOfferGrant{
		TemplateID: "AccountResource:AthenaSeasonalXP",
		Quantity: 100000,
		ProfileType: "athena",
	})
	xp.Price.PriceType = StorefrontCatalogOfferPriceTypeMtxCurrency
	xp.Price.SaleType = StorefrontCatalogOfferPriceSaleTypeNone
	xp.Price.OriginalPrice = 1000
	xp.Price.FinalPrice = 1000
	weeklySection.AddOffer(xp)
	
	moneySection.AddOffer(newMoneyOffer(1000, 0, "https://cdn1.epicgames.com/offer/fn/EGS_VBucks_1000_1200x1600-c8a13f66ba88744d5216f884855e2a4d", 3))
	moneySection.AddOffer(newMoneyOffer(2800, 300, "https://cdn1.epicgames.com/offer/fn/EGS_VBucks_2800_1200x1600-055112a56c0fb d65989470ece7c653f", 2))
	moneySection.AddOffer(newMoneyOffer(7500, 1500, "https://cdn1.epicgames.com/offer/fn/EGS_VBucks_5000_1200x1600-8ea53bb4ea3d75821153075df8e3ca95", 1))
	moneySection.AddOffer(newMoneyOffer(13500, 3500, "https://cdn1.epicgames.com/offer/fn/EGS_VBucks_13500_1200x1600-39489a289769bc6c1d14f4a8b53b48f4", 0))

	lagunaKit := newKitOffer("The Laguna Pack", 499, 8, []*StorefrontCatalogOfferGrant{
		{TemplateID: "AthenaCharacter:CID_367_Athena_Commando_F_Tropical", Quantity: 1, ProfileType: "athena"},
		{TemplateID: "AthenaBackpack:BID_231_TropicalFemale", Quantity: 1, ProfileType: "athena"},
		{TemplateID: "AthenaItemWrap:Wrap_033_TropicalGirl", Quantity: 1, ProfileType: "athena"},
		{TemplateID: "Currency:MtxPurchased", Quantity: 600, ProfileType: "common_core"},
	}...)
	lagunaKit.Meta.DisplayAssetPath = "/Game/Catalog/DisplayAssets/DA_Featured_CID_367_Athena_Commando_F_Tropical.DA_Featured_CID_367_Athena_Commando_F_Tropical"
	lagunaKit.Meta.FeaturedImageURL = "https://fortnite-api.com/images/cosmetics/br/CID_367_Athena_Commando_F_Tropical/icon.png"
	kitSection.AddOffer(lagunaKit)

	ikonikKit := newKitOffer("Ikonik Pack", 3999, 8, []*StorefrontCatalogOfferGrant{
		{TemplateID: "AthenaCharacter:CID_313_Athena_Commando_M_KpopFashion", Quantity: 1, ProfileType: "athena"},
		{TemplateID: "AthenaDance:EID_KPopDance03", Quantity: 1, ProfileType: "athena"},
		{TemplateID: "Currency:MtxPurchased", Quantity: 600, ProfileType: "common_core"},
	}...)
	ikonikKit.Meta.DisplayAssetPath = "/Game/Catalog/DisplayAssets/DA_Featured_CID_313_Athena_Commando_M_KpopFashion.DA_Featured_CID_313_Athena_Commando_M_KpopFashion"
	ikonikKit.Meta.FeaturedImageURL = "https://fortnite-api.com/images/cosmetics/br/CID_313_Athena_Commando_M_KpopFashion/icon.png"
	kitSection.AddOffer(ikonikKit)

	return shop
}

func newItemOffer(item *fortnite.APICosmeticDefinition, addAssets, giftable bool) *StorefrontCatalogOfferTypeItem {
	displayAsset := regexp.MustCompile(`[^/]+$`).FindString(item.DisplayAssetPath)
	
	offer := NewItemCatalogOffer()
	offer.Meta.TileSize = aid.Ternary[string](item.Type.BackendValue == "AthenaCharacter", "Small", "Normal")
	offer.Meta.Giftable = giftable
	offer.Meta.Refundable = true
	if addAssets {
		offer.Meta.DisplayAssetPath = aid.Ternary[string](displayAsset != "", "/Game/Catalog/DisplayAssets/" + displayAsset + "." + displayAsset, "")
		offer.Meta.NewDisplayAssetPath = aid.Ternary[string](item.NewDisplayAssetPath != "", "/Game/Catalog/NewDisplayAssets/" + item.NewDisplayAssetPath + "." + item.NewDisplayAssetPath, "")
	}

	offer.Rewards = append(offer.Rewards, &StorefrontCatalogOfferGrant{
		TemplateID: item.Type.BackendValue + ":" + item.ID,
		Quantity: 1,
		ProfileType: "athena",
	})

	offer.Price.PriceType = StorefrontCatalogOfferPriceTypeMtxCurrency
	offer.Price.SaleType = StorefrontCatalogOfferPriceSaleTypeNone
	offer.Price.OriginalPrice = fortnite.DataClient.GetStorefrontCosmeticOfferPrice(item.Rarity.BackendValue, item.Type.BackendValue)
	offer.Price.FinalPrice = offer.Price.OriginalPrice

	offer.OfferID = fmt.Sprintf("item://%s", aid.Hash([]byte(offer.OfferID)))

	return offer
}

func newMoneyOffer(real, bonus int, imgUrl string, position int) *StorefrontCatalogOfferTypeCurrency {
	format := aid.FormatNumber(real)
	offer := NewCurrencyCatalogOffer()

	offer.Meta.IconSize = "Small"
	offer.Meta.CurrencyAnalyticsName = fmt.Sprintf("MtxPack%d", real)
	offer.Meta.OriginalOffer = real
	offer.Meta.ExtraBonus = bonus
	offer.Meta.DisplayAssetPath = fmt.Sprintf("/Game/Catalog/DisplayAssets/DA_MtxPack%d.DA_MtxPack%d", real, real)
	offer.Meta.FeaturedImageURL = imgUrl
	offer.Meta.Priority = position

	offer.Diplay.Title = fmt.Sprintf("%s V-Bucks", format)
	offer.Diplay.Description = fmt.Sprintf("Buy %s Fortnite V-Bucks, the in-game currency that can be spent in Fortnite Battle Royale and Creative modes. You can purchase new customization items like Outfits, Gliders, Pickaxes, Emotes, Wraps and the latest season's Battle Pass! Gliders and Contrails may not be used in Save the World mode.", format)
	offer.Diplay.LongDescription = fmt.Sprintf("Buy %s Fortnite V-Bucks, the in-game currency that can be spent in Fortnite Battle Royale and Creative modes. You can purchase new customization items like Outfits, Gliders, Pickaxes, Emotes, Wraps and the latest season's Battle Pass! Gliders and Contrails may not be used in Save the World mode.\n\nAll V-Bucks purchased on the Epic Games Store are not redeemable or usable on Nintendo Switch™.", format)
	
	offer.Price.PriceType = StorefrontCatalogOfferPriceTypeRealMoney
	offer.Price.BasePrice = float64(fortnite.DataClient.GetStorefrontCurrencyOfferPrice("GBP", real))
	offer.Price.LocalPrice = float64(fortnite.DataClient.GetStorefrontCurrencyOfferPrice("USD", real))

	offer.Rewards = append(offer.Rewards, &StorefrontCatalogOfferGrant{
		TemplateID: "Currency:MtxPurchased",
		Quantity: real,
		ProfileType: "common_core",
	})

	return offer
}

func newKitOffer(title string, basePrice, season int, rewards ...*StorefrontCatalogOfferGrant) *StorefrontCatalogOfferTypeStarterKit {
	description := fmt.Sprintf("Jump into Fortnite Battle Royale with the %s. Includes:\n\n- 600 V-Bucks", strings.ReplaceAll(title, "The ", ""))
	for _, reward := range rewards {
		item := fortnite.DataClient.FortniteItems[strings.Split(reward.TemplateID, ":")[1]]
		if item != nil {
			description += fmt.Sprintf("\n- %s %s", item.Name, item.Type.DisplayValue)
		}
	}

	offer := NewStarterKitCatalogOffer()

	offer.Meta.ReleaseSeason = season
	offer.Meta.OriginalOffer = 600
	offer.Meta.ExtraBonus = 100

	offer.Diplay.Title = title
	offer.Diplay.Description = description
	offer.Diplay.LongDescription = fmt.Sprintf("%s\n\nV-Bucks are an in-game currency that can be spent in both the Battle Royale PvP mode and the Save the World PvE campaign. In Battle Royale, you can use V-bucks to purchase new customization items like outfits, emotes, pickaxes, gliders, and more! In Save the World you can purchase Llama Pinata card packs that contain weapon, trap and gadget schematics as well as new Heroes and more! \n\nNote: Items do not transfer between the Battle Royale mode and the Save the World campaign.", description)

	offer.Price.PriceType = StorefrontCatalogOfferPriceTypeRealMoney
	offer.Price.BasePrice = float64(fortnite.DataClient.GetStorefrontLocalizedOfferPrice("GBP", basePrice))
	offer.Price.LocalPrice = float64(fortnite.DataClient.GetStorefrontLocalizedOfferPrice("USD", basePrice))

	offer.Rewards = rewards

	return offer
}

func newBookOffer(customId string, ogPrice, finalprice int, rewards ...*StorefrontCatalogOfferGrant) *StorefrontCatalogOfferTypeBattlePass {
	offer := NewBattlePassCatalogOffer()
	offer.OfferID	= customId

	offer.Meta.TileSize = "Normal"
	offer.Meta.SectionID = "BattlePass"

	offer.Price.PriceType = StorefrontCatalogOfferPriceTypeMtxCurrency
	offer.Price.SaleType = StorefrontCatalogOfferPriceSaleTypeStrikethrough
	offer.Price.OriginalPrice = ogPrice
	offer.Price.FinalPrice = finalprice

	offer.Rewards = rewards

	return offer
}