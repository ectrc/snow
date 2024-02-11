package fortnite

type FortniteVariantChannel struct {
	Tag   string `json:"tag"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type FortniteVariant struct {
	Channel string                   `json:"channel"`
	Type    string                   `json:"type"`
	Options []FortniteVariantChannel `json:"options"`
}

type FortniteItem struct {
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
	Variants            []FortniteVariant `json:"variants"`
	GameplayTags        []string          `json:"gameplayTags"`
	SearchTags          []string          `json:"searchTags"`
	MetaTags            []string          `json:"metaTags"`
	ShowcaseVideo       string            `json:"showcaseVideo"`
	DynamicPakID        string            `json:"dynamicPakId"`
	DisplayAssetPath    string            `json:"displayAssetPath"`
	DisplayAssetPath2   string
	ItemPreviewHeroPath string        `json:"itemPreviewHeroPath"`
	Backpack            *FortniteItem `json:"backpack"`
	Path                string        `json:"path"`
	Added               string        `json:"added"`
	ShopHistory         []string      `json:"shopHistory"`
	BattlePass          bool          `json:"battlePass"`
}

type FortniteSet struct {
	BackendName string          `json:"backendName"`
	DisplayName string          `json:"displayName"`
	Items       []*FortniteItem `json:"items"`
}

type FortniteCosmeticsResponse struct {
	Status int            `json:"status"`
	Data   []FortniteItem `json:"data"`
}

type SnowCosmeticVariantToken struct {
	Grants []struct {
		Channel string `json:"channel"`
		Value   string `json:"value"`
	} `json:"grants"`
	Name   string `json:"name"`
	Gift   bool   `json:"gift"`
	Equip  bool   `json:"equip"`
	Unseen bool   `json:"unseen"`
}