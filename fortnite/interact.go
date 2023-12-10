package fortnite

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
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
	Variants []struct {
		Channel string `json:"channel"`
		Type string `json:"type"`
		Options []struct {
			Tag string `json:"tag"`
			Name string `json:"name"`
			Image string `json:"image"`
		} `json:"options"`
	} `json:"variants"`
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
	aid.Print("Fortnite Assets from", StaticAPI.URL)

	list, err := StaticAPI.GetAllCosmetics()
	if err != nil {
		return err
	}

	battlePassSkins := make([]FAPI_Cosmetic, 0)
	for _, item := range list {
		if item.Introduction.BackendValue > max {
			continue
		}

		if len(item.ShopHistory) == 0 && item.Type.Value == "outfit" {
			item.BattlePass = true
			battlePassSkins = append(battlePassSkins, item)
		}

		Cosmetics.Items[item.ID] = item

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

	found := make([]string, 0)
	characters := make([]string, 0)
	for id, item := range Cosmetics.Items {
		if item.Type.Value == "outfit" {
			characters = append(characters, id)
		}

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
		Cosmetics.Items[characterId] = character

		if _, ok := Cosmetics.Sets[character.Set.BackendValue]; !ok {
			Cosmetics.Sets[character.Set.BackendValue] = Set{
				Items: make(map[string]FAPI_Cosmetic),
				Name: character.Set.Value,
				BackendName: character.Set.BackendValue,
			}
		}
		Cosmetics.Sets[character.Set.BackendValue].Items[characterId] = character
		found = append(found, id)
	}

	aid.Print("Preloaded", len(found), "backpacks with characters", "(", float64(len(found))/float64(len(characters))*100, "% ) coverage")

	assets := storage.HttpAsset("QKnwROGzQjYm1W9xu9uL3VrbSA0tnVj6NJJtEChUdAb3DF8uN.json")
	if assets == nil {
		panic("Failed to load assets")
	}

	var assetData []string
	err = json.Unmarshal(*assets, &assetData)
	if err != nil {
		return err
	}

	for _, asset := range assetData {
		asset := strings.ReplaceAll(asset, "DAv2_", "")
		parts := strings.Split(asset, "_")

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

	withDisplayAssets := 0
	for _, item := range Cosmetics.Items {
		if item.DisplayAssetPath2 == "" {
			continue
		}

		withDisplayAssets++
	}
	aid.Print("Preloaded", len(Cosmetics.Items), "cosmetics with", withDisplayAssets, "display assets", "(", float64(withDisplayAssets)/float64(len(assetData))*100, "% ) coverage" )
	
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
	Cosmetics.Items[character.ID] = character

	if _, ok := Cosmetics.Sets[character.Set.BackendValue]; !ok {
		Cosmetics.Sets[character.Set.BackendValue] = Set{
			Items: make(map[string]FAPI_Cosmetic),
			Name: character.Set.Value,
			BackendName: character.Set.BackendValue,
		}
	}
	Cosmetics.Sets[character.Set.BackendValue].Items[character.ID] = character
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
	Cosmetics.Items[backpack.ID] = backpack

	if _, ok := Cosmetics.Sets[backpack.Set.BackendValue]; !ok {
		Cosmetics.Sets[backpack.Set.BackendValue] = Set{
			Items: make(map[string]FAPI_Cosmetic),
			Name: backpack.Set.Value,
			BackendName: backpack.Set.BackendValue,
		}
	}
	Cosmetics.Sets[backpack.Set.BackendValue].Items[backpack.ID] = backpack
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
	Cosmetics.Items[emote.ID] = emote

	if _, ok := Cosmetics.Sets[emote.Set.BackendValue]; !ok {
		Cosmetics.Sets[emote.Set.BackendValue] = Set{
			Items: make(map[string]FAPI_Cosmetic),
			Name: emote.Set.Value,
			BackendName: emote.Set.BackendValue,
		}
	}
	Cosmetics.Sets[emote.Set.BackendValue].Items[emote.ID] = emote
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
	Cosmetics.Items[pickaxe.ID] = pickaxe

	if _, ok := Cosmetics.Sets[pickaxe.Set.BackendValue]; !ok {
		Cosmetics.Sets[pickaxe.Set.BackendValue] = Set{
			Items: make(map[string]FAPI_Cosmetic),
			Name: pickaxe.Set.Value,
			BackendName: pickaxe.Set.BackendValue,
		}
	}
	Cosmetics.Sets[pickaxe.Set.BackendValue].Items[pickaxe.ID] = pickaxe
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
	Cosmetics.Items[glider.ID] = glider

	if _, ok := Cosmetics.Sets[glider.Set.BackendValue]; !ok {
		Cosmetics.Sets[glider.Set.BackendValue] = Set{
			Items: make(map[string]FAPI_Cosmetic),
			Name: glider.Set.Value,
			BackendName: glider.Set.BackendValue,
		}
	}
	Cosmetics.Sets[glider.Set.BackendValue].Items[glider.ID] = glider
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
	Cosmetics.Items[wrap.ID] = wrap

	if _, ok := Cosmetics.Sets[wrap.Set.BackendValue]; !ok {
		Cosmetics.Sets[wrap.Set.BackendValue] = Set{
			Items: make(map[string]FAPI_Cosmetic),
			Name: wrap.Set.Value,
			BackendName: wrap.Set.BackendValue,
		}
	}
	Cosmetics.Sets[wrap.Set.BackendValue].Items[wrap.ID] = wrap
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
	Cosmetics.Items[music.ID] = music

	if _, ok := Cosmetics.Sets[music.Set.BackendValue]; !ok {
		Cosmetics.Sets[music.Set.BackendValue] = Set{
			Items: make(map[string]FAPI_Cosmetic),
			Name: music.Set.Value,
			BackendName: music.Set.BackendValue,
		}
	}
	Cosmetics.Sets[music.Set.BackendValue].Items[music.ID] = music
}