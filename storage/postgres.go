package storage

import (
	"github.com/ectrc/snow/aid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresStorage struct {
	Postgres *gorm.DB
}

func NewPostgresStorage() *PostgresStorage {
	l := logger.Default
	if aid.Config.Output.Level == "time" {
		l = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(aid.Config.Database.URI), &gorm.Config{
		Logger: l,
	})
	if err != nil {
		panic(err)
	}

	return &PostgresStorage{
		Postgres: db,
	}
}

func (s *PostgresStorage) Migrate(table interface{}, tableName string) {
	s.Postgres.Table(tableName).AutoMigrate(table)
}

func (s *PostgresStorage) MigrateAll() {
	s.Migrate(&DB_Person{}, "Persons")
	s.Migrate(&DB_Profile{}, "Profiles")
	s.Migrate(&DB_Item{}, "Items")
	s.Migrate(&DB_Gift{}, "Gifts")
	s.Migrate(&DB_Quest{}, "Quests")
	s.Migrate(&DB_Loadout{}, "Loadouts")
	s.Migrate(&DB_Loot{}, "Loot")
	s.Migrate(&DB_VariantChannel{}, "Variants")
	s.Migrate(&DB_PAttribute{}, "Attributes")
	s.Migrate(&DB_TemporaryCode{}, "Exchanges")
	s.Migrate(&DB_DiscordPerson{}, "Discords")
	s.Migrate(&DB_SeasonStat{}, "Stats")
}

func (s *PostgresStorage) DropTables() {
	s.Postgres.Exec(`DROP SCHEMA public CASCADE; CREATE SCHEMA public; GRANT ALL ON SCHEMA public TO postgres; GRANT ALL ON SCHEMA public TO public;`)
}

func (s *PostgresStorage) GetPerson(personId string) *DB_Person {
	var dbPerson DB_Person
	s.Postgres.
		Model(&DB_Person{}).
		Preload("Profiles").
		Preload("Profiles.Loadouts").
		Preload("Profiles.Items.Variants").
		Preload("Profiles.Gifts.Loot").
		Preload("Profiles.Attributes").
		Preload("Profiles.Items").
		Preload("Profiles.Gifts").
		Preload("Profiles.Quests").
		Preload("Discord").
		Where("id = ?", personId).
		Find(&dbPerson)

	if dbPerson.ID == "" {
		return nil
	}

	return &dbPerson
}

func (s *PostgresStorage) GetPersonByDisplay(displayName string) *DB_Person {
	var dbPerson DB_Person
	s.Postgres.
		Model(&DB_Person{}).
		Preload("Profiles").
		Preload("Profiles.Loadouts").
		Preload("Profiles.Items.Variants").
		Preload("Profiles.Gifts.Loot").
		Preload("Profiles.Attributes").
		Preload("Profiles.Items").
		Preload("Profiles.Gifts").
		Preload("Profiles.Quests").
		Preload("Discord").
		Where("display_name = ?", displayName).
		Find(&dbPerson)

	if dbPerson.ID == "" {
		return nil
	}

	return &dbPerson
}

func (s *PostgresStorage) GetPersonByDiscordID(discordId string) *DB_Person {
	var discordEntry DB_DiscordPerson
	s.Postgres.Model(&DB_DiscordPerson{}).Where("id = ?", discordId).Find(&discordEntry)

	if discordEntry.ID == "" {
		return nil
	}

	return s.GetPerson(discordEntry.PersonID)
}

func (s *PostgresStorage) GetAllPersons() []*DB_Person {
	var dbPersons []*DB_Person

	s.Postgres.
		Model(&DB_Person{}).
		Preload("Profiles").
		Preload("Profiles.Loadouts").
		Preload("Profiles.Items.Variants").
		Preload("Profiles.Gifts.Loot").
		Preload("Profiles.Attributes").
		Preload("Profiles.Items").
		Preload("Profiles.Gifts").
		Preload("Profiles.Quests").
		Preload("Discord").
		Find(&dbPersons)

	return dbPersons
}

func (s *PostgresStorage) SavePerson(person *DB_Person) {
	s.Postgres.Save(person)
}

func (s *PostgresStorage) DeletePerson(personId string) {
	s.Postgres.Delete(&DB_Person{}, "id = ?", personId)
}

func (s *PostgresStorage) SaveProfile(profile *DB_Profile) {
	s.Postgres.Save(profile)
}

func (s *PostgresStorage) DeleteProfile(profileId string) {
	s.Postgres.Delete(&DB_Profile{}, "id = ?", profileId)
}

func (s *PostgresStorage) SaveItem(item *DB_Item) {
	s.Postgres.Save(item)
}

func (s *PostgresStorage) DeleteItem(itemId string) {
	s.Postgres.Delete(&DB_Item{}, "id = ?", itemId)
}

func (s *PostgresStorage) SaveVariant(variant *DB_VariantChannel) {
	s.Postgres.Save(variant)
}

func (s *PostgresStorage) DeleteVariant(variantId string) {
	s.Postgres.Delete(&DB_VariantChannel{}, "id = ?", variantId)
}

func (s *PostgresStorage) SaveQuest(quest *DB_Quest) {
	s.Postgres.Save(quest)
}

func (s *PostgresStorage) DeleteQuest(questId string) {
	s.Postgres.Delete(&DB_Quest{}, "id = ?", questId)
}

func (s *PostgresStorage) SaveLoot(loot *DB_Loot) {
	s.Postgres.Save(loot)
}

func (s *PostgresStorage) DeleteLoot(lootId string) {
	s.Postgres.Delete(&DB_Loot{}, "id = ?", lootId)
}

func (s *PostgresStorage) SaveGift(gift *DB_Gift) {
	s.Postgres.Save(gift)
}

func (s *PostgresStorage) DeleteGift(giftId string) {
	s.Postgres.Delete(&DB_Gift{}, "id = ?", giftId)
}

func (s *PostgresStorage) SaveAttribute(attribute *DB_PAttribute) {
	s.Postgres.Save(attribute)
}

func (s *PostgresStorage) DeleteAttribute(attributeId string) {
	s.Postgres.Delete(&DB_PAttribute{}, "id = ?", attributeId)
}

func (s *PostgresStorage) SaveLoadout(loadout *DB_Loadout) {
	s.Postgres.Save(loadout)
}

func (s *PostgresStorage) DeleteLoadout(loadoutId string) {
	s.Postgres.Delete(&DB_Loadout{}, "id = ?", loadoutId)
}

func (s *PostgresStorage) SaveTemporaryCode(code *DB_TemporaryCode) {
	s.Postgres.Save(code)
}

func (s *PostgresStorage) DeleteTemporaryCode(codeId string) {
	s.Postgres.Delete(&DB_TemporaryCode{}, "id = ?", codeId)
}

func (s *PostgresStorage) SaveDiscordPerson(discordPerson *DB_DiscordPerson) {
	s.Postgres.Save(discordPerson)
}

func (s *PostgresStorage) DeleteDiscordPerson(discordPersonId string) {
	s.Postgres.Delete(&DB_DiscordPerson{}, "id = ?", discordPersonId)
}