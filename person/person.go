package person

import (
	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/storage"
	"github.com/google/uuid"
)

type Person struct {
	ID string
	DisplayName string
	Permissions []string
	IsBanned bool
	AthenaProfile *Profile
	CommonCoreProfile *Profile
	CommonPublicProfile *Profile
	Profile0Profile *Profile
	CollectionsProfile *Profile
	CreativeProfile *Profile
	Discord *storage.DB_DiscordPerson
	Relationships aid.GenericSyncMap[Relationship]
}

func NewPerson() *Person {
	return &Person{
		ID: uuid.New().String(),
		DisplayName: uuid.New().String(),
		Permissions: []string{},
		IsBanned: false,
		AthenaProfile: NewProfile("athena"),
		CommonCoreProfile: NewProfile("common_core"),
		CommonPublicProfile: NewProfile("common_public"),
		Profile0Profile: NewProfile("profile0"),
		CollectionsProfile: NewProfile("collections"),
		CreativeProfile: NewProfile("creative"),
	}
}

func NewPersonWithCustomID(id string) *Person {
	return &Person{
		ID: id,
		DisplayName: uuid.New().String(),
		Permissions: []string{},
		IsBanned: false,
		AthenaProfile: NewProfile("athena"),
		CommonCoreProfile: NewProfile("common_core"),
		CommonPublicProfile: NewProfile("common_public"),
		Profile0Profile: NewProfile("profile0"),
		CollectionsProfile: NewProfile("collections"),
		CreativeProfile: NewProfile("creative"),
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

	return findHelper(person, false, true)
}

func FindShallow(personId string) *Person {
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

	return findHelper(person, true, true)
}

func FindByDisplay(displayName string) *Person {
	if cache == nil {
		cache = NewPersonsCacheMutex()
	}

	cachedPerson := cache.GetPersonByDisplay(displayName)
	if cachedPerson != nil {
		return cachedPerson
	}

	person := storage.Repo.GetPersonByDisplayFromDB(displayName)
	if person == nil {
		return nil
	}

	return findHelper(person, false, true)
}

func FindByDisplayShallow(displayName string) *Person {
	if cache == nil {
		cache = NewPersonsCacheMutex()
	}

	cachedPerson := cache.GetPersonByDisplay(displayName)
	if cachedPerson != nil {
		return cachedPerson
	}

	person := storage.Repo.GetPersonByDisplayFromDB(displayName)
	if person == nil {
		return nil
	}

	return findHelper(person, true, true)
}

func FindByDiscord(discordId string) *Person {
	if cache == nil {
		cache = NewPersonsCacheMutex()
	}

	cachedPerson := cache.GetPersonByDiscordID(discordId)
	if cachedPerson != nil {
		return cachedPerson
	}

	person := storage.Repo.GetPersonByDiscordIDFromDB(discordId)
	if person == nil {
		return nil
	}

	return findHelper(person, false, true)
}

func findHelper(databasePerson *storage.DB_Person, shallow bool, save bool) *Person {
	athenaProfile := NewProfile("athena")
	commonCoreProfile := NewProfile("common_core")
	commonPublicProfile := NewProfile("common_public")
	profile0 := NewProfile("profile0")
	collectionsProfile := NewProfile("collections")
	creativeProfile := NewProfile("creative")

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

		if profile.Type == "collections" {
			collectionsProfile.ID = profile.ID
			collectionsProfile = FromDatabaseProfile(&profile)
		}

		if profile.Type == "creative" {
			creativeProfile.ID = profile.ID
			creativeProfile = FromDatabaseProfile(&profile)
		}
	}

	person := &Person{
		ID: databasePerson.ID,
		DisplayName: databasePerson.DisplayName,
		Permissions: databasePerson.Permissions,
		IsBanned: databasePerson.IsBanned,
		AthenaProfile: athenaProfile,
		CommonCoreProfile: commonCoreProfile,
		CommonPublicProfile: commonPublicProfile,
		Profile0Profile: profile0,
		CollectionsProfile: collectionsProfile,
		CreativeProfile: creativeProfile,
		Discord: &databasePerson.Discord,
		Relationships: aid.GenericSyncMap[Relationship]{},
	}

	if !shallow {
		person.LoadRelationships()
	}

	if save {
		cache.SavePerson(person)
	}
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
		return p.Profile0Profile
	case "collections":
		return p.CollectionsProfile
	case "creative":
		return p.CreativeProfile
	}

	return nil
}

func (p *Person) Save() {
	dbPerson := p.ToDatabase()
	storage.Repo.SavePerson(dbPerson)
}

func (p *Person) Ban() {
	p.IsBanned = true
	p.Save()
}

func (p *Person) Unban() {
	p.IsBanned = false
	p.Save()
}

func (p *Person) AddPermission(permission string) {
	p.Permissions = append(p.Permissions, permission)
	p.Save()
}

func (p *Person) RemovePermission(permission string) {
	for i, perm := range p.Permissions {
		if perm == permission {
			p.Permissions = append(p.Permissions[:i], p.Permissions[i+1:]...)
			break
		}
	}
	p.Save()
}

func (p *Person) HasPermission(permission Permission) bool {
	for _, perm := range p.Permissions {
		if perm == string(PermissionAll) {
			return true
		}

		if perm == string(permission) {
			return true
		}
	}

	return false
}

func (p *Person) ToDatabase() *storage.DB_Person {
	dbPerson := storage.DB_Person{
		ID: p.ID,
		DisplayName: p.DisplayName,
		Permissions: p.Permissions,
		IsBanned: p.IsBanned,
		Profiles: []storage.DB_Profile{},
		Stats: []storage.DB_SeasonStat{},
		Discord: storage.DB_DiscordPerson{},
	}

	if p.Discord != nil {
		dbPerson.Discord = *p.Discord
	}

	profilesToConvert := map[string]*Profile{
		"common_core": p.CommonCoreProfile,
		"athena": p.AthenaProfile,
		"common_public": p.CommonPublicProfile,
		"profile0": p.Profile0Profile,
		"collections": p.CollectionsProfile,
		"creative": p.CreativeProfile,
	}

	for profileType, profile := range profilesToConvert {
		dbProfile := storage.DB_Profile{
			ID: profile.ID,
			PersonID: p.ID,
			Type: profileType,
			Items: []storage.DB_Item{},
			Gifts: []storage.DB_Gift{},
			Quests: []storage.DB_Quest{},
			Loadouts: []storage.DB_Loadout{},
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

		profile.Loadouts.RangeLoadouts(func(id string, loadout *Loadout) bool {
			dbProfile.Loadouts = append(dbProfile.Loadouts, *loadout.ToDatabase(p.ID))
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
		Permissions: p.Permissions,
		IsBanned: p.IsBanned,
		AthenaProfile: *p.AthenaProfile.Snapshot(),
		CommonCoreProfile: *p.CommonCoreProfile.Snapshot(),
		CommonPublicProfile: *p.CommonPublicProfile.Snapshot(),
		Profile0Profile: *p.Profile0Profile.Snapshot(),
		CollectionsProfile: *p.CollectionsProfile.Snapshot(),
		CreativeProfile: *p.CreativeProfile.Snapshot(),
		Discord: *p.Discord,
	}
} 

func (p *Person) Delete() {
	storage.Repo.DeletePerson(p.ID)
}