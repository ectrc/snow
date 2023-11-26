package fortnite

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strings"

	"github.com/ectrc/snow/aid"
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
}

type Set struct {
	Items map[string]FAPI_Cosmetic `json:"items"`
	Name string `json:"name"`
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

func PreloadCosmetics(max int) error {
	list, err := StaticAPI.GetAllCosmetics()
	if err != nil {
		return err
	}

	for _, item := range list {
		if item.Introduction.BackendValue > max {
			continue
		}

		Cosmetics.Items[item.ID] = item

		if item.Set.BackendValue != "" {
			if _, ok := Cosmetics.Sets[item.Set.BackendValue]; !ok {
				Cosmetics.Sets[item.Set.BackendValue] = Set{
					Items: make(map[string]FAPI_Cosmetic),
					Name: item.Set.Value,
				}
			}

			Cosmetics.Sets[item.Set.BackendValue].Items[item.ID] = item
		}
	}

	aid.Print("Preloaded", len(Cosmetics.Items), "cosmetics")
	aid.Print("Preloaded", len(Cosmetics.Sets), "sets")

	notFound := make([]string, 0)
	for id, item := range Cosmetics.Items {
		if item.ItemPreviewHeroPath == "" {
			continue
		}

		if item.Type.Value != "AthenaBackpack" {
			continue
		}

		previewHeroPath := strings.Split(item.ItemPreviewHeroPath, "/")
		characterId := previewHeroPath[len(previewHeroPath)-1]

		character, ok := Cosmetics.Items[characterId]
		if !ok {
			notFound = append(notFound, characterId)
			continue
		}

		character.Backpack = id

		Cosmetics.Items[characterId] = character
		Cosmetics.Sets[character.Set.BackendValue].Items[characterId] = character
	}

	aid.Print("Could not find", len(notFound), "items with backpacks")

	return nil
}