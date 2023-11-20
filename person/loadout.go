package person

import (
	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/storage"
	"github.com/google/uuid"
)

type Loadout struct {
	ID string
	PersonID string
	ProfileID string
	TemplateID string
	LockerName string
	BannerID string
	BannerColorID string
	CharacterID string
	PickaxeID string
	BackpackID string
	GliderID string
	DanceID []string
	ItemWrapID []string
	ContrailID string
	LoadingScreenID string
	MusicPackID string
}

func NewLoadout(name string, athena *Profile) *Loadout {
	character := athena.Attributes.GetAttributeByKey("favorite_character")
	pickaxe := athena.Attributes.GetAttributeByKey("favorite_pickaxe")
	backpack := athena.Attributes.GetAttributeByKey("favorite_backpack")
	glider := athena.Attributes.GetAttributeByKey("favorite_glider")
	contrail := athena.Attributes.GetAttributeByKey("favorite_skydivecontrail")
	screen := athena.Attributes.GetAttributeByKey("favorite_loadingscreen")
	music := athena.Attributes.GetAttributeByKey("favorite_musicpack")
	
	icon := athena.Attributes.GetAttributeByKey("banner_icon")
	color := athena.Attributes.GetAttributeByKey("banner_color")

	dances := athena.Attributes.GetAttributeByKey("favorite_dance")
	wraps := athena.Attributes.GetAttributeByKey("favorite_itemwraps")

	dancesReal := aid.JSONParse(dances.ValueJSON).([]any)
	wrapsReal := aid.JSONParse(wraps.ValueJSON).([]any)

	dancesValue := make([]string, len(dancesReal))
	wrapsValue := make([]string, len(wrapsReal))

	for i, v := range dancesReal {
		value, ok := v.(string)
		if !ok {
			continue
		}

		dancesValue[i] = value
	}

	for i, v := range wrapsReal {
		value, ok := v.(string)
		if !ok {
			continue
		}

		wrapsValue[i] = value
	}

	return &Loadout{
		ID: uuid.New().String(),
		PersonID: athena.ID,
		TemplateID: "CosmeticLocker:CosmeticLocker_Athena",
		LockerName: name,
		CharacterID: aid.JSONParse(character.ValueJSON).(string),
		PickaxeID: aid.JSONParse(pickaxe.ValueJSON).(string),
		BackpackID: aid.JSONParse(backpack.ValueJSON).(string),
		GliderID: aid.JSONParse(glider.ValueJSON).(string),
		ContrailID: aid.JSONParse(contrail.ValueJSON).(string),
		LoadingScreenID: aid.JSONParse(screen.ValueJSON).(string),
		MusicPackID: aid.JSONParse(music.ValueJSON).(string),
		BannerID: aid.JSONParse(icon.ValueJSON).(string),
		BannerColorID: aid.JSONParse(color.ValueJSON).(string),
		DanceID: dancesValue,
		ItemWrapID: wrapsValue,
	}
}

func FromDatabaseLoadout(loadout *storage.DB_Loadout) *Loadout {
	return &Loadout{
		ID: loadout.ID,
		ProfileID: loadout.ProfileID,
		TemplateID: loadout.TemplateID,
		LockerName: loadout.LockerName,
		BannerID: loadout.BannerID,
		BannerColorID: loadout.BannerColorID,
		CharacterID: loadout.CharacterID,
		PickaxeID: loadout.PickaxeID,
		BackpackID: loadout.BackpackID,
		GliderID: loadout.GliderID,
		DanceID: loadout.DanceID,
		ItemWrapID: loadout.ItemWrapID,
		ContrailID: loadout.ContrailID,
		LoadingScreenID: loadout.LoadingScreenID,
		MusicPackID: loadout.MusicPackID,
	}
}

func (l *Loadout) GenerateFortniteLoadoutEntry() aid.JSON {
	bannerItem := Find(l.PersonID).AthenaProfile.Items.GetItem(l.BannerID)
	if bannerItem == nil {
		bannerItem = &Item{
			TemplateID: "HomebaseBannerIcon:StandardBanner1",
		}
	}

	bannerColorItem := Find(l.PersonID).AthenaProfile.Items.GetItem(l.BannerColorID)
	if bannerColorItem == nil {
		bannerColorItem = &Item{
			TemplateID: "HomebaseBannerColor:DefaultColor1",
		}
	}

	json := aid.JSON{
		"templateId": l.TemplateID,
		"attributes": aid.JSON{
			"locker_name": l.LockerName,
			"banner_icon_template": bannerItem.TemplateID,
			"banner_color_template": bannerColorItem.TemplateID,
			"locker_slots_data": l.GenerateFortniteLockerSlotsData(),
			"item_seen": true,
		},
		"quantity": 1,
	}
	return json
}

func (l *Loadout) GetAttribute(attribute string) interface{} {
	bannerItem := Find(l.PersonID).AthenaProfile.Items.GetItem(l.BannerID)
	if bannerItem == nil {
		bannerItem = &Item{
			TemplateID: "HomebaseBannerIcon:StandardBanner1",
		}
	}

	bannerColorItem := Find(l.PersonID).AthenaProfile.Items.GetItem(l.BannerColorID)
	if bannerColorItem == nil {
		bannerColorItem = &Item{
			TemplateID: "HomebaseBannerColor:DefaultColor5",
		}
	}

	switch attribute {
	case "locker_name":
		return l.LockerName
	case "banner_icon_template":
		return bannerItem.TemplateID
	case "banner_color_template":
		return bannerColorItem.TemplateID
	case "locker_slots_data":
		return l.GenerateFortniteLockerSlotsData()
	}

	return nil
}

func (l *Loadout) GenerateFortniteLockerSlotsData() aid.JSON {
	return aid.JSON{
		"slots": aid.JSON{
			"Character": l.GetItemSlotData(l.CharacterID),
			"Backpack": l.GetItemSlotData(l.BackpackID),
			"Pickaxe": l.GetItemSlotData(l.PickaxeID),
			"Glider": l.GetItemSlotData(l.GliderID),
			"ItemWrap": l.GetItemsSlotData(l.ItemWrapID),
			"Dance": l.GetItemsSlotData(l.DanceID),
			"SkyDiveContrail": l.GetItemSlotData(l.ContrailID),
			"LoadingScreen": l.GetItemSlotData(l.LoadingScreenID),
			"MusicPack": l.GetItemSlotData(l.MusicPackID),
		},
	}
}

func (l *Loadout) GetItemSlotData(itemId string) aid.JSON {
	json := aid.JSON{
		"items": []string{},
		"activeVariants": []aid.JSON{},
	}

	person := Find(l.PersonID)
	if person == nil {
		return json
	}

	item := person.AthenaProfile.Items.GetItem(itemId)
	if item == nil {
		return json
	}

	items := json["items"].([]string)
	items = append(items, item.TemplateID)
	activeVariants := json["activeVariants"].([]aid.JSON)
	activeVariants = append(activeVariants, aid.JSON{
		"variants": item.GenerateFortniteItemVariantChannels(),
	})
	
	json["items"] = items
	json["activeVariants"] = activeVariants

	return json
}

func (l *Loadout) GetItemsSlotData(itemIds []string) aid.JSON {
	json := aid.JSON{
		"items": make([]string, len(itemIds)),
		"activeVariants": make([]aid.JSON, len(itemIds)),
	}

	person := Find(l.PersonID)
	if person == nil {
		return json
	}

	for idx, itemId := range itemIds {
		item := person.AthenaProfile.Items.GetItem(itemId)
		if item == nil {
			item = &Item{
				TemplateID: "",
				Variants: []*VariantChannel{},
			}
		}
		
		items := json["items"].([]string)
		items[idx] = item.TemplateID
		

		activeVariants := json["activeVariants"].([]aid.JSON)
		activeVariants[idx] = aid.JSON{
			"variants": []aid.JSON{},
		}

		json["items"] = items
		json["activeVariants"] = activeVariants
	}

	return json
}

func (l *Loadout) Delete() {
	storage.Repo.DeleteLoadout(l.ID)
}

func (l *Loadout) ToDatabase(profileId string) *storage.DB_Loadout {
	return &storage.DB_Loadout{
		ID: l.ID,
		ProfileID: profileId,
		TemplateID: l.TemplateID,
		LockerName: l.LockerName,
		BannerID: l.BannerID,
		BannerColorID: l.BannerColorID,
		CharacterID: l.CharacterID,
		PickaxeID: l.PickaxeID,
		BackpackID: l.BackpackID,
		GliderID: l.GliderID,
		DanceID: l.DanceID,
		ItemWrapID: l.ItemWrapID,
		ContrailID: l.ContrailID,
		LoadingScreenID: l.LoadingScreenID,
		MusicPackID: l.MusicPackID,
	}
}

func (q *Loadout) Save() {
	storage.Repo.SaveLoadout(q.ToDatabase(q.ProfileID))
}