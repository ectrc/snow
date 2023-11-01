package person

import (
	"sync"
)

type ItemMutex struct {
	sync.Map
	ProfileType string
}

func NewItemMutex(profile string) *ItemMutex {
	return &ItemMutex{
		ProfileType: profile,
	}
}

func (m *ItemMutex) AddItem(item *Item) {
	item.ProfileType = m.ProfileType
	m.Store(item.ID, item)
	// storage.Repo.SaveItem(item)
}

func (m *ItemMutex) DeleteItem(id string) {
	item := m.GetItem(id)
	if item == nil {
		return
	}

	item.Delete()
	m.Delete(id)
	// storage.Repo.DeleteItem(id)
}

func (m *ItemMutex) GetItem(id string) *Item {
	item, ok := m.Load(id)
	if !ok {
		return nil
	}

	return item.(*Item)
}

func (m *ItemMutex) GetItemByTemplateID(templateID string) *Item {
	var item *Item

	m.Range(func(key, value interface{}) bool {
		if value.(*Item).TemplateID == templateID {
			item = value.(*Item)
			return false
		}

		return true
	})

	return item
}

func (m *ItemMutex) RangeItems(f func(key string, value *Item) bool) {
	m.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(*Item))
	})
}

func (m *ItemMutex) Count() int {
	count := 0
	m.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

type GiftMutex struct {
	sync.Map
	ProfileType string
}

func NewGiftMutex() *GiftMutex {
	return &GiftMutex{}
}

func (m *GiftMutex) AddGift(gift *Gift) {
	m.Store(gift.ID, gift)
	// storage.Repo.SaveGift(gift)
}

func (m *GiftMutex) DeleteGift(id string) {
	m.Delete(id)
	// storage.Repo.DeleteGift(id)
}

func (m *GiftMutex) GetGift(id string) *Gift {
	gift, ok := m.Load(id)
	if !ok {
		return nil
	}

	return gift.(*Gift)
}

func (m *GiftMutex) RangeGifts(f func(key string, value *Gift) bool) {
	m.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(*Gift))
	})
}

func (m *GiftMutex) Count() int {
	count := 0
	m.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

type QuestMutex struct {
	sync.Map
}

func NewQuestMutex() *QuestMutex {
	return &QuestMutex{}
}

func (m *QuestMutex) AddQuest(quest *Quest) {
	m.Store(quest.ID, quest)
	// storage.Repo.SaveQuest(quest)
}

func (m *QuestMutex) DeleteQuest(id string) {
	m.Delete(id)
	// storage.Repo.DeleteQuest(id)
}

func (m *QuestMutex) GetQuest(id string) *Quest {
	quest, ok := m.Load(id)
	if !ok {
		return nil
	}

	return quest.(*Quest)
}

func (m *QuestMutex) RangeQuests(f func(key string, value *Quest) bool) {
	m.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(*Quest))
	})
}

func (m *QuestMutex) Count() int {
	count := 0
	m.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

type AttributeMutex struct {
	sync.Map
}

func NewAttributeMutex() *AttributeMutex {
	return &AttributeMutex{}
}

func (m *AttributeMutex) AddAttribute(attribute *Attribute) {
	m.Store(attribute.Key, attribute)
	// storage.Repo.SaveAttribute(key, value)
}

func (m *AttributeMutex) DeleteAttribute(key string) {
	m.Delete(key)
	// storage.Repo.DeleteAttribute(key)
}

func (m *AttributeMutex) GetAttribute(key string) *Attribute {
	value, ok := m.Load(key)
	if !ok {
		return nil
	}

	return value.(*Attribute)
}

func (m *AttributeMutex) RangeAttributes(f func(key string, value *Attribute) bool) {
	m.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(*Attribute))
	})
}