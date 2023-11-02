package storage

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Entry interface{}
	LastAccessed time.Time
}

type PersonsCache struct {
	sync.Map
}

func NewPersonsCacheMutex() *PersonsCache {
	return &PersonsCache{}
}

func (m *PersonsCache) CacheKiller() {
	for {
		if Cache.Count() == 0 {
			continue
		}

		Cache.Range(func(key, value interface{}) bool {
			cacheEntry := value.(*CacheEntry)
			
			if time.Since(cacheEntry.LastAccessed) >= 30 * time.Minute {
				Cache.Delete(key)
			}

			return true
		})

		time.Sleep(5000 * time.Minute)
	}
}

func (m *PersonsCache) GetPerson(id string) *DB_Person {
	if p, ok := m.Load(id); ok {
		cacheEntry := p.(*CacheEntry)
		return cacheEntry.Entry.(*DB_Person)
	}

	return nil
}

func (m *PersonsCache) SavePerson(p *DB_Person) {
	m.Store(p.ID, &CacheEntry{
		Entry: p,
		LastAccessed: time.Now(),
	})
}

func (m *PersonsCache) DeletePerson(id string) {
	m.Delete(id)
}

func (m *PersonsCache) RangePersons(f func(key string, value *DB_Person) bool) {
	m.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(*CacheEntry).Entry.(*DB_Person))
	})
}

func (m *PersonsCache) Count() int {
	count := 0
	m.Range(func(key, value interface{}) bool {
		count++
		return true
	})

	return count
}