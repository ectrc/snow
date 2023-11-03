package storage

import "github.com/lib/pq"

type Tabler interface {
	TableName() string
}

type DB_Person struct {
	ID          string
	DisplayName string
	AccessKey  string
	Profiles    []DB_Profile `gorm:"foreignkey:PersonID"`
	Loadout     DB_Loadout   `gorm:"foreignkey:PersonID"`
}

func (DB_Person) TableName() string {
	return "Persons"
}

type DB_Loadout struct {
	ID        			string `gorm:"primary_key"`
	PersonID  			string
	Character 			string
	Backpack  			string
	Pickaxe   			string
	Glider				  string
	Dances					pq.StringArray `gorm:"type:text[]"`
	ItemWraps 			pq.StringArray `gorm:"type:text[]"`
	LoadingScreen		string
	SkyDiveContrail	string
	MusicPack				string
	BannerIcon			string
	BannerColor			string
}

func (DB_Loadout) TableName() string {
	return "Loadouts"
}

type DB_Profile struct {
	ID         string `gorm:"primary_key"`
	PersonID   string
	Items      []DB_Item       `gorm:"foreignkey:ProfileID"`
	Gifts      []DB_Gift       `gorm:"foreignkey:ProfileID"`
	Quests     []DB_Quest      `gorm:"foreignkey:ProfileID"`
	Attributes []DB_PAttribute `gorm:"foreignkey:ProfileID"`
	Type       string
	Revision	 int
}

func (DB_Profile) TableName() string {
	return "Profiles"
}

type DB_PAttribute struct {
	ID        string `gorm:"primary_key"`
	ProfileID string
	Key       string
	ValueJSON string
	Type      string
}

func (DB_PAttribute) TableName() string {
	return "Attributes"
}

type DB_Item struct {
	ID         string `gorm:"primary_key"`
	ProfileID  string
	TemplateID string
	Quantity   int
	Favorite   bool
	HasSeen    bool
	Variants   []DB_VariantChannel `gorm:"foreignkey:ItemID"`
}

func (DB_Item) TableName() string {
	return "Items"
}

type DB_VariantChannel struct {
	ID      string `gorm:"primary_key"`
	ItemID  string
	Channel string
	Owned   pq.StringArray `gorm:"type:text[]"`
	Active  string
}

func (DB_VariantChannel) TableName() string {
	return "Variants"
}

type DB_Quest struct {
	ID         string `gorm:"primary_key"`
	ProfileID  string
	TemplateID string
	State      string
	Objectives pq.StringArray `gorm:"type:text[]"`
	ObjectiveCounts pq.Int64Array `gorm:"type:bigint[]"`
	BundleID   string
	ScheduleID string
}

func (DB_Quest) TableName() string {
	return "Quests"
}

type DB_Gift struct {
	ID         string `gorm:"primary_key"`
	ProfileID  string
	TemplateID string
	Quantity   int
	FromID     string
	GiftedAt   int64
	Message    string
	Loot       []DB_Loot `gorm:"foreignkey:GiftID"`
}

func (DB_Gift) TableName() string {
	return "Gifts"
}

type DB_Loot struct {
	ID         string `gorm:"primary_key"`
	GiftID     string
	TemplateID string
	Quantity   int
	ProfileType  string
}

func (DB_Loot) TableName() string {
	return "Loot"
}