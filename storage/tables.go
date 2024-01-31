package storage

import "github.com/lib/pq"

type Tabler interface {
	TableName() string
}

type DB_Person struct {
	ID string
	DisplayName string
	Permissions pq.StringArray `gorm:"type:text[]"`
	IsBanned bool
	Profiles []DB_Profile `gorm:"foreignkey:PersonID"`
	Stats []DB_SeasonStat `gorm:"foreignkey:PersonID"`
	Discord DB_DiscordPerson `gorm:"foreignkey:PersonID"`
}

func (DB_Person) TableName() string {
	return "Persons"
}

type DB_Relationship struct {
	FromPersonID string `gorm:"primary_key"`
	TowardsPersonID string `gorm:"primary_key"`
	Status string
}

func (DB_Relationship) TableName() string {
	return "Relationships"
}

type DB_Profile struct {
	ID string `gorm:"primary_key"`
	PersonID string `gorm:"index"`
	Items []DB_Item `gorm:"foreignkey:ProfileID"`
	Gifts []DB_Gift `gorm:"foreignkey:ProfileID"`
	Quests []DB_Quest `gorm:"foreignkey:ProfileID"`
	Attributes []DB_PAttribute `gorm:"foreignkey:ProfileID"`
	Loadouts []DB_Loadout `gorm:"foreignkey:ProfileID"`
	Type string
	Revision int
}

func (DB_Profile) TableName() string {
	return "Profiles"
}

type DB_PAttribute struct {
	ID string `gorm:"primary_key"`
	ProfileID string
	Key string
	ValueJSON string
	Type string
}

func (DB_PAttribute) TableName() string {
	return "Attributes"
}

type DB_Loadout struct {
	ID string `gorm:"primary_key"`
	ProfileID string `gorm:"index"`
	TemplateID string
	LockerName string
	BannerID string
	BannerColorID string
	CharacterID string
	PickaxeID string
	BackpackID string
	GliderID string
	DanceID pq.StringArray `gorm:"type:text[]"`
	ItemWrapID pq.StringArray `gorm:"type:text[]"`
	ContrailID string
	LoadingScreenID string
	MusicPackID string
}

func (DB_Loadout) TableName() string {
	return "Loadouts"
}

type DB_Item struct {
	ID string `gorm:"primary_key"`
	ProfileID string `gorm:"index"`
	TemplateID string
	Quantity int
	Favorite bool
	HasSeen bool
	Variants []DB_VariantChannel `gorm:"foreignkey:ItemID"`
}

func (DB_Item) TableName() string {
	return "Items"
}

type DB_VariantChannel struct {
	ID string `gorm:"primary_key"`
	ItemID string `gorm:"index"`
	Channel string
	Owned pq.StringArray `gorm:"type:text[]"`
	Active string
}

func (DB_VariantChannel) TableName() string {
	return "Variants"
}

type DB_Quest struct {
	ID string `gorm:"primary_key"`
	ProfileID string `gorm:"index"`
	TemplateID string
	State string
	Objectives pq.StringArray `gorm:"type:text[]"`
	ObjectiveCounts pq.Int64Array `gorm:"type:bigint[]"`
	BundleID string
	ScheduleID string
}

func (DB_Quest) TableName() string {
	return "Quests"
}

type DB_Gift struct {
	ID string `gorm:"primary_key"`
	ProfileID string `gorm:"index"`
	TemplateID string
	Quantity int
	FromID string
	GiftedAt int64
	Message string
	Loot []DB_Loot `gorm:"foreignkey:GiftID"`
}

func (DB_Gift) TableName() string {
	return "Gifts"
}

type DB_Loot struct {
	ID string `gorm:"primary_key"`
	GiftID string `gorm:"index"`
	TemplateID string
	Quantity int
	ProfileType string
}

func (DB_Loot) TableName() string {
	return "Loot"
}

type DB_DiscordPerson struct {
	ID string `gorm:"primary_key"`
	PersonID string
	Username string
	Avatar string
	Banner string
	AccessToken string
	RefreshToken string
}

func (DB_DiscordPerson) TableName() string {
	return "Discords"
}

type DB_SeasonStat struct {
	ID string `gorm:"primary_key"`
	PersonID string
	Build string
	XP int
	Level int
	LevelClaimed int
	Stars int
	Tier int
	TierClaimed int
}

func (DB_SeasonStat) TableName() string {
	return "Stats"
}