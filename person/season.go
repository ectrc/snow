package person

import (
	"github.com/ectrc/snow/storage"
	"github.com/google/uuid"
)

type SeasonStats struct {
	ID string
	PersonID string
	Season int
	SeasonXP int
	BookXP int
	BookPurchased bool
	Hype int
}

func NewSeasonStats(season int) *SeasonStats {
	return &SeasonStats{
		ID: uuid.New().String(),
		Season: season,
	}
}

func (s *SeasonStats) ToDatabase(personId string) *storage.DB_SeasonStat {
	return &storage.DB_SeasonStat{
		ID: s.ID,
		PersonID: personId,
		Season: s.Season,
		SeasonXP: s.SeasonXP,
		BookXP: s.BookXP,
		BookPurchased: s.BookPurchased,
		Hype: s.Hype,
	}
}

func (s *SeasonStats) Save() {
	storage.Repo.SaveSeasonStats(s.ToDatabase(s.PersonID))
}

func (s *SeasonStats) Delete() {
	storage.Repo.DeleteSeasonStats(s.ID)
}

func FromDatabaseSeasonStats(db storage.DB_SeasonStat) *SeasonStats {
	return &SeasonStats{
		ID: db.ID,
		PersonID: db.PersonID,
		Season: db.Season,
		SeasonXP: db.SeasonXP,
		BookXP: db.BookXP,
		BookPurchased: db.BookPurchased,
		Hype: db.Hype,
	}
}