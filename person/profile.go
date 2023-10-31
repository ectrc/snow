package person

import (
	"sync"

	"github.com/ectrc/snow/storage"
	"github.com/google/uuid"
	"github.com/r3labs/diff/v3"
)

type Profile struct {
	ID         string
	Items      *ItemMutex
	Gifts      *GiftMutex
	Quests		 *QuestMutex
	Attributes *sync.Map
	Type			 string
	Changes 	 []diff.Change
}

func NewProfile(profile string) *Profile {
	return &Profile{
		ID:         uuid.New().String(),
		Items:      NewItemMutex(profile),
		Gifts:      NewGiftMutex(),
		Quests:		  NewQuestMutex(),
		Attributes: &sync.Map{},
		Type:			 profile,
	}
}

func FromDatabaseProfile(profile *storage.DB_Profile) *Profile {
	items := NewItemMutex(profile.Type)
	gifts := NewGiftMutex()
	quests := NewQuestMutex()

	for _, item := range profile.Items {
		items.AddItem(FromDatabaseItem(&item, &profile.Type))
	}

	for _, gift := range profile.Gifts {
		gifts.AddGift(FromDatabaseGift(&gift))
	}

	for _, quest := range profile.Quests {
		quests.AddQuest(FromDatabaseQuest(&quest, &profile.Type))
	}

	attributes := &sync.Map{}
	for _, attribute := range profile.Attributes {
		attributes.Store(attribute.Key, attribute.Value)
	}

	return &Profile{
		ID:         profile.ID,
		Items:      items,
		Gifts:      gifts,
		Quests:			quests,
		Attributes: attributes,
		Type:			  profile.Type,
	}
}

func (p *Profile) Save() {
	//storage.Repo.SaveProfile(p.ToDatabase())
}

func (p *Profile) Snapshot() *ProfileSnapshot {
	items := map[string]ItemSnapshot{}
	gifts := map[string]GiftSnapshot{}
	quests := map[string]Quest{}
	attributes := map[string]string{}

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

	p.Attributes.Range(func(key, value interface{}) bool {
		attributes[key.(string)] = value.(string)
		return true
	})

	return &ProfileSnapshot{
		ID:         p.ID,
		Items:      items,
		Gifts:      gifts,
		Quests:			quests,
		Attributes:	attributes,
	}
}

func (p *Profile) Diff(snapshot *ProfileSnapshot) []diff.Change {
	changes, _ := diff.Diff(p.Snapshot(), snapshot)
	p.Changes = changes
	return changes
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
		Dances: []string{"", "", "", "", "", ""},
		ItemWraps: []string{"", "", "", "", "", "", ""},
		LoadingScreen: "",
		SkyDiveContrail: "",
		MusicPack: "",
		BannerIcon: "",
		BannerColor: "",
	}
}

func FromDatabaseLoadout(loadout *storage.DB_Loadout) *Loadout {
	return &Loadout{
		ID: loadout.ID,
		Character: loadout.Character,
		Backpack: loadout.Backpack,
		Pickaxe: loadout.Pickaxe,
		Glider: loadout.Glider,
		Dances: loadout.Dances,
		ItemWraps: loadout.ItemWraps,
		LoadingScreen: loadout.LoadingScreen,
		SkyDiveContrail: loadout.SkyDiveContrail,
		MusicPack: loadout.MusicPack,
		BannerIcon: loadout.BannerIcon,
		BannerColor: loadout.BannerColor,
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