package storage

import (
	"time"

	"github.com/lib/pq"
)

type Tabler interface {
	TableName() string
}

type DB_Person struct {
	ID string
	DisplayName string
	RefundTickets int
	Permissions int64
	Receipts []DB_Receipt `gorm:"foreignkey:PersonID"`
	Profiles []DB_Profile `gorm:"foreignkey:PersonID"`
	Stats []DB_SeasonStat `gorm:"foreignkey:PersonID"`
	Discord DB_DiscordPerson `gorm:"foreignkey:PersonID"`
	BanHistory []DB_BanStatus `gorm:"foreignkey:PersonID"` 
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

type DB_Receipt struct {
	ID string `gorm:"primary_key"`
	PersonID string `gorm:"index"`
	OfferID string
	PurchaseDate int64
	TotalPaid int
	State string
	Loot []DB_ReceiptLoot `gorm:"foreignkey:ReceiptID"`
}

func (DB_Receipt) TableName() string {
	return "Receipts"
}

type DB_ReceiptLoot struct {
	ID string `gorm:"primary_key"`
	ReceiptID string `gorm:"index"`
	TemplateID string
	Quantity int
	ProfileType string
}

func (DB_ReceiptLoot) TableName() string {
	return "ReceiptLoot"
}

type DB_Profile struct {
	ID string `gorm:"primary_key"`
	PersonID string `gorm:"index"`
	Items []DB_Item `gorm:"foreignkey:ProfileID"`
	Gifts []DB_Gift `gorm:"foreignkey:ProfileID"`
	Quests []DB_Quest `gorm:"foreignkey:ProfileID"`
	Attributes []DB_Attribute `gorm:"foreignkey:ProfileID"`
	Loadouts []DB_Loadout `gorm:"foreignkey:ProfileID"`
	Purchases []DB_Purchase `gorm:"foreignkey:ProfileID"`
	VariantTokens []DB_VariantToken `gorm:"foreignkey:ProfileID"`
	Type string
	Revision int
}

func (DB_Profile) TableName() string {
	return "Profiles"
}

type DB_Attribute struct {
	ID string `gorm:"primary_key"`
	ProfileID string
	Key string
	ValueJSON string
	Type string
}

func (DB_Attribute) TableName() string {
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

type DB_Purchase struct {
	ID string `gorm:"primary_key"`
	ProfileID string `gorm:"index"`
	OfferID string
	PurchaseDate int64
	FreeRefundExpiry int64
	RefundExpiry int64
	RefundedAt int64
	TotalPaid int
	Loot []DB_PurchaseLoot `gorm:"foreignkey:PurchaseID"`
}

func (DB_Purchase) TableName() string {
	return "Purchases"
}

type DB_PurchaseLoot struct {
	ID string `gorm:"primary_key"`
	PurchaseID string `gorm:"index"`
	TemplateID string
	Quantity int
	ProfileType string
}

func (DB_PurchaseLoot) TableName() string {
	return "PurchaseLoot"
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
	Loot []DB_GiftLoot `gorm:"foreignkey:GiftID"`
}

func (DB_Gift) TableName() string {
	return "Gifts"
}

type DB_GiftLoot struct {
	ID string `gorm:"primary_key"`
	GiftID string `gorm:"index"`
	TemplateID string
	Quantity int
	ProfileType string
}

func (DB_GiftLoot) TableName() string {
	return "GiftLoot"
}

type DB_VariantToken struct {
	ID string `gorm:"primary_key"`
	ProfileID string `gorm:"index"`
	TemplateID string
	Name string
	AutoEquipOnGrant bool
	CreateGiftboxOnGrant bool
	MarkItemUnseenOnGrant bool
	VariantGrants []DB_VariantTokenGrant `gorm:"foreignkey:VariantTokenID"`
}

func (DB_VariantToken) TableName() string {
	return "VariantTokens"
}

type DB_VariantTokenGrant struct {
	ID string `gorm:"primary_key"`
	VariantTokenID string `gorm:"index"`
	Channel string
	Value string
}

func (DB_VariantTokenGrant) TableName() string {
	return "VariantTokenGrants"
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
	Season int
	SeasonXP int
	BookXP int
	BookPurchased bool
	Hype int
}

func (DB_SeasonStat) TableName() string {
	return "Stats"
}

type DB_BanStatus struct {
	ID string `gorm:"primary_key"`
	PersonID string `gorm:"index"`
	IssuedBy string
	Expiry time.Time
	Reason string
}

func (DB_BanStatus) TableName() string {
	return "Bans"
}