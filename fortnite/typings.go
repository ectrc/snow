package fortnite

import (
	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/person"
)

type APICosmeticDefinitionVariant struct {
	Tag   string `json:"tag"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type APICosmeticDefinitionVariantChannel struct {
	Channel string                         `json:"channel"`
	Type    string                         `json:"type"`
	Options []APICosmeticDefinitionVariant `json:"options"`
}

type APICosmeticDefinition struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        struct {
		Value        string `json:"value"`
		DisplayValue string `json:"displayValue"`
		BackendValue string `json:"backendValue"`
	} `json:"type"`
	Rarity struct {
		Value        string `json:"value"`
		DisplayValue string `json:"displayValue"`
		BackendValue string `json:"backendValue"`
	} `json:"rarity"`
	Series struct {
		Value        string `json:"value"`
		Image        string `json:"image"`
		BackendValue string `json:"backendValue"`
	} `json:"series"`
	Set struct {
		Value        string `json:"value"`
		Text         string `json:"text"`
		BackendValue string `json:"backendValue"`
	} `json:"set"`
	Introduction struct {
		Chapter      string `json:"chapter"`
		Season       string `json:"season"`
		Text         string `json:"text"`
		BackendValue int    `json:"backendValue"`
	} `json:"introduction"`
	Images struct {
		Icon      string            `json:"icon"`
		Featured  string            `json:"featured"`
		SmallIcon string            `json:"smallIcon"`
		Other     map[string]string `json:"other"`
	} `json:"images"`
	Variants            []APICosmeticDefinitionVariantChannel `json:"variants"`
	GameplayTags        []string                              `json:"gameplayTags"`
	SearchTags          []string                              `json:"searchTags"`
	MetaTags            []string                              `json:"metaTags"`
	ShowcaseVideo       string                                `json:"showcaseVideo"`
	DynamicPakID        string                                `json:"dynamicPakId"`
	DisplayAssetPath    string                                `json:"displayAssetPath"`
	NewDisplayAssetPath string
	ItemPreviewHeroPath string                 `json:"itemPreviewHeroPath"`
	BackpackDefinition  *APICosmeticDefinition `json:"backpack"`
	Path                string                 `json:"path"`
	Added               string                 `json:"added"`
	ShopHistory         []string               `json:"shopHistory"`
	BattlePass          bool                   `json:"battlePass"`
}

type APISetDefinition struct {
	BackendName string                   `json:"backendName"`
	DisplayName string                   `json:"displayName"`
	Items       []*APICosmeticDefinition `json:"items"`
}

type APICosmeticsResponse struct {
	Status int                     `json:"status"`
	Data   []APICosmeticDefinition `json:"data"`
}

type SnowCosmeticVariantDefinition struct {
	Grants []struct {
		Channel string `json:"channel"`
		Value   string `json:"value"`
	} `json:"grants"`
	Item   string `json:"item"`
	Name   string `json:"name"`
	Gift   bool   `json:"gift"`
	Equip  bool   `json:"equip"`
	Unseen bool   `json:"unseen"`
}

type FortniteVariantToken struct {
	Grants []struct {
		Channel string `json:"channel"`
		Value   string `json:"value"`
	} `json:"grants"`
	Item   *APICosmeticDefinition `json:"item"`
	Name   string                 `json:"name"`
	Gift   bool                   `json:"gift"`
	Equip  bool                   `json:"equip"`
	Unseen bool                   `json:"unseen"`
}

type ItemGrant struct {
	TemplateID  string
	Quantity    int
	ProfileType string
}

func NewItemGrant(templateId string, quantity int) *ItemGrant {
	return &ItemGrant{
		TemplateID: templateId,
		Quantity:   quantity,
	}
}

type LootResultLoot struct {
	TemplateID      string
	ItemID          string
	Quantity        int
	ItemProfileType string
}

type LootResult struct {
	Items []*LootResultLoot
}

func NewLootResult() *LootResult {
	return &LootResult{
		Items: make([]*LootResultLoot, 0),
	}
}

func (l *LootResult) AddItem(i *person.Item) {
	l.Items = append(l.Items, &LootResultLoot{
		TemplateID:      i.TemplateID,
		ItemID:          i.ID,
		Quantity:        i.Quantity,
		ItemProfileType: i.ProfileType,
	})
}

func (l *LootResult) GenerateFortniteLootResultEntry() []aid.JSON {
	loot := []aid.JSON{}

	for _, item := range l.Items {
		loot = append(loot, aid.JSON{
			"itemType":    item.TemplateID,
			"itemGuid":    item.ItemID,
			"itemProfile": item.ItemProfileType,
			"quantity":    item.Quantity,
		})
	}

	return loot
}