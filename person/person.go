package person

import (
	"time"

	"github.com/ectrc/snow/storage"
	"github.com/google/uuid"
)

type Person struct {
	ID string
	DisplayName string
	AccessKey string
	AthenaProfile *Profile
	CommonCoreProfile *Profile
	CommonPublicProfile *Profile
	Profile0 *Profile
}

type Option struct {
	Key string
	Value string
}

func NewPerson() *Person {
	return &Person{
		ID: uuid.New().String(),
		DisplayName: uuid.New().String(),
		AccessKey: "",
		AthenaProfile: NewProfile("athena"),
		CommonCoreProfile: NewProfile("common_core"),
		CommonPublicProfile: NewProfile("common_public"),
		Profile0: NewProfile("profile0"),
	}
}

func Find(personId string) *Person {
	if cache == nil {
		cache = NewPersonsCacheMutex()
	}

	cachedPerson := cache.GetPerson(personId)
	if cachedPerson != nil {
		return cachedPerson
	}

	person := storage.Repo.GetPersonFromDB(personId)
	if person == nil {
		return nil
	}

	return findHelper(person)
}

func FindByDisplay(displayName string) *Person {
	if cache == nil {
		cache = NewPersonsCacheMutex()
	}

	person := storage.Repo.GetPersonByDisplayFromDB(displayName)
	if person == nil {
		return nil
	}

	cachedPerson := cache.GetPerson(person.ID)
	if cachedPerson != nil {
		return cachedPerson
	}

	return findHelper(person)
}

func findHelper(databasePerson *storage.DB_Person) *Person {
	athenaProfile := NewProfile("athena")
	commonCoreProfile := NewProfile("common_core")
	commonPublicProfile := NewProfile("common_public")
	profile0 := NewProfile("profile0")

	for _, profile := range databasePerson.Profiles {
		if profile.Type == "athena" {
			athenaProfile.ID = profile.ID
			athenaProfile = FromDatabaseProfile(&profile)
		}

		if profile.Type == "common_core" {
			commonCoreProfile.ID = profile.ID
			commonCoreProfile = FromDatabaseProfile(&profile)
		}

		if profile.Type == "common_public" {
			commonPublicProfile.ID = profile.ID
			commonPublicProfile = FromDatabaseProfile(&profile)
		}

		if profile.Type == "profile0" {
			profile0.ID = profile.ID
			profile0 = FromDatabaseProfile(&profile)
		}
	}

	person := &Person{
		ID: databasePerson.ID,
		DisplayName: databasePerson.DisplayName,
		AccessKey: databasePerson.AccessKey,
		AthenaProfile: athenaProfile,
		CommonCoreProfile: commonCoreProfile,
		CommonPublicProfile: commonPublicProfile,
		Profile0: profile0,
	}

	cache.Store(person.ID, &CacheEntry{
		Entry: person,
		LastAccessed: time.Now(),
	})
	
	return person
}

func AllFromDatabase() []*Person {
	var persons []*Person
	for _, person := range storage.Repo.GetAllPersons() {
		persons = append(persons, Find(person.ID))
	}

	return persons
}

func AllFromCache() []*Person {
	if cache == nil {
		cache = NewPersonsCacheMutex()
	}

	var persons []*Person
	cache.RangeEntry(func(key string, value *CacheEntry) bool {
		persons = append(persons, value.Entry)
		return true
	})

	return persons
}

func (p *Person) GetProfileFromType(profileType string) *Profile {
	switch profileType {
	case "athena":
		return p.AthenaProfile
	case "common_core":
		return p.CommonCoreProfile
	case "common_public":
		return p.CommonPublicProfile
	case "profile0":
		return p.Profile0
	}

	return nil
}

func (p *Person) Save() {
	dbPerson := p.ToDatabase()
	storage.Repo.SavePerson(dbPerson)
}

func (p *Person) ToDatabase() *storage.DB_Person {
	dbPerson := storage.DB_Person{
		ID: p.ID,
		DisplayName: p.DisplayName,
		Profiles: []storage.DB_Profile{},
		AccessKey: p.AccessKey,
	}

	profilesToConvert := map[string]*Profile{
		"common_core": p.CommonCoreProfile,
		"athena": p.AthenaProfile,
		"common_public": p.CommonPublicProfile,
		"profile0": p.Profile0,
	}

	for profileType, profile := range profilesToConvert {
		dbProfile := storage.DB_Profile{
			ID: profile.ID,
			PersonID: p.ID,
			Type: profileType,
			Items: []storage.DB_Item{},
			Gifts: []storage.DB_Gift{},
			Quests: []storage.DB_Quest{},
			Attributes: []storage.DB_PAttribute{},
			Revision: profile.Revision,
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

		profile.Attributes.RangeAttributes(func(key string, value *Attribute) bool {
			dbProfile.Attributes = append(dbProfile.Attributes, *value.ToDatabase(p.ID))
			return true
		})

		dbPerson.Profiles = append(dbPerson.Profiles, dbProfile)
	}

	return &dbPerson
}

func (p *Person) AddAttribute(value *Attribute) {
	p.AthenaProfile.Attributes.AddAttribute(value)
}

func (p *Person) GetAttribute(key string) *Attribute {
	attribute := p.AthenaProfile.Attributes.GetAttribute(key)
	return attribute
}

func (p *Person) RemoveAttribute(key string) {
	p.AthenaProfile.Attributes.DeleteAttribute(key)
}

func (p *Person) Snapshot() *PersonSnapshot {
	return &PersonSnapshot{
		ID: p.ID,
		DisplayName: p.DisplayName,
		AthenaProfile: *p.AthenaProfile.Snapshot(),
		CommonCoreProfile: *p.CommonCoreProfile.Snapshot(),
	}
} 