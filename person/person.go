package person

import (
	"time"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/storage"
	"github.com/google/uuid"
)

type Person struct {
	ID string
	DisplayName string
	RefundTickets int
	Permissions Permission
	AthenaProfile *Profile
	CommonCoreProfile *Profile
	CommonPublicProfile *Profile
	Profile0Profile *Profile
	CollectionsProfile *Profile
	CreativeProfile *Profile
	Discord *storage.DB_DiscordPerson
	BanHistory aid.GenericSyncMap[storage.DB_BanStatus]
	Relationships aid.GenericSyncMap[Relationship]
}

func NewPerson() *Person {
	return &Person{
		ID: uuid.New().String(),
		DisplayName: uuid.New().String(),
		Permissions: 0,
		RefundTickets: 3,
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
		Permissions: 0,
		RefundTickets: 3,
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

	return findHelper(person, true, false)
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

	return findHelper(person, true, false)
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
		Permissions: Permission(databasePerson.Permissions),
		BanHistory: aid.GenericSyncMap[storage.DB_BanStatus]{},
		AthenaProfile: athenaProfile,
		CommonCoreProfile: commonCoreProfile,
		CommonPublicProfile: commonPublicProfile,
		Profile0Profile: profile0,
		CollectionsProfile: collectionsProfile,
		CreativeProfile: creativeProfile,
		Discord: &databasePerson.Discord,
		RefundTickets: databasePerson.RefundTickets,
		Relationships: aid.GenericSyncMap[Relationship]{},
	}

	for _, ban := range databasePerson.BanHistory {
		person.BanHistory.Set(ban.ID, &ban)
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

func (p *Person) SaveShallow() {
	dbPerson := p.ToDatabaseShallow()
	storage.Repo.SavePerson(dbPerson)
}

func (p *Person) AddBan(reason string, issuedBy string, expiry ...string) {
	t := time.Now().AddDate(0, 0, 7)

	if len(expiry) > 0 && expiry[0] != "" {
		parsed, err := aid.ParseDuration(expiry[0])
		if err == nil {
			t = time.Now().Add(parsed)
		}
	}

	ban := &storage.DB_BanStatus{
		ID: uuid.New().String(),
		PersonID: p.ID,
		IssuedBy: issuedBy,
		Reason: reason,
		Expiry: t,
	}

	p.BanHistory.Set(ban.ID, ban)
	storage.Repo.SaveBanStatus(ban)
}

func (p *Person) ClearBans() {
	p.BanHistory.Range(func(key string, ban *storage.DB_BanStatus) bool {
		ban.Expiry = time.Now()
		storage.Repo.SaveBanStatus(ban)
		return true
	})
}

func (p *Person) GetLatestActiveBan() *storage.DB_BanStatus {
	var latestBan *storage.DB_BanStatus
	p.BanHistory.Range(func(key string, ban *storage.DB_BanStatus) bool {
		if latestBan == nil || ban.Expiry.After(latestBan.Expiry) {
			latestBan = ban
		}
		return true
	})

	if latestBan != nil && latestBan.Expiry.Before(time.Now()) {
		return nil
	}

	return latestBan
}

func (p *Person) AddPermission(permission Permission) {
	p.Permissions |= permission
	p.SaveShallow()
}

func (p *Person) RemovePermission(permission Permission) {
	p.Permissions &= ^permission
	p.SaveShallow()
}

func (p *Person) HasPermission(permission Permission) bool {
	// if permission == PermissionAll && permission != PermissionOwner {
	// 	return p.Permissions == PermissionAll
	// }

	return p.Permissions & permission != 0
}

func (p *Person) ToDatabase() *storage.DB_Person {
	dbPerson := storage.DB_Person{
		ID: p.ID,
		DisplayName: p.DisplayName,
		Permissions: int64(p.Permissions),
		BanHistory: []storage.DB_BanStatus{},
		RefundTickets: p.RefundTickets,
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

	p.BanHistory.Range(func(key string, ban *storage.DB_BanStatus) bool {
		dbPerson.BanHistory = append(dbPerson.BanHistory, *ban)
		return true
	})

	for profileType, profile := range profilesToConvert {
		dbProfile := storage.DB_Profile{
			ID: profile.ID,
			PersonID: p.ID,
			Type: profileType,
			Items: []storage.DB_Item{},
			Gifts: []storage.DB_Gift{},
			Quests: []storage.DB_Quest{},
			Loadouts: []storage.DB_Loadout{},
			Attributes: []storage.DB_Attribute{},
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

func (p *Person) ToDatabaseShallow() *storage.DB_Person {
	dbPerson := storage.DB_Person{
		ID: p.ID,
		DisplayName: p.DisplayName,
		Permissions: int64(p.Permissions),
		BanHistory: []storage.DB_BanStatus{},
		RefundTickets: p.RefundTickets,
		Profiles: []storage.DB_Profile{},
		Stats: []storage.DB_SeasonStat{},
		Discord: storage.DB_DiscordPerson{},
	}

	if p.Discord != nil {
		dbPerson.Discord = *p.Discord
	}

	p.BanHistory.Range(func(key string, ban *storage.DB_BanStatus) bool {
		dbPerson.BanHistory = append(dbPerson.BanHistory, *ban)
		return true
	})

	return &dbPerson
}

func (p *Person) Snapshot() *PersonSnapshot {
	snapshot := &PersonSnapshot{
		ID: p.ID,
		DisplayName: p.DisplayName,
		Permissions: int64(p.Permissions),
		AthenaProfile: *p.AthenaProfile.Snapshot(),
		CommonCoreProfile: *p.CommonCoreProfile.Snapshot(),
		CommonPublicProfile: *p.CommonPublicProfile.Snapshot(),
		Profile0Profile: *p.Profile0Profile.Snapshot(),
		CollectionsProfile: *p.CollectionsProfile.Snapshot(),
		CreativeProfile: *p.CreativeProfile.Snapshot(),
		BanHistory: []storage.DB_BanStatus{},
		Discord: *p.Discord,
	}

	p.BanHistory.Range(func(key string, ban *storage.DB_BanStatus) bool {
		snapshot.BanHistory = append(snapshot.BanHistory, *ban)
		return true
	})

	return snapshot
} 

func (p *Person) Delete() {
	storage.Repo.DeletePerson(p.ID)
	cache.DeletePerson(p.ID)
}

func (p *Person) SetPurchaseHistoryAttribute() {
	purchases := []aid.JSON{}

	p.AthenaProfile.Purchases.RangePurchases(func(key string, value *Purchase) bool {
		purchases = append(purchases, value.GenerateFortnitePurchaseEntry())
		return true
	})

	purchaseAttribute := p.CommonCoreProfile.Attributes.GetAttributeByKey("mtx_purchase_history")
	purchaseAttribute.ValueJSON = aid.JSONStringify(aid.JSON{
		"refundsUsed": p.AthenaProfile.Purchases.CountRefunded(),
		"refundCredits": p.RefundTickets,
		"purchases": purchases,
	})
	purchaseAttribute.Save()
}