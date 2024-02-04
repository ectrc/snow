package fortnite

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"slices"
	"strings"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/storage"
)

type FortniteAPI struct {
	URL string
	C *http.Client
}

type FAPI_Response struct {
	Status int `json:"status"`
	Data []FAPI_Cosmetic `json:"data"`
}

type FAPI_Error struct {
	Status int `json:"status"`
	Error string `json:"error"`
}

type FAPI_Cosmetic_Variant struct {
	Channel string `json:"channel"`
	Type string `json:"type"`
	Options []FAPI_Cosmetic_VariantChannel  `json:"options"`
}

type FAPI_Cosmetic_VariantChannel struct {
	Tag string `json:"tag"`
	Name string `json:"name"`
	Image string `json:"image"`
}

type FAPI_Cosmetic struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Type struct {
		Value string `json:"value"`
		DisplayValue string `json:"displayValue"`
		BackendValue string `json:"backendValue"`
	} `json:"type"`
	Rarity struct {
		Value string `json:"value"`
		DisplayValue string `json:"displayValue"`
		BackendValue string `json:"backendValue"`
	} `json:"rarity"`
	Series struct {
		Value string `json:"value"`
		Image string `json:"image"`
		BackendValue string `json:"backendValue"`
	} `json:"series"`
	Set struct {
		Value string `json:"value"`
		Text string `json:"text"`
		BackendValue string `json:"backendValue"`
	} `json:"set"`
	Introduction struct {
		Chapter string `json:"chapter"`
		Season string `json:"season"`
		Text string `json:"text"`
		BackendValue int `json:"backendValue"`
	} `json:"introduction"`
	Images struct {
		Icon string `json:"icon"`
		Featured string `json:"featured"`
		SmallIcon string `json:"smallIcon"`
		Other map[string]string `json:"other"`
	} `json:"images"`
	Variants []FAPI_Cosmetic_Variant `json:"variants"`
	GameplayTags []string `json:"gameplayTags"`
	SearchTags []string `json:"searchTags"`
	MetaTags []string `json:"metaTags"`
	ShowcaseVideo string `json:"showcaseVideo"`
	DynamicPakID string `json:"dynamicPakId"`
	DisplayAssetPath string `json:"displayAssetPath"`
	DisplayAssetPath2 string
	ItemPreviewHeroPath string `json:"itemPreviewHeroPath"`
	Backpack string `json:"backpack"`
	Path string `json:"path"`
	Added string `json:"added"`
	ShopHistory []string `json:"shopHistory"`
	BattlePass bool `json:"battlePass"`
}

type Set struct {
	Items map[string]FAPI_Cosmetic `json:"items"`
	Name string `json:"name"`
	BackendName string `json:"backendName"`
}

type CosmeticData struct {
	Items map[string]FAPI_Cosmetic `json:"items"`
	Sets map[string]Set `json:"sets"`
}

func (c *CosmeticData) GetRandomItem() FAPI_Cosmetic {
	randomInt := rand.Intn(len(c.Items))

	i := 0
	for _, item := range c.Items {
		if i == randomInt {
			return item
		}

		i++
	}

	return c.GetRandomItem()
}

func (c *CosmeticData) GetRandomItemByType(itemType string) FAPI_Cosmetic {
	randomInt := rand.Intn(len(c.Items))

	i := 0
	for _, item := range c.Items {
		if item.Type.BackendValue != itemType {
			continue
		}

		if i == randomInt {
			return item
		}

		i++
	}

	return c.GetRandomItemByType(itemType)
}

func (c *CosmeticData) GetRandomItemByNotType(itemType string) FAPI_Cosmetic {
	randomInt := rand.Intn(len(c.Items))

	i := 0
	for _, item := range c.Items {
		if item.Type.BackendValue == itemType {
			continue
		}

		if i == randomInt {
			return item
		}

		i++
	}

	return c.GetRandomItemByNotType(itemType)
}

func (c *CosmeticData) GetRandomSet() Set {
	randomInt := rand.Intn(len(c.Sets))

	i := 0
	for _, set := range c.Sets {
		if i == randomInt {
			return set
		}

		i++
	}

	return c.GetRandomSet()
}

var EXTRA_NUMERIC_STYLES = []string{"Soccer", "Football", "ScaryBall"}

func (c *CosmeticData) AddItem(item FAPI_Cosmetic) {
	if slices.Contains(EXTRA_NUMERIC_STYLES, item.Set.BackendValue) {
		item = c.AddNumericVariantChannelToItem(item)
	}

	c.Items[item.ID] = item

	if item.Set.BackendValue != "" {
		if _, ok := Cosmetics.Sets[item.Set.BackendValue]; !ok {
			Cosmetics.Sets[item.Set.BackendValue] = Set{
				Items: make(map[string]FAPI_Cosmetic),
				Name: item.Set.Value,
				BackendName: item.Set.BackendValue,
			}
		}

		Cosmetics.Sets[item.Set.BackendValue].Items[item.ID] = item
	}
}

func (c *CosmeticData) AddNumericVariantChannelToItem(item FAPI_Cosmetic) FAPI_Cosmetic {
	owned := []FAPI_Cosmetic_VariantChannel{}
	for i := 0; i < 100; i++ {
		owned = append(owned, FAPI_Cosmetic_VariantChannel{
			Tag: fmt.Sprint(i),
		})
	}

	item.Variants = append(item.Variants, FAPI_Cosmetic_Variant{
		Channel: "Numeric",
		Type: "int",
		Options: owned,
	})

	return item
}

var (
	StaticAPI = NewFortniteAPI()
	Cosmetics = CosmeticData{
		Items: make(map[string]FAPI_Cosmetic),
		Sets: make(map[string]Set),
	}
)

func NewFortniteAPI() *FortniteAPI {
	return &FortniteAPI{
		URL: "https://fortnite-api.com",
		C: &http.Client{},
	}
}

func (f *FortniteAPI) Get(path string) (*FAPI_Response, error) {
	req, err := http.NewRequest("GET", f.URL + path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := f.C.Do(req)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data FAPI_Response
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (f *FortniteAPI) GetAllCosmetics() ([]FAPI_Cosmetic, error) {
	resp, err := f.Get("/v2/cosmetics/br")
	if err != nil {
		return nil, err
	}
	
	return resp.Data, nil
}

func (f *FortniteAPI) GetPlaylistImage(playlist string) (any, error) {
	return nil, nil
}

func PreloadCosmetics(max int) error {
	aid.Print("(external) assets from", StaticAPI.URL)

	list, err := StaticAPI.GetAllCosmetics()
	if err != nil {
		return err
	}

	for _, item := range list {
		if item.Introduction.BackendValue > max {
			continue
		}

		if len(item.ShopHistory) == 0 && item.Type.Value == "outfit" {
			item.BattlePass = true
		}

		Cosmetics.AddItem(item)
	}

	for id, item := range Cosmetics.Items {
		if item.Type.Value != "backpack" {
			continue
		}

		if item.ItemPreviewHeroPath == "" {
			continue
		}

		previewHeroPath := strings.Split(item.ItemPreviewHeroPath, "/")
		characterId := previewHeroPath[len(previewHeroPath)-1]

		character, ok := Cosmetics.Items[characterId]
		if !ok {
			continue
		}
		character.Backpack = id
		Cosmetics.AddItem(character)
	}

	assets := storage.HttpAsset("QKnwROGzQjYm1W9xu9uL3VrbSA0tnVj6NJJtEChUdAb3DF8uN.json")
	if assets == nil {
		panic("Failed to load assets")
	}

	var assetData []string
	err = json.Unmarshal(*assets, &assetData)
	if err != nil {
		return err
	}
	withDisplayAssets := 0

	for _, asset := range assetData {
		asset := strings.ReplaceAll(asset, "DAv2_", "")
		parts := strings.Split(asset, "_")

		if strings.Contains(asset, "Bundle") {
			withDisplayAssets++
			continue
		}

		switch {
		case parts[0] == "CID":
			addCharacterAsset(parts)
		case parts[0] == "Character":
			addCharacterAsset(parts)
		case parts[0] == "BID":
			addBackpackAsset(parts)
		case parts[0] == "EID":
			addEmoteAsset(parts)
		case parts[0] == "Emote":
			addEmoteAsset(parts)
		case parts[0] == "Pickaxe":
			addPickaxeAsset(parts)
		case parts[0] == "Wrap":
			addWrapAsset(parts)
		case parts[0] == "Glider":
			addGliderAsset(parts)
		case parts[0] == "MusicPack":
			addMusicAsset(parts)
		}
	}

	for _, item := range Cosmetics.Items {
		if item.DisplayAssetPath2 == "" {
			continue
		}

		withDisplayAssets++
	}
	aid.Print("(snow) preloaded", len(Cosmetics.Items), "cosmetics")
	
	return nil
}

func addCharacterAsset(parts []string) {
	character := FAPI_Cosmetic{}

	for _, item := range Cosmetics.Items {
		if item.Type.Value != "outfit" {
			continue
		}

		if parts[0] == "CID" {
			cid := ""
			if parts[1] != "A" {
				cid = parts[0] + "_" + parts[1]
			}

			if parts[1] == "A" {
				cid = parts[0] + "_A_" + parts[2]
			}

			if strings.Contains(item.ID, cid) {
				character = item
				break
			}
		}

		if parts[0] == "Character" {
			if strings.Contains(item.ID, parts[1]) {
				character = item
				break
			}
		}
	}

	if character.ID == "" {
		return
	}

	character.DisplayAssetPath2 = "DAv2_" + strings.Join(parts, "_")
	Cosmetics.AddItem(character)
}

func addBackpackAsset(parts []string) {
	backpack := FAPI_Cosmetic{}

	for _, item := range Cosmetics.Items {
		if item.Type.Value != "backpack" {
			continue
		}

		bid := ""
		if parts[1] != "A" {
			bid = parts[0] + "_" + parts[1]
		}

		if parts[1] == "A" {
			bid = parts[0] + "_A_" + parts[2]
		}

		if strings.Contains(item.ID, bid) {
			backpack = item
			break
		}
	}

	if backpack.ID == "" {
		return
	}

	backpack.DisplayAssetPath2 = "DAv2_" + strings.Join(parts, "_")
	Cosmetics.AddItem(backpack)
}

func addEmoteAsset(parts []string) {
	emote := FAPI_Cosmetic{}

	for _, item := range Cosmetics.Items {
		if item.Type.Value != "emote" {
			continue
		}

		if strings.Contains(item.ID, parts[1]) {
			emote = item
			break
		}
	}

	if emote.ID == "" {
		return
	}

	emote.DisplayAssetPath2 = "DAv2_" + strings.Join(parts, "_")
	Cosmetics.AddItem(emote)
}

func addPickaxeAsset(parts []string) {
	pickaxe := FAPI_Cosmetic{}

	for _, item := range Cosmetics.Items {
		if item.Type.Value != "pickaxe" {
			continue
		}

		pickaxeId := ""
		if parts[1] != "ID" {
			pickaxeId = parts[0] + "_" + parts[1]
		}

		if parts[1] == "ID" {
			pickaxeId = parts[0] + "_ID_" + parts[2]
		}

		if strings.Contains(item.ID, pickaxeId) {
			pickaxe = item
			break
		}
	}

	if pickaxe.ID == "" {
		return
	}

	pickaxe.DisplayAssetPath2 = "DAv2_" + strings.Join(parts, "_")
	Cosmetics.AddItem(pickaxe)
}

func addGliderAsset(parts []string) {
	glider := FAPI_Cosmetic{}

	for _, item := range Cosmetics.Items {
		if item.Type.Value != "glider" {
			continue
		}

		gliderId := ""
		if parts[1] != "ID" {
			gliderId = parts[0] + "_" + parts[1]
		}

		if parts[1] == "ID" {
			gliderId = parts[0] + "_ID_" + parts[2]
		}

		if strings.Contains(item.ID, gliderId) {
			glider = item
			break
		}
	}

	if glider.ID == "" {
		return
	}

	glider.DisplayAssetPath2 = "DAv2_" + strings.Join(parts, "_")
	Cosmetics.AddItem(glider)
}

func addWrapAsset(parts []string) {
	wrap := FAPI_Cosmetic{}

	for _, item := range Cosmetics.Items {
		if item.Type.Value != "wrap" {
			continue
		}

		if strings.Contains(item.ID, parts[1]) {
			wrap = item
			break
		}
	}

	if wrap.ID == "" {
		return
	}

	wrap.DisplayAssetPath2 = "DAv2_" + strings.Join(parts, "_")
	Cosmetics.AddItem(wrap)
}

func addMusicAsset(parts []string) {
	music := FAPI_Cosmetic{}

	for _, item := range Cosmetics.Items {
		if item.Type.Value != "music" {
			continue
		}

		if strings.Contains(item.ID, parts[1]) {
			music = item
			break
		}
	}

	if music.ID == "" {
		return
	}

	music.DisplayAssetPath2 = "DAv2_" + strings.Join(parts, "_")
	Cosmetics.AddItem(music)
}