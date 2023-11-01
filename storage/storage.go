package storage

var (
	Repo  *Repository
	Cache *PersonsCache
)

type Storage interface {
	Migrate(table interface{}, tableName string)

	GetPerson(personId string) *DB_Person
	GetAllPersons() []*DB_Person
	SavePerson(person *DB_Person)

	DeleteItem(itemId string)
	DeleteVariant(variantId string)
	DeleteQuest(questId string)
	DeleteLoot(lootId string)
	DeleteGift(giftId string)
	DeleteAttribute(attributeId string)
}

type Repository struct {
	Storage Storage
}

func NewStorage(s Storage) *Repository {
	return &Repository{
		Storage: s,
	}
}

func (r *Repository) GetPerson(personId string) *DB_Person {
	cachePerson := Cache.GetPerson(personId)
	if cachePerson != nil {
		return cachePerson
	}

	storagePerson := r.Storage.GetPerson(personId)
	if storagePerson != nil {
		Cache.SavePerson(storagePerson)
		return storagePerson
	}

	return nil
}

func (r *Repository) GetAllPersons() []*DB_Person {
	return r.Storage.GetAllPersons()
}

func (r *Repository) SavePerson(person *DB_Person) {
	Cache.SavePerson(person)
	r.Storage.SavePerson(person)
}

func (r *Repository) DeleteItem(itemId string) {
	r.Storage.DeleteItem(itemId)
}

func (r *Repository) DeleteVariant(variantId string) {
	r.Storage.DeleteVariant(variantId)
}

func (r *Repository) DeleteQuest(questId string) {
	r.Storage.DeleteQuest(questId)
}

func (r *Repository) DeleteLoot(lootId string) {
	r.Storage.DeleteLoot(lootId)
}

func (r *Repository) DeleteGift(giftId string) {
	r.Storage.DeleteGift(giftId)
}

func (r *Repository) DeleteAttribute(attributeId string) {
	r.Storage.DeleteAttribute(attributeId)
}