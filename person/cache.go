package person

import (
	"fmt"
	"sync"
	"time"
)

var (
	cache *PersonsCache
)

type CacheEntry struct {
	Entry *Person
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
		if m.Count() == 0 {
			continue
		}

		m.Range(func(key, value interface{}) bool {
			cacheEntry := value.(*CacheEntry)
			
			if time.Since(cacheEntry.LastAccessed) >= 30 * time.Minute {
				m.Delete(key)
			}

			return true
		})

		time.Sleep(5000 * time.Minute)
	}
}

func (m *PersonsCache) GetPerson(id string) *Person {
	if p, ok := m.Load(id); ok {
		fmt.Println("Cache hit", id)
		cacheEntry := p.(*CacheEntry)
		return cacheEntry.Entry
	}

	return nil
}

func (m *PersonsCache) GetPersonByDisplay(displayName string) *Person {
	var person *Person
	m.RangeEntry(func(key string, value *CacheEntry) bool {
		if value.Entry.DisplayName == displayName {
			person = value.Entry
			return false
		}

		return true
	})

	return person
}

func (m *PersonsCache) SavePerson(p *Person) {
	m.Store(p.ID, &CacheEntry{
		Entry: p,
		LastAccessed: time.Now(),
	})
}

func (m *PersonsCache) DeletePerson(id string) {
	m.Delete(id)
}

func (m *PersonsCache) RangeEntry(f func(key string, value *CacheEntry) bool) {
	m.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(*CacheEntry))
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