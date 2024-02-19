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
	FortniteSets map[string]*FortniteSet `json:"sets"`
	FortniteItems map[string]*FortniteItem `json:"items"`
	FortniteItemsWithDisplayAssets map[string]*FortniteItem `json:"-"`
	FortniteItemsWithFeaturedImage []*FortniteItem `json:"-"`
	TypedFortniteItems map[string][]*FortniteItem `json:"-"`
	TypedFortniteItemsWithDisplayAssets map[string][]*FortniteItem `json:"-"`
	SnowVariantTokens map[string]*FortniteVariantToken `json:"variants"`
	StorefrontCosmeticOfferPriceLookup map[string]map[string]int `json:"-"`
	StorefrontDailyItemCountLookup []struct{Season int;Items int} `json:"-"`
	StorefrontWeeklySetCountLookup []struct{Season int;Sets int} `json:"-"`
	StorefrontCurrencyOfferPriceLookup map[string]map[int]int `json:"-"`
}

func NewDataClient() *dataClient {
	return &dataClient{
		h: &http.Client{},
		FortniteItems: make(map[string]*FortniteItem),
		FortniteSets: make(map[string]*FortniteSet),
		FortniteItemsWithDisplayAssets: make(map[string]*FortniteItem),
		FortniteItemsWithFeaturedImage: []*FortniteItem{},
		TypedFortniteItems: make(map[string][]*FortniteItem),
		TypedFortniteItemsWithDisplayAssets: make(map[string][]*FortniteItem),
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
		},
		StorefrontCurrencyOfferPriceLookup: map[string]map[int]int{
			"USD": {
				1000: 999,
				2800: 2499,
				5000: 3999,
				13500: 9999,
			},
			"GBP": {
				1000: 799,
				2800: 1999,
				5000: 3499,
				13500: 7999,
			},
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

	content := &FortniteCosmeticsResponse{}
	err = json.Unmarshal(bodyBytes, content)
	if err != nil {
		return
	}
	
	for _, item := range content.Data {
		c.LoadItem(&item)
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

	variantTokens := storage.HttpAsset[map[string]SnowCosmeticVariantToken]("variants.snow.json")
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
}

func (c *dataClient) LoadItem(item *FortniteItem) {
	if item.Introduction.BackendValue > aid.Config.Fortnite.Season || item.Introduction.BackendValue == 0 {
		return
	}
	
	if c.FortniteSets[item.Set.BackendValue] == nil {
		c.FortniteSets[item.Set.BackendValue] = &FortniteSet{
			BackendName: item.Set.Value,
			DisplayName: item.Set.Text,
			Items: []*FortniteItem{},
		}
	}

	if c.TypedFortniteItems[item.Type.BackendValue] == nil {
		c.TypedFortniteItems[item.Type.BackendValue] = []*FortniteItem{}
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

func (c *dataClient) AddBackpackToItem(backpack *FortniteItem) {
	if backpack.ItemPreviewHeroPath == "" {
		return
	}

	splitter := strings.Split(backpack.ItemPreviewHeroPath, "/")
	character, found := c.FortniteItems[splitter[len(splitter) - 1]]
	if !found {
		return
	}

	character.Backpack = backpack
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

	found.DisplayAssetPath2 = displayAsset
	c.FortniteItemsWithDisplayAssets[found.ID] = found
	c.TypedFortniteItemsWithDisplayAssets[found.Type.BackendValue] = append(c.TypedFortniteItemsWithDisplayAssets[found.Type.BackendValue], found)
}

func (c *dataClient) AddNumericStylesToItem(item *FortniteItem) {
	ownedStyles := []FortniteVariantChannel{}
	for i := 0; i < 100; i++ {
		ownedStyles = append(ownedStyles, FortniteVariantChannel{
			Tag: fmt.Sprint(i),
		})
	}

	item.Variants = append(item.Variants, FortniteVariant{
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
	for _, item := range c.StorefrontWeeklySetCountLookup {
		if item.Season > season {
			continue
		}
		currentValue = item.Sets
	}
	return currentValue
}

func (c *dataClient) GetStorefrontCosmeticOfferPrice(rarity string, type_ string) int {
	return c.StorefrontCosmeticOfferPriceLookup[rarity][type_]
}

func (c *dataClient) GetStorefrontCurrencyOfferPrice(currency string, amount int) int {
	return c.StorefrontCurrencyOfferPriceLookup[currency][amount]
}

func PreloadCosmetics() error {
	DataClient = NewDataClient()
	DataClient.LoadExternalData()

	aid.Print("(snow) " + fmt.Sprint(len(DataClient.FortniteItems)) + " cosmetics loaded from fortnite-api.com")
	return nil
}

func GetItemByShallowID(shallowID string) *FortniteItem {
	for _, item := range DataClient.TypedFortniteItems["AthenaCharacter"] {
		if strings.Contains(item.ID, shallowID) {
			return item
		}
	}

	return nil
}

func GetRandomItemWithDisplayAsset() *FortniteItem {
	items := DataClient.FortniteItemsWithDisplayAssets
	if len(items) == 0 {
		return nil
	}

	flat := []FortniteItem{}
	for _, item := range items {
		flat = append(flat, *item)
	}

	slices.SortFunc[[]FortniteItem](flat, func(a, b FortniteItem) int {
		return strings.Compare(a.ID, b.ID)
	})

	return &flat[aid.RandomInt(0, len(flat))]
}

func GetRandomItemWithDisplayAssetOfNotType(notType string) *FortniteItem {
	flat := []FortniteItem{}
	
	for t, items := range DataClient.TypedFortniteItemsWithDisplayAssets {
		if t == notType {
			continue
		}

		for _, item := range items {
			flat = append(flat, *item)
		}
	}

	slices.SortFunc[[]FortniteItem](flat, func(a, b FortniteItem) int {
		return strings.Compare(a.ID, b.ID)
	})

	return &flat[aid.RandomInt(0, len(flat))]
}

func GetRandomSet() *FortniteSet {
	sets := []FortniteSet{}
	for _, set := range DataClient.FortniteSets {
		if set.BackendName == "" {
			continue
		}
		sets = append(sets, *set)
	}

	slices.SortFunc[[]FortniteSet](sets, func(a, b FortniteSet) int {
		return strings.Compare(a.BackendName, b.BackendName)
	})

	return &sets[aid.RandomInt(0, len(sets))]
}