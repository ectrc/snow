package storage

import (
	"sync"
)

type personsMutex struct {
	sync.Map
}

func newPersonsMutex() *personsMutex {
	return &personsMutex{}
}

func (m *personsMutex) GetPerson(id string) *DB_Person {
	p, ok := m.Load(id)
	if !ok {
		return nil
	}

	return p.(*DB_Person)
}

func (m *personsMutex) SavePerson(person *DB_Person) {
	m.Store(person.ID, person)
}

type MemoryStorage struct {
	Persons *personsMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		Persons: newPersonsMutex(),
	}
}

func (s *MemoryStorage) Migrate(table interface{}, tableName string) {} // not needed for memory storage as there is no db

func (s *MemoryStorage) GetPerson(id string) *DB_Person {
	return s.Persons.GetPerson(id)
}

func (s *MemoryStorage) SavePerson(person *DB_Person) {
	s.Persons.SavePerson(person)
}