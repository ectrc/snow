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
	External *ExternalDataClient
)

type ExternalDataClient struct {
	h *http.Client
	FortniteSets map[string]*FortniteSet `json:"sets"`
	FortniteItems map[string]*FortniteItem `json:"items"`
	FortniteItemsWithDisplayAssets map[string]*FortniteItem `json:"-"`
	FortniteItemsWithFeaturedImage []*FortniteItem `json:"-"`
	TypedFortniteItems map[string][]*FortniteItem `json:"-"`
	TypedFortniteItemsWithDisplayAssets map[string][]*FortniteItem `json:"-"`
	SnowVariantTokens map[string]SnowCosmeticVariantToken `json:"-"`
}

func NewExternalDataClient() *ExternalDataClient {
	return &ExternalDataClient{
		h: &http.Client{},
		FortniteItems: make(map[string]*FortniteItem),
		FortniteSets: make(map[string]*FortniteSet),
		FortniteItemsWithDisplayAssets: make(map[string]*FortniteItem),
		FortniteItemsWithFeaturedImage: []*FortniteItem{},
		TypedFortniteItems: make(map[string][]*FortniteItem),
		TypedFortniteItemsWithDisplayAssets: make(map[string][]*FortniteItem),
		SnowVariantTokens: make(map[string]SnowCosmeticVariantToken),
	}
}

func (c *ExternalDataClient) LoadExternalData() {
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

	displayAssets := storage.HttpAsset[[]string]("QKnwROGzQjYm1W9xu9uL3VrbSA0tnVj6NJJtEChUdAb3DF8uN.json")
	if displayAssets == nil {
		return
	}

	for _, displayAsset := range *displayAssets {
		c.AddDisplayAssetToItem(displayAsset)
	}

	variantTokens := storage.HttpAsset[map[string]SnowCosmeticVariantToken]("QF3nHCFt1vhELoU4q1VKTmpxnk20c2iAiBEBzlbzQAY.json")
	if variantTokens == nil {
		return
	}
	c.SnowVariantTokens = *variantTokens

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

func (c *ExternalDataClient) LoadItem(item *FortniteItem) {
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

func (c *ExternalDataClient) AddBackpackToItem(backpack *FortniteItem) {
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

func (c *ExternalDataClient) AddDisplayAssetToItem(displayAsset string) {
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

func (c *ExternalDataClient) AddNumericStylesToItem(item *FortniteItem) {
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

func PreloadCosmetics() error {
	External = NewExternalDataClient()
	External.LoadExternalData()

	aid.Print("(snow) " + fmt.Sprint(len(External.FortniteItems)) + " cosmetics loaded from fortnite-api.com")
	return nil
}

func GetItemByShallowID(shallowID string) *FortniteItem {
	for _, item := range External.TypedFortniteItems["AthenaCharacter"] {
		if strings.Contains(item.ID, shallowID) {
			return item
		}
	}

	return nil
}

func GetRandomItemWithDisplayAsset() *FortniteItem {
	items := External.FortniteItemsWithDisplayAssets
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
	
	for t, items := range External.TypedFortniteItemsWithDisplayAssets {
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
	for _, set := range External.FortniteSets {
		sets = append(sets, *set)
	}

	slices.SortFunc[[]FortniteSet](sets, func(a, b FortniteSet) int {
		return strings.Compare(a.BackendName, b.BackendName)
	})

	return &sets[aid.RandomInt(0, len(sets))]
}