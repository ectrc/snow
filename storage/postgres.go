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

func (s *PostgresStorage) DropTables() {
	s.Postgres.Exec(`DROP SCHEMA public CASCADE; CREATE SCHEMA public; GRANT ALL ON SCHEMA public TO postgres; GRANT ALL ON SCHEMA public TO public;`)
}