package person

import (
	"fmt"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/storage"
	"github.com/google/uuid"
	"github.com/r3labs/diff/v3"
)

type Profile struct {
	ID string
	PersonID string
	Items *ItemMutex
	Gifts *GiftMutex
	Quests *QuestMutex
	Attributes *AttributeMutex
	Loadouts *LoadoutMutex
	Type string
	Revision int
	Changes []interface{}
}

func NewProfile(profile string) *Profile {
	id := uuid.New().String()
	return &Profile{
		ID: id,
		PersonID: "",
		Items: NewItemMutex(&storage.DB_Profile{ID: id, Type: profile}),
		Gifts: NewGiftMutex(&storage.DB_Profile{ID: id, Type: profile}),
		Quests:	NewQuestMutex(&storage.DB_Profile{ID: id, Type: profile}),
		Attributes: NewAttributeMutex(&storage.DB_Profile{ID: id, Type: profile}),
		Loadouts: NewLoadoutMutex(&storage.DB_Profile{ID: id, Type: profile}),
		Type:			 profile,
		Revision: 0,
		Changes: 	 []interface{}{},
	}
}

func FromDatabaseProfile(profile *storage.DB_Profile) *Profile {
	items := NewItemMutex(profile)
	gifts := NewGiftMutex(profile)
	quests := NewQuestMutex(profile)
	attributes := NewAttributeMutex(profile)
	loadouts := NewLoadoutMutex(profile)

	for _, item := range profile.Items {
		items.AddItem(FromDatabaseItem(&item))
	}

	for _, gift := range profile.Gifts {
		gifts.AddGift(FromDatabaseGift(&gift))
	}

	for _, quest := range profile.Quests {
		quests.AddQuest(FromDatabaseQuest(&quest))
	}

	for _, loadout := range profile.Loadouts {
		loadouts.AddLoadout(FromDatabaseLoadout(&loadout))
	}

	for _, attribute := range profile.Attributes {
		parsed := FromDatabaseAttribute(&attribute)
		if parsed == nil {
			fmt.Printf("error getting attribute from database")
			continue
		}

		attributes.AddAttribute(parsed)
	}

	return &Profile{
		ID: profile.ID,
		PersonID: profile.PersonID,
		Items: items,
		Gifts: gifts,
		Quests: quests,
		Attributes: attributes,
		Loadouts: loadouts,
		Type: profile.Type,
		Revision: profile.Revision,
		Changes: []interface{}{},
	}
}

func (p *Profile) GenerateFortniteProfileEntry() aid.JSON {
	items := aid.JSON{}
	attributes := aid.JSON{}

	p.Items.RangeItems(func(id string, item *Item) bool {
		items[id] = item.GenerateFortniteItemEntry()
		return true
	})

	p.Quests.RangeQuests(func(id string, quest *Quest) bool {
		items[id] = quest.GenerateFortniteQuestEntry()
		return true
	})

	p.Gifts.RangeGifts(func(id string, gift *Gift) bool {
		items[id] = gift.GenerateFortniteGiftEntry()
		return true
	})

	p.Attributes.RangeAttributes(func(id string, attribute *Attribute) bool {
		attributes[attribute.Key] = aid.JSONParse(attribute.ValueJSON)
		return true
	})

	p.Loadouts.RangeLoadouts(func(id string, loadout *Loadout) bool {
		items[id] = loadout.GenerateFortniteLoadoutEntry()
		return true
	})

	return aid.JSON{
		"profileId": p.Type,
		"accountId": p.PersonID,
		"rvn": p.Revision,
		"commandRevision": p.Revision,
		"wipeNumber": 0,
		"version": "",
		"items": items,
		"stats": aid.JSON{
			"attributes": attributes,
		},
	}
}

func (p *Profile) Save() {
	// storage.Repo.SaveProfile(p.ToDatabase())
}

func (p *Profile) Snapshot() *ProfileSnapshot {
	items := map[string]ItemSnapshot{}
	gifts := map[string]GiftSnapshot{}
	quests := map[string]Quest{}
	attributes := map[string]Attribute{}
	loadouts := map[string]Loadout{}

	p.Items.RangeItems(func(id string, item *Item) bool {
		items[id] = item.Snapshot()
		return true
	})

	p.Gifts.RangeGifts(func(id string, gift *Gift) bool {
		gifts[id] = gift.Snapshot()
		return true
	})

	p.Quests.RangeQuests(func(id string, quest *Quest) bool {
		quests[id] = *quest
		return true
	})

	p.Attributes.RangeAttributes(func(key string, attribute *Attribute) bool {
		attributes[key] = *attribute
		return true
	})

	p.Loadouts.RangeLoadouts(func(id string, loadout *Loadout) bool {
		loadouts[id] = *loadout
		return true
	})

	return &ProfileSnapshot{
		ID: p.ID,
		Items: items,
		Gifts: gifts,
		Quests: quests,
		Attributes: attributes,
		Loadouts: loadouts,
		Type: p.Type,
		Revision: p.Revision,
	}
}

func (p *Profile) Diff(b *ProfileSnapshot) []diff.Change {
	changes, err := diff.Diff(*b, *p.Snapshot())
	if err != nil {
		fmt.Printf("error diffing profile: %v\n", err)
		return nil
	}

	aid.PrintJSON(changes)

	for _, change := range changes {
		switch change.Path[0] {
		case "Items":
			if change.Type == "create" && change.Path[2] == "ID" {
				p.CreateItemAddedChange(p.Items.GetItem(change.Path[1]))
			}

			if change.Type == "delete" && change.Path[2] == "ID" {
				p.CreateItemRemovedChange(change.Path[1])
			}

			if change.Type == "update" && change.Path[2] == "Quantity" {
				p.CreateItemQuantityChangedChange(p.Items.GetItem(change.Path[1]))
			}

			if change.Type == "update" && change.Path[2] != "Quantity" {
				p.CreateItemAttributeChangedChange(p.Items.GetItem(change.Path[1]), change.Path[2])
			}
		case "Quests":
			if change.Type == "create" && change.Path[2] == "ID" {
				p.CreateQuestAddedChange(p.Quests.GetQuest(change.Path[1]))
			}

			if change.Type == "delete" && change.Path[2] == "ID" {
				p.CreateItemRemovedChange(change.Path[1])
			}
		case "Gifts":
			if change.Type == "create" && change.Path[2] == "ID" {
				p.CreateGiftAddedChange(p.Gifts.GetGift(change.Path[1]))
			}

			if change.Type == "delete" && change.Path[2] == "ID" {
				p.CreateItemRemovedChange(change.Path[1])
			}
		case "Attributes":
			if change.Type == "create" && change.Path[2] == "ID" {
				p.CreateStatModifiedChange(p.Attributes.GetAttribute(change.Path[1]))
			}

			if change.Type == "update" && change.Path[2] == "ValueJSON" {
				p.CreateStatModifiedChange(p.Attributes.GetAttribute(change.Path[1]))
			}
		case "Loadouts":
			if change.Type == "create" && change.Path[2] == "ID" {
				p.CreateLoadoutAddedChange(p.Loadouts.GetLoadout(change.Path[1]))
			}

			if change.Type == "delete" && change.Path[2] == "ID" {
				p.CreateLoadoutRemovedChange(change.Path[1])
			}

			if change.Type == "update" && change.Path[2] != "ID" {
				p.CreateLoadoutChangedChange(p.Loadouts.GetLoadout(change.Path[1]), change.Path[2])
			}
		}
	}

	return changes
}

func (p *Profile) CreateAttribute(key string, value interface{}) *Attribute {
	p.Attributes.AddAttribute(NewAttribute(key, value))
	return p.Attributes.GetAttribute(key)
}

func (p *Profile) CreateStatModifiedChange(attribute *Attribute) {
	if attribute == nil {
		fmt.Println("error getting attribute from profile", attribute.ID)
		return
	}

	p.Changes = append(p.Changes, StatModified{
		ChangeType: "statModified",
		Name: attribute.Key,
		Value: aid.JSONParse(attribute.ValueJSON),
	})
}

func (p *Profile) CreateGiftAddedChange(gift *Gift) {
	if gift == nil {
		fmt.Println("error getting gift from profile", gift.ID)
		return
	}

	p.Changes = append(p.Changes, ItemAdded{
		ChangeType: "itemAdded",
		ItemId: gift.ID,
		Item: gift.GenerateFortniteGiftEntry(),
	})
}

func (p *Profile) CreateQuestAddedChange(quest *Quest) {
	if quest == nil {
		fmt.Println("error getting quest from profile", quest.ID)
		return
	}

	p.Changes = append(p.Changes, ItemAdded{
		ChangeType: "itemAdded",
		ItemId: quest.ID,
		Item: quest.GenerateFortniteQuestEntry(),
	})
}

func (p *Profile) CreateItemAddedChange(item *Item) {
	if item == nil {
		fmt.Println("error getting item from profile", item.ID)
		return
	}

	p.Changes = append(p.Changes, ItemAdded{
		ChangeType: "itemAdded",
		ItemId: item.ID,
		Item: item.GenerateFortniteItemEntry(),
	})
}

func (p *Profile) CreateItemRemovedChange(itemId string) {
	p.Changes = append(p.Changes, ItemRemoved{
		ChangeType: "itemRemoved",
		ItemId: itemId,
	})
}

func (p *Profile) CreateItemQuantityChangedChange(item *Item) {
	if item == nil {
		fmt.Println("error getting item from profile", item.ID)
		return
	}

	p.Changes = append(p.Changes, ItemQuantityChanged{
		ChangeType: "itemQuantityChanged",
		ItemId: item.ID,
		Quantity: item.Quantity,
	})
}

func (p *Profile) CreateItemAttributeChangedChange(item *Item, attribute string) {
	if item == nil {
		fmt.Println("error getting item from profile", item.ID)
		return
	}

	lookup := map[string]string{
		"Favorite": "favorite",
		"HasSeen": "item_seen",
		"Variants": "variants",
	}

	p.Changes = append(p.Changes, ItemAttributeChanged{
		ChangeType: "itemAttrChanged",
		ItemId: item.ID,
		AttributeName: lookup[attribute],
		AttributeValue: item.GetAttribute(attribute),
	})
}

func (p *Profile) CreateLoadoutAddedChange(loadout *Loadout) {
	if loadout == nil {
		fmt.Println("error getting item from profile", loadout.ID)
		return
	}

	p.Changes = append(p.Changes, ItemAdded{
		ChangeType: "itemAdded",
		ItemId: loadout.ID,
		Item: loadout.GenerateFortniteLoadoutEntry(),
	})
}

func (p *Profile) CreateLoadoutRemovedChange(loadoutId string) {
	p.Changes = append(p.Changes, ItemRemoved{
		ChangeType: "itemRemoved",
		ItemId: loadoutId,
	})
}

func (p *Profile) CreateLoadoutChangedChange(loadout *Loadout, attribute string) {
	if loadout == nil {
		fmt.Println("error getting item from profile", loadout.ID)
		return
	}

	lookup := map[string]string{
		"LockerName": "locker_name",
		"BannerID": "banner_icon_template",
		"BannerColorID": "banner_color_template",
		"CharacterID": "locker_slots_data",
		"PickaxeID": "locker_slots_data",
		"BackpackID": "locker_slots_data",
		"GliderID": "locker_slots_data",
		"DanceID": "locker_slots_data",
		"ItemWrapID": "locker_slots_data",
		"ContrailID": "locker_slots_data",
		"LoadingScreenID": "locker_slots_data",
		"MusicPackID": "locker_slots_data",
	}

	p.Changes = append(p.Changes, ItemAttributeChanged{
		ChangeType: "itemAttrChanged",
		ItemId: loadout.ID,
		AttributeName: lookup[attribute],
		AttributeValue: loadout.GetAttribute(lookup[attribute]),
	})
}

func (p *Profile) CreateFullProfileUpdateChange() {
	p.Changes = []interface{}{FullProfileUpdate{
		ChangeType: "fullProfileUpdate",
		Profile: p.GenerateFortniteProfileEntry(),
	}}
}

func (p *Profile) ClearProfileChanges() {
	p.Changes = []interface{}{}
}