package fortnite

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/storage"
)

var (
	DataClient *dataClient
)

type dataClient struct {
	h *http.Client
	FortniteSets map[string]*APISetDefinition `json:"sets"`
	FortniteItems map[string]*APICosmeticDefinition `json:"items"`
	FortniteItemsWithDisplayAssets map[string]*APICosmeticDefinition `json:"-"`
	FortniteItemsWithFeaturedImage []*APICosmeticDefinition `json:"-"`
	TypedFortniteItems map[string][]*APICosmeticDefinition `json:"-"`
	TypedFortniteItemsWithDisplayAssets map[string][]*APICosmeticDefinition `json:"-"`
	SnowVariantTokens map[string]*FortniteVariantToken `json:"variants"`
	StorefrontCosmeticOfferPriceLookup map[string]map[string]int `json:"-"`
	StorefrontDailyItemCountLookup []struct{Season int;Items int} `json:"-"`
	StorefrontWeeklySetCountLookup []struct{Season int;Sets int} `json:"-"`
	StorefrontCurrencyOfferPriceLookup map[string]map[int]int `json:"-"`
	StorefrontCurrencyMultiplier map[string]float64 `json:"-"`
	SnowSeason *SnowSeasonDefinition `json:"season"`
}

func NewDataClient() *dataClient {
	return &dataClient{
		h: &http.Client{},
		FortniteItems: make(map[string]*APICosmeticDefinition),
		FortniteSets: make(map[string]*APISetDefinition),
		FortniteItemsWithDisplayAssets: make(map[string]*APICosmeticDefinition),
		FortniteItemsWithFeaturedImage: []*APICosmeticDefinition{},
		TypedFortniteItems: make(map[string][]*APICosmeticDefinition),
		TypedFortniteItemsWithDisplayAssets: make(map[string][]*APICosmeticDefinition),
		SnowVariantTokens: make(map[string]*FortniteVariantToken),
		StorefrontDailyItemCountLookup: []struct{Season int;Items int}{
			{2, 4},
			{4, 6},
			{13, 10},
		},
		StorefrontWeeklySetCountLookup: []struct{Season int;Sets int}{
			{2, 2},
			{4, 3},
			{13, 5},
		},
		StorefrontCosmeticOfferPriceLookup: map[string]map[string]int{
			"EFortRarity::Legendary": {
				"AthenaCharacter": 2000,
				"AthenaBackpack":  300,
				"AthenaPickaxe":   1500,
				"AthenaGlider":    1800,
				"AthenaDance":     500,
				"AthenaItemWrap":  800,
			},
			"EFortRarity::Epic": {
				"AthenaCharacter": 1500,
				"AthenaBackpack":  250,
				"AthenaPickaxe":   1200,
				"AthenaGlider":    1500,
				"AthenaDance":     800,
				"AthenaItemWrap":  800,
			},
			"EFortRarity::Rare": {
				"AthenaCharacter": 1200,
				"AthenaBackpack":  200,
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
		},
		StorefrontCurrencyOfferPriceLookup: map[string]map[int]int{
			"USD": {
				1000: 999,
				2800: 2499,
				5000: 3999,
				7500: 5999,
				13500: 9999,
			},
			"GBP": {
				1000: 799,
				2800: 1999,
				5000: 3499,
				7500: 4999,
				13500: 7999,
			},
		},
		StorefrontCurrencyMultiplier: map[string]float64{
			"USD": 1.2503128911,
			"GBP": 1.0,
		},
	}
}

func (c *dataClient) LoadExternalData() {
	req, err := http.NewRequest("GET", "https://fortnite-api.com/v2/cosmetics/br", nil)
	if err != nil {
		return
	}

	resp, err := c.h.Do(req)
	if err != nil {
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	content := &APICosmeticsResponse{}
	err = json.Unmarshal(bodyBytes, content)
	if err != nil {
		return
	}
	
	for _, item := range content.Data {
		c.LoadItemDefinition(&item)
	}

	for _, item := range c.TypedFortniteItems["AthenaBackpack"] {
		c.AddBackpackToItem(item)
	}

	displayAssets := storage.HttpAsset[[]string]("assets.snow.json")
	if displayAssets == nil {
		return
	}

	for _, displayAsset := range *displayAssets {
		c.AddDisplayAssetToItem(displayAsset)
	}

	variantTokens := storage.HttpAsset[map[string]SnowCosmeticVariantDefinition]("variants.snow.json")
	if variantTokens == nil {
		return
	}
	
	for k, v := range *variantTokens {
		item := c.FortniteItems[v.Item]
		if item == nil {
			continue
		}

		c.SnowVariantTokens[k] = &FortniteVariantToken{
			Grants: v.Grants,
			Item: item,
			Name: v.Name,
			Gift: v.Gift,
			Equip: v.Equip,
			Unseen: v.Unseen,
		}
	}

	addNumericStylesToSets := []string{"Soccer", "Football", "ScaryBall"} 
	for _, setValue := range addNumericStylesToSets {
		set, found := c.FortniteSets[setValue]
		if !found {
			continue
		}

		for _, item := range set.Items {
			c.AddNumericStylesToItem(item)
		}
	}

	athenaSeasonObj := storage.HttpAsset[[]UnrealEngineObject[UnrealSeasonProperties, UnrealNoRows]]("season.snow.json")
	if athenaSeasonObj == nil {
		return
	}

	levelProgressionObj := storage.HttpAsset[[]UnrealEngineObject[UnrealProgressionProperties, UnrealProgressionRows]]("progression.levels.snow.json")
	if levelProgressionObj == nil {
		return
	}

	bookProgressionObj := storage.HttpAsset[[]UnrealEngineObject[UnrealProgressionProperties, UnrealProgressionRows]]("progression.book.snow.json")
	if bookProgressionObj == nil {
		return
	}

	c.SnowSeason = NewSeasonDefinition()

	for index, tier := range (*athenaSeasonObj)[0].Properties.BookXpSchedulePaid.Levels {
		c.SnowSeason.TierRewardsPremium[index] = []*ItemGrant{}
		for _, reward := range tier.Rewards {
			templateId := aid.Ternary[string](reward.ItemDefinition.AssetPathName != "None", convertAssetPathToTemplateId(reward.ItemDefinition.AssetPathName), reward.TemplateID)
			c.SnowSeason.TierRewardsPremium[index] = append(c.SnowSeason.TierRewardsPremium[index], NewItemGrant(templateId, reward.Quantity))
		}
	}
	c.SnowSeason.TierRewardsPremium[0] = []*ItemGrant{}

	for index, tier := range (*athenaSeasonObj)[0].Properties.BookXpScheduleFree.Levels {
		c.SnowSeason.TierRewardsFree[index] = []*ItemGrant{}
		for _, reward := range tier.Rewards {
			templateId := aid.Ternary[string](reward.ItemDefinition.AssetPathName != "None", convertAssetPathToTemplateId(reward.ItemDefinition.AssetPathName), reward.TemplateID)
			c.SnowSeason.TierRewardsFree[index] = append(c.SnowSeason.TierRewardsFree[index], NewItemGrant(templateId, reward.Quantity))
		}
	}
	c.SnowSeason.TierRewardsFree[0] = []*ItemGrant{}

	for index, level := range (*athenaSeasonObj)[0].Properties.SeasonXpScheduleFree.Levels {
		c.SnowSeason.LevelRewards[index] = []*ItemGrant{}
		for _, reward := range level.Rewards {
			templateId := aid.Ternary[string](reward.ItemDefinition.AssetPathName != "None", convertAssetPathToTemplateId(reward.ItemDefinition.AssetPathName), reward.TemplateID)
			c.SnowSeason.LevelRewards[index] = append(c.SnowSeason.LevelRewards[index], NewItemGrant(templateId, reward.Quantity))
		}
	}
	c.SnowSeason.LevelRewards[0] = []*ItemGrant{}

	for _, token := range (*athenaSeasonObj)[0].Properties.TokensToRemoveAtSeasonEnd {
		c.SnowSeason.SeasonTokenRemoval = append(c.SnowSeason.SeasonTokenRemoval, NewItemGrant(convertAssetPathToTemplateId(token.AssetPathName), 1))
	}

	for _, reward := range (*athenaSeasonObj)[0].Properties.SeasonFirstWinRewards.Rewards {
		templateId := aid.Ternary[string](reward.ItemDefinition.AssetPathName != "None", convertAssetPathToTemplateId(reward.ItemDefinition.AssetPathName), reward.TemplateID)
		c.SnowSeason.VictoryRewards = append(c.SnowSeason.VictoryRewards, NewItemGrant(templateId, reward.Quantity))
	}

	for _, token := range (*athenaSeasonObj)[0].Properties.ExpiringRewardTypes {
		c.SnowSeason.SeasonTokenRemoval = append(c.SnowSeason.SeasonTokenRemoval, NewItemGrant(convertAssetPathToTemplateId(token.AssetPathName), 1))
	}

	for _, reward := range (*athenaSeasonObj)[0].Properties.SeasonGrantsToEveryone.Rewards {
		templateId := aid.Ternary[string](reward.ItemDefinition.AssetPathName != "None", convertAssetPathToTemplateId(reward.ItemDefinition.AssetPathName), reward.TemplateID)
		c.SnowSeason.TierRewardsFree[0] = append(c.SnowSeason.TierRewardsFree[0], NewItemGrant(templateId, reward.Quantity))
	}

	for _, replacement := range (*athenaSeasonObj)[0].Properties.BattleStarSubstitutionReward.Rewards {
		templateId := aid.Ternary[string](replacement.ItemDefinition.AssetPathName != "None", convertAssetPathToTemplateId(replacement.ItemDefinition.AssetPathName), replacement.TemplateID)
		c.SnowSeason.BookXPReplacements = append(c.SnowSeason.BookXPReplacements, NewItemGrant(templateId, replacement.Quantity))
	}

	for indexString, row := range *(*levelProgressionObj)[0].Rows {
		index := aid.ToInt(indexString)
		c.SnowSeason.LevelProgression[index] = &row
	}
	c.SnowSeason.LevelProgression[0] = &UnrealProgressionRows{
		Level: 0,
		XpToNextLevel: 0,
		XpTotal: 0,
	}
	c.SnowSeason.LevelProgression[100] = &UnrealProgressionRows{
		Level: 100,
		XpToNextLevel: 0,
		XpTotal: c.SnowSeason.LevelProgression[99].XpTotal + c.SnowSeason.LevelProgression[99].XpToNextLevel,
	}

	for indexString, row := range *(*bookProgressionObj)[0].Rows {
		index := aid.ToInt(indexString)
		c.SnowSeason.BookProgression[index] = &row
	}
	c.SnowSeason.BookProgression[0] = &UnrealProgressionRows{
		Level: 0,
		XpToNextLevel: 0,
		XpTotal: 0,
	}
	c.SnowSeason.BookProgression[100] = &UnrealProgressionRows{
		Level: 100,
		XpToNextLevel: 0,
		XpTotal: c.SnowSeason.BookProgression[99].XpTotal + c.SnowSeason.BookProgression[99].XpToNextLevel,
	}

	c.SnowSeason.DefaultOfferID = (*athenaSeasonObj)[0].Properties.BattlePassOfferId
	c.SnowSeason.BundleOfferID = (*athenaSeasonObj)[0].Properties.BattlePassBundleOfferId
	c.SnowSeason.TierOfferID = (*athenaSeasonObj)[0].Properties.BattlePassLevelOfferID
}

func (c *dataClient) LoadItemDefinition(item *APICosmeticDefinition) {
	if item.Introduction.BackendValue > aid.Config.Fortnite.Season || item.Introduction.BackendValue == 0 {
		return
	}

	typeLookup := map[string]string{
		"AthenaCharacter": "AthenaCharacter",
		"AthenaBackpack": "AthenaBackpack",
		"AthenaPickaxe": "AthenaPickaxe",
		"AthenaGlider": "AthenaGlider",
		"AthenaDance": "AthenaDance",
		"AthenaToy": "AthenaDance",
		"AthenaEmoji": "AthenaEmoji",
		"AthenaItemWrap": "AthenaItemWrap",
		"AthenaMusicPack": "AthenaMusicPack",
		"AthenaPet": "AthenaBackpack",
		"AthenaPetCarrier": "AthenaBackpack",
		"AthenaLoadingScreen": "AthenaLoadingScreen",
		"AthenaSkyDiveContrail": "AthenaSkyDiveContrail",
	}

	item.Type.BackendValue = aid.Ternary[string](typeLookup[item.Type.BackendValue] != "", typeLookup[item.Type.BackendValue], item.Type.BackendValue)
	
	if c.FortniteSets[item.Set.BackendValue] == nil {
		c.FortniteSets[item.Set.BackendValue] = &APISetDefinition{
			BackendName: item.Set.Value,
			DisplayName: item.Set.Text,
			Items: []*APICosmeticDefinition{},
		}
	}

	if c.TypedFortniteItems[item.Type.BackendValue] == nil {
		c.TypedFortniteItems[item.Type.BackendValue] = []*APICosmeticDefinition{}
	}

	c.FortniteItems[item.ID] = item
	c.FortniteSets[item.Set.BackendValue].Items = append(c.FortniteSets[item.Set.BackendValue].Items, item)
	c.TypedFortniteItems[item.Type.BackendValue] = append(c.TypedFortniteItems[item.Type.BackendValue], item)

	if item.Type.BackendValue != "AthenaCharacter" || item.Images.Featured == "" || slices.Contains[[]string]([]string{
		"Soccer",
		"Football",
		"Waypoint",
	}, item.Set.BackendValue) {
		return
	}

	for _, tag := range item.GameplayTags {
		if strings.Contains(tag, "StarterPack") {
			return
		}
	}

	c.FortniteItemsWithFeaturedImage = append(c.FortniteItemsWithFeaturedImage, item)
}

func (c *dataClient) AddBackpackToItem(backpack *APICosmeticDefinition) {
	if backpack.ItemPreviewHeroPath == "" {
		return
	}

	splitter := strings.Split(backpack.ItemPreviewHeroPath, "/")
	character, found := c.FortniteItems[splitter[len(splitter) - 1]]
	if !found {
		return
	}

	character.BackpackDefinition = backpack
}

func (c *dataClient) AddDisplayAssetToItem(displayAsset string) {
	split := strings.Split(displayAsset, "_")[1:]
	found := c.FortniteItems[strings.Join(split[:], "_")]

	if found == nil && split[0] == "CID" {
		r := aid.Regex(strings.Join(split[:], "_"), `(?:CID_)(\d+|A_\d+)(?:_.+)`)
		if r != nil {
			found = GetItemByShallowID(*r)
		}
	}

	if found == nil {
		return
	}

	found.NewDisplayAssetPath = displayAsset
	c.FortniteItemsWithDisplayAssets[found.ID] = found
	c.TypedFortniteItemsWithDisplayAssets[found.Type.BackendValue] = append(c.TypedFortniteItemsWithDisplayAssets[found.Type.BackendValue], found)
}

func (c *dataClient) AddNumericStylesToItem(item *APICosmeticDefinition) {
	ownedStyles := []APICosmeticDefinitionVariant{}
	for i := 0; i < 100; i++ {
		ownedStyles = append(ownedStyles, APICosmeticDefinitionVariant{
			Tag: fmt.Sprint(i),
		})
	}

	item.Variants = append(item.Variants, APICosmeticDefinitionVariantChannel{
		Channel: "Numeric",
		Type: "int",
		Options: ownedStyles,
	})
}

func (c *dataClient) GetStorefrontDailyItemCount(season int) int {
	currentValue := 4
	for _, item := range c.StorefrontDailyItemCountLookup {
		if item.Season > season {
			continue
		}
		currentValue = item.Items
	}
	return currentValue
}

func (c *dataClient) GetStorefrontWeeklySetCount(season int) int {
	currentValue := 2
	// for _, item := range c.StorefrontWeeklySetCountLookup {
	// 	if item.Season > season {
	// 		continue
	// 	}
	// 	currentValue = item.Sets
	// }
	return currentValue
}

func (c *dataClient) GetStorefrontCosmeticOfferPrice(rarity string, type_ string) int {
	return c.StorefrontCosmeticOfferPriceLookup[rarity][type_]
}

func (c *dataClient) GetStorefrontCurrencyOfferPrice(currency string, amount int) int {
	return c.StorefrontCurrencyOfferPriceLookup[currency][amount]
}

func (c *dataClient) GetStorefrontLocalizedOfferPrice(currency string, amount int) int {
	return int(float64(amount) * c.StorefrontCurrencyMultiplier[currency])
}

func PreloadCosmetics() error {
	DataClient = NewDataClient()
	DataClient.LoadExternalData()

	aid.Print("(snow) " + fmt.Sprint(len(DataClient.FortniteItems)) + " cosmetics loaded from fortnite-api.com")
	return nil
}

func GetItemByShallowID(shallowID string) *APICosmeticDefinition {
	for _, item := range DataClient.TypedFortniteItems["AthenaCharacter"] {
		if strings.Contains(item.ID, shallowID) {
			return item
		}
	}

	return nil
}

func GetRandomItemWithDisplayAsset() *APICosmeticDefinition {
	items := DataClient.FortniteItemsWithDisplayAssets
	if len(items) == 0 {
		return nil
	}

	flat := []APICosmeticDefinition{}
	for _, item := range items {
		flat = append(flat, *item)
	}

	slices.SortFunc[[]APICosmeticDefinition](flat, func(a, b APICosmeticDefinition) int {
		return strings.Compare(a.ID, b.ID)
	})

	return &flat[aid.RandomInt(0, len(flat))]
}

func GetRandomItemWithDisplayAssetOfNotType(notType string) *APICosmeticDefinition {
	flat := []APICosmeticDefinition{}
	
	for t, items := range DataClient.TypedFortniteItemsWithDisplayAssets {
		if t == notType {
			continue
		}

		for _, item := range items {
			flat = append(flat, *item)
		}
	}

	slices.SortFunc[[]APICosmeticDefinition](flat, func(a, b APICosmeticDefinition) int {
		return strings.Compare(a.ID, b.ID)
	})

	return &flat[aid.RandomInt(0, len(flat))]
}

func GetRandomSet() *APISetDefinition {
	sets := []APISetDefinition{}
	for _, set := range DataClient.FortniteSets {
		if set.BackendName == "" {
			continue
		}
		sets = append(sets, *set)
	}

	slices.SortFunc[[]APISetDefinition](sets, func(a, b APISetDefinition) int {
		return strings.Compare(a.BackendName, b.BackendName)
	})

	return &sets[aid.RandomInt(0, len(sets))]
}