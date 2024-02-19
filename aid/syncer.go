package aid

import (
	"sync"
)

type GenericSyncMap[T any] struct {
	m sync.Map
}

func (s *GenericSyncMap[T]) Set(key string, value *T) {
	s.m.Store(key, value)
}

func (s *GenericSyncMap[T]) Get(key string) (*T, bool) {
	v, ok := s.m.Load(key)
	if !ok {
		return nil, false
	}

	return v.(*T), true
}

func (s *GenericSyncMap[T]) MustGet(key string) *T {
	v, ok := s.m.Load(key)
	if !ok {
		return nil
	}

	return v.(*T)
}

func (s *GenericSyncMap[T]) Delete(key string) {
	s.m.Delete(key)
}

func (s *GenericSyncMap[T]) Range(f func(key string, value *T) bool) {
	s.m.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(*T))
	})
}

func (s *GenericSyncMap[T]) Has(key string) bool {
	_, ok := s.Get(key)
	return ok
}

// change the key of a value in the map

func (s *GenericSyncMap[T]) ChangeKey(oldKey, newKey string) {
	v, ok := s.Get(oldKey)
	if !ok {
		return
	}

	s.Set(newKey, v)
	s.Delete(oldKey)
}

func (s *GenericSyncMap[T]) Len() int {
	count := 0
	s.m.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}

func (s *GenericSyncMap[T]) Snapshot() *map[string]*T {
	m := make(map[string]*T)
	s.Range(func(key string, value *T) bool {
		m[key] = value
		return true
	})
	return &m
}