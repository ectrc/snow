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