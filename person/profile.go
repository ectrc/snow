package person

import (
	"fmt"

	"github.com/ectrc/snow/storage"
	"github.com/google/uuid"
	"github.com/r3labs/diff/v3"
)

type Profile struct {
	ID         string
	Items      *ItemMutex
	Gifts      *GiftMutex
	Quests		 *QuestMutex
	Attributes *AttributeMutex
	Type			 string
	Revision				 int
	Changes 	 []interface{}
}

func NewProfile(profile string) *Profile {
	return &Profile{
		ID:         uuid.New().String(),
		Items:      NewItemMutex(profile),
		Gifts:      NewGiftMutex(),
		Quests:		  NewQuestMutex(),
		Attributes: NewAttributeMutex(),
		Type:			  profile,
		Revision:   0,
	}
}

func FromDatabaseProfile(profile *storage.DB_Profile) *Profile {
	items := NewItemMutex(profile.Type)
	gifts := NewGiftMutex()
	quests := NewQuestMutex()
	attributes := NewAttributeMutex()

	for _, item := range profile.Items {
		items.AddItem(FromDatabaseItem(&item, &profile.Type))
	}

	for _, gift := range profile.Gifts {
		gifts.AddGift(FromDatabaseGift(&gift))
	}

	for _, quest := range profile.Quests {
		quests.AddQuest(FromDatabaseQuest(&quest, &profile.Type))
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
		ID:         profile.ID,
		Items:      items,
		Gifts:      gifts,
		Quests:			quests,
		Attributes: attributes,
		Type:			  profile.Type,
		Revision:   profile.Revision,
	}
}

func (p *Profile) Save() {
	//storage.Repo.SaveProfile(p.ToDatabase())
}

func (p *Profile) Snapshot() *ProfileSnapshot {
	items := map[string]ItemSnapshot{}
	gifts := map[string]GiftSnapshot{}
	quests := map[string]Quest{}
	attributes := map[string]Attribute{}

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

	return &ProfileSnapshot{
		ID:         p.ID,
		Items:      items,
		Gifts:      gifts,
		Quests:			quests,
		Attributes: attributes,
		Type:			  p.Type,
		Revision:   p.Revision,
	}
}

func (p *Profile) Diff(snapshot *ProfileSnapshot) []diff.Change {
	changes, err := diff.Diff(snapshot, p.Snapshot())
	if err != nil {
		fmt.Printf("error diffing profile: %v\n", err)
		return nil
	}

	// aid.PrintJSON(changes)

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
		}
	}

	return changes
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
		ChangeType: "itemAttributeChanged",
		ItemId: item.ID,
		AttributeName: lookup[attribute],
		AttributeValue: item.GetAttribute(attribute),
	})
}

type Loadout struct {
	ID string
	Character string
	Backpack string
	Pickaxe string
	Glider string
	Dances []string
	ItemWraps []string
	LoadingScreen string
	SkyDiveContrail string
	MusicPack string
	BannerIcon string
	BannerColor string
}

func NewLoadout() *Loadout {
	return &Loadout{
		ID: uuid.New().String(),
		Character: "",
		Backpack: "",
		Pickaxe: "",
		Glider: "",
		Dances: make([]string, 6),
		ItemWraps: make([]string, 7),
		LoadingScreen: "",
		SkyDiveContrail: "",
		MusicPack: "",
		BannerIcon: "",
		BannerColor: "",
	}
}

func FromDatabaseLoadout(l *storage.DB_Loadout) *Loadout {
	return &Loadout{
		ID: l.ID,
		Character: l.Character,
		Backpack: l.Backpack,
		Pickaxe: l.Pickaxe,
		Glider: l.Glider,
		Dances: l.Dances,
		ItemWraps: l.ItemWraps,
		LoadingScreen: l.LoadingScreen,
		SkyDiveContrail: l.SkyDiveContrail,
		MusicPack: l.MusicPack,
		BannerIcon: l.BannerIcon,
		BannerColor: l.BannerColor,
	}
}

func (l *Loadout) ToDatabase() *storage.DB_Loadout {
	return &storage.DB_Loadout{
		ID: l.ID,
		Character: l.Character,
		Backpack: l.Backpack,
		Pickaxe: l.Pickaxe,
		Glider: l.Glider,
		Dances: l.Dances,
		ItemWraps: l.ItemWraps,
		LoadingScreen: l.LoadingScreen,
		SkyDiveContrail: l.SkyDiveContrail,
		MusicPack: l.MusicPack,
		BannerIcon: l.BannerIcon,
		BannerColor: l.BannerColor,
	}
}

func (l *Loadout) Save() {
	//storage.Repo.SaveLoadout(l.ToDatabase())
}