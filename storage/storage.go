package storage

var (
	Repo *Repository
)

type Storage interface {
	Migrate(table interface{}, tableName string)

	GetAllPersons() []*DB_Person
	GetPersonsCount() int
	TotalVBucks() int

	GetPerson(personId string) *DB_Person
	GetPersonByDisplay(displayName string) *DB_Person
	GetPersonsByPartialDisplay(displayName string) []*DB_Person
	GetPersonByDiscordID(discordId string) *DB_Person
	SavePerson(person *DB_Person)
	DeletePerson(personId string)

	GetIncomingRelationships(personId string) []*DB_Relationship
	GetOutgoingRelationships(personId string) []*DB_Relationship
	SaveRelationship(relationship *DB_Relationship)
	DeleteRelationship(relationship *DB_Relationship)

	SaveProfile(profile *DB_Profile)
	DeleteProfile(profileId string)

	SaveItem(item *DB_Item)
	BulkCreateItems(items *[]DB_Item)
	DeleteItem(itemId string)

	SaveVariant(variant *DB_VariantChannel)
	BulkCreateVariants(variants *[]DB_VariantChannel)
	DeleteVariant(variantId string)

	SaveQuest(quest *DB_Quest)
	DeleteQuest(questId string)

	SaveLoot(loot *DB_GiftLoot)
	DeleteLoot(lootId string)

	SaveGift(gift *DB_Gift)
	DeleteGift(giftId string)

	SaveAttribute(attribute *DB_Attribute)
	DeleteAttribute(attributeId string)

	SaveLoadout(loadout *DB_Loadout)
	DeleteLoadout(loadoutId string)

	SavePurchase(purchase *DB_Purchase)
	DeletePurchase(purchaseId string)

	SaveDiscordPerson(person *DB_DiscordPerson)
	DeleteDiscordPerson(personId string)

	SaveBanStatus(ban *DB_BanStatus)
	DeleteBanStatus(banId string)
}

type Repository struct {
	Storage Storage
	Amazon  *AmazonClient
}

func NewStorage(s Storage) *Repository {
	return &Repository{
		Storage: s,
	}
}

func (r *Repository) GetPersonFromDB(personId string) *DB_Person {
	storagePerson := r.Storage.GetPerson(personId)
	if storagePerson != nil {
		return storagePerson
	}

	return nil
}

func (r *Repository) GetPersonByDisplayFromDB(displayName string) *DB_Person {
	storagePerson := r.Storage.GetPersonByDisplay(displayName)
	if storagePerson != nil {
		return storagePerson
	}

	return nil
}

func (r *Repository) GetPersonsByPartialDisplayFromDB(displayName string) []*DB_Person {
	storagePerson := r.Storage.GetPersonsByPartialDisplay(displayName)
	if storagePerson != nil {
		return storagePerson
	}

	return nil
}

func (r *Repository) GetPersonByDiscordIDFromDB(discordId string) *DB_Person {
	storagePerson := r.Storage.GetPersonByDiscordID(discordId)
	if storagePerson != nil {
		return storagePerson
	}

	return nil
}

func (r *Repository) TotalVBucks() int {
	return r.Storage.TotalVBucks()
}

func (r *Repository) GetAllPersons() []*DB_Person {
	return r.Storage.GetAllPersons()
}

func (r *Repository) GetPersonsCount() int {
	return r.Storage.GetPersonsCount()
}

func (r *Repository) SavePerson(person *DB_Person) {
	r.Storage.SavePerson(person)
}

func (r *Repository) DeletePerson(personId string) {
	r.Storage.DeletePerson(personId)
}

func (r *Repository) SaveProfile(profile *DB_Profile) {
	r.Storage.SaveProfile(profile)
}

func (r *Repository) DeleteProfile(profileId string) {
	r.Storage.DeleteProfile(profileId)
}

func (r *Repository) GetIncomingRelationships(personId string) []*DB_Relationship {
	return r.Storage.GetIncomingRelationships(personId)
}

func (r *Repository) GetOutgoingRelationships(personId string) []*DB_Relationship {
	return r.Storage.GetOutgoingRelationships(personId)
}

func (r *Repository) SaveRelationship(relationship *DB_Relationship) {
	r.Storage.SaveRelationship(relationship)
}

func (r *Repository) DeleteRelationship(relationship *DB_Relationship) {
	r.Storage.DeleteRelationship(relationship)
}

func (r *Repository) SaveItem(item *DB_Item) {
	r.Storage.SaveItem(item)
}

func (r *Repository) BulkCreateItems(items *[]DB_Item) {
	r.Storage.BulkCreateItems(items)
}

func (r *Repository) DeleteItem(itemId string) {
	r.Storage.DeleteItem(itemId)
}

func (r *Repository) BulkCreateVariants(variants *[]DB_VariantChannel) {
	r.Storage.BulkCreateVariants(variants)
}

func (r *Repository) SaveVariant(variant *DB_VariantChannel) {
	r.Storage.SaveVariant(variant)
}

func (r *Repository) DeleteVariant(variantId string) {
	r.Storage.DeleteVariant(variantId)
}

func (r *Repository) SaveQuest(quest *DB_Quest) {
	r.Storage.SaveQuest(quest)
}

func (r *Repository) DeleteQuest(questId string) {
	r.Storage.DeleteQuest(questId)
}

func (r *Repository) SaveLoot(loot *DB_GiftLoot) {
	r.Storage.SaveLoot(loot)
}

func (r *Repository) DeleteLoot(lootId string) {
	r.Storage.DeleteLoot(lootId)
}

func (r *Repository) SaveGift(gift *DB_Gift) {
	r.Storage.SaveGift(gift)
}

func (r *Repository) DeleteGift(giftId string) {
	r.Storage.DeleteGift(giftId)
}

func (r *Repository) SaveAttribute(attribute *DB_Attribute) {
	r.Storage.SaveAttribute(attribute)
}

func (r *Repository) DeleteAttribute(attributeId string) {
	r.Storage.DeleteAttribute(attributeId)
}

func (r *Repository) SaveLoadout(loadout *DB_Loadout) {
	r.Storage.SaveLoadout(loadout)
}

func (r *Repository) DeleteLoadout(loadoutId string) {
	r.Storage.DeleteLoadout(loadoutId)
}

func (r *Repository) SavePurchase(purchase *DB_Purchase) {
	r.Storage.SavePurchase(purchase)
}

func (r *Repository) DeletePurchase(purchaseId string) {
	r.Storage.DeletePurchase(purchaseId)
}

func (r *Repository) SaveDiscordPerson(person *DB_DiscordPerson) {
	r.Storage.SaveDiscordPerson(person)
}

func (r *Repository) DeleteDiscordPerson(personId string) {
	r.Storage.DeleteDiscordPerson(personId)
}

func (r *Repository) SaveBanStatus(ban *DB_BanStatus) {
	r.Storage.SaveBanStatus(ban)
}

func (r *Repository) DeleteBanStatus(banId string) {
	r.Storage.DeleteBanStatus(banId)
}