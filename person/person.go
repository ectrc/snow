package person

import (
	"github.com/ectrc/snow/storage"
	"github.com/google/uuid"
)

type Person struct {
	ID       string
	DisplayName string
	AthenaProfile  *Profile
	CommonCoreProfile *Profile
	Loadout *Loadout
}

type Option struct {
	Key string
	Value string
}

func NewPerson() *Person {
	return &Person{
		ID: uuid.New().String(),
		DisplayName: "Hello, Bully!",
		AthenaProfile: NewProfile("athena"),
		CommonCoreProfile: NewProfile("common_core"),
		Loadout: NewLoadout(),
	}
}

func FromDatabase(personId string) *Person {
	person := storage.Repo.GetPerson(personId)
	if person == nil {
		return nil
	}

	loadout := FromDatabaseLoadout(&person.Loadout)
	athenaProfile := &Profile{}
	commonCoreProfile := &Profile{}

	for _, profile := range person.Profiles {
		if profile.Type == "athena" {
			athenaProfile := FromDatabaseProfile(&profile)
			athenaProfile.ID = profile.ID
		}

		if profile.Type == "common_core" {
			commonCoreProfile.ID = profile.ID
			commonCoreProfile = FromDatabaseProfile(&profile)
		}
	}

	return &Person{
		ID: person.ID,
		DisplayName: person.DisplayName,
		AthenaProfile: athenaProfile,
		CommonCoreProfile: commonCoreProfile,
		Loadout: loadout,
	}
}

func AllFromDatabase() []*Person {
	var persons []*Person

	for _, person := range storage.Repo.GetAllPersons() {
		persons = append(persons, FromDatabase(person.ID))
	}

	return persons
}

func (p *Person) Save() {
	storage.Repo.SavePerson(p.ToDatabase())
}

func (p *Person) ToDatabase() *storage.DB_Person {
	dbPerson := storage.DB_Person{
		ID: p.ID,
		DisplayName: p.DisplayName,
		Profiles: []storage.DB_Profile{},
		Loadout: *p.Loadout.ToDatabase(),
	}

	profilesToConvert := map[string]*Profile{"athena": p.AthenaProfile, "common_core": p.CommonCoreProfile}

	for profileType, profile := range profilesToConvert {
		dbProfile := storage.DB_Profile{
			ID: profile.ID,
			PersonID: p.ID,
			Type: profileType,
			Items: []storage.DB_Item{},
			Gifts: []storage.DB_Gift{},
			Attributes: []storage.DB_PAttribute{},
		}

		profile.Items.RangeItems(func(id string, item *Item) bool {
			dbProfile.Items = append(dbProfile.Items, *item.ToDatabase(p.ID))
			return true
		})

		profile.Gifts.RangeGifts(func(id string, gift *Gift) bool {
			dbProfile.Gifts = append(dbProfile.Gifts, *gift.ToDatabase(p.ID))
			return true
		})

		profile.Quests.RangeQuests(func(id string, quest *Quest) bool {
			dbProfile.Quests = append(dbProfile.Quests, *quest.ToDatabase(p.ID))
			return true
		})

		profile.Attributes.Range(func(key, value interface{}) bool {
			dbProfile.Attributes = append(dbProfile.Attributes, storage.DB_PAttribute{
				ProfileID: profile.ID,
				Key: key.(string),
				Value: value.(string),
			})
			return true
		})

		dbPerson.Profiles = append(dbPerson.Profiles, dbProfile)
	}

	return &dbPerson
}