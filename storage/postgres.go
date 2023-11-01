package storage

import (
	"github.com/ectrc/snow/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresStorage struct {
	Postgres *gorm.DB
}

func NewPostgresStorage() *PostgresStorage {
	db, err := gorm.Open(postgres.Open(config.Get().Database.URI), &gorm.Config{})
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

func (s *PostgresStorage) DropTables() {
	s.Postgres.Exec(`DROP SCHEMA public CASCADE; CREATE SCHEMA public; GRANT ALL ON SCHEMA public TO postgres; GRANT ALL ON SCHEMA public TO public;`)
}

func (s *PostgresStorage) GetPerson(personId string) *DB_Person {
	var dbPerson DB_Person
	s.Postgres.
		Preload("Profiles").
		Preload("Profiles.Items.Variants").
		Preload("Profiles.Gifts.Loot").
		Preload("Profiles.Attributes").
		Preload("Profiles.Items").
		Preload("Profiles.Gifts").
		Preload("Profiles.Quests").
		Find(&dbPerson)

	if dbPerson.ID == "" {
		return nil
	}

	return &dbPerson
}

func (s *PostgresStorage) GetAllPersons() []*DB_Person {
	var dbPersons []*DB_Person

	s.Postgres.
		Preload("Profiles").
		Preload("Profiles.Items.Variants").
		Preload("Profiles.Gifts.Loot").
		Preload("Profiles.Attributes").
		Preload("Profiles.Items").
		Preload("Profiles.Gifts").
		Preload("Profiles.Quests").
		Find(&dbPersons)

	return dbPersons
}

func (s *PostgresStorage) SavePerson(person *DB_Person) {
	s.Postgres.Save(person)
}

func (s *PostgresStorage) DeleteItem(itemId string) {
	s.Postgres.Delete(&DB_Item{}, "id = ?", itemId)
}

func (s *PostgresStorage) DeleteVariant(variantId string) {
	s.Postgres.Delete(&DB_VariantChannel{}, "id = ?", variantId)
}

func (s *PostgresStorage) DeleteQuest(questId string) {
	s.Postgres.Delete(&DB_Quest{}, "id = ?", questId)
}

func (s *PostgresStorage) DeleteLoot(lootId string) {
	s.Postgres.Delete(&DB_Loot{}, "id = ?", lootId)
}

func (s *PostgresStorage) DeleteGift(giftId string) {
	s.Postgres.Delete(&DB_Gift{}, "id = ?", giftId)
}

func (s *PostgresStorage) DeleteAttribute(attributeId string) {
	s.Postgres.Delete(&DB_PAttribute{}, "id = ?", attributeId)
}