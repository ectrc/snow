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
	KnownDisplayAssets = make(map[string]bool)
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

	aid.Print("Preloaded", len(Cosmetics.Items), "cosmetics")
	aid.Print("Preloaded", len(Cosmetics.Sets), "sets")
	aid.Print("Preloaded", len(battlePassSkins), "battle pass skins")

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

	// print the perecentage of backpacks that have a character
	aid.Print("Preloaded", len(found), "backpacks with characters", "(", float64(len(found))/float64(len(characters))*100, "% )")

	DAv2 := *storage.Asset("assets.json")
	if DAv2 == nil {
		aid.Print("Couldn't find DAv2.json")
	}

	var DAv2Data map[string]bool
	err = json.Unmarshal(DAv2, &DAv2Data)
	if err != nil {
		return err
	}

	KnownDisplayAssets = DAv2Data

	aid.Print("Preloaded", len(KnownDisplayAssets), "display assets")

	return nil
}