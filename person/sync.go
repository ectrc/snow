package person

import (
	"sync"

	"github.com/ectrc/snow/storage"
)

type ItemMutex struct {
	sync.Map
	ProfileType string
	ProfileID	 string
}

func NewItemMutex(profile *storage.DB_Profile) *ItemMutex {
	return &ItemMutex{
		ProfileType: profile.Type,
		ProfileID:	 profile.ID,
	}
}

func (m *ItemMutex) AddItem(item *Item) {
	item.ProfileType = m.ProfileType
	item.ProfileID = m.ProfileID
	m.Store(item.ID, item)
	storage.Repo.SaveItem(item.ToDatabase(m.ProfileID))
}

func (m *ItemMutex) DeleteItem(id string) {
	item := m.GetItem(id)
	if item == nil {
		return
	}

	item.Delete()
	m.Delete(id)
	storage.Repo.DeleteItem(id)
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
	ProfileID	 string
}

func NewGiftMutex(profile *storage.DB_Profile) *GiftMutex {
	return &GiftMutex{
		ProfileType: profile.Type,
		ProfileID:	 profile.ID,
	}
}

func (m *GiftMutex) AddGift(gift *Gift) {
	gift.ProfileID = m.ProfileID
	m.Store(gift.ID, gift)
	storage.Repo.SaveGift(gift.ToDatabase(m.ProfileID))
}

func (m *GiftMutex) DeleteGift(id string) {
	m.Delete(id)
	storage.Repo.DeleteGift(id)
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
	ProfileType string
	ProfileID	 string
}

func NewQuestMutex(profile *storage.DB_Profile) *QuestMutex {
	return &QuestMutex{
		ProfileType: profile.Type,
		ProfileID:	 profile.ID,
	}
}

func (m *QuestMutex) AddQuest(quest *Quest) {
	quest.ProfileID = m.ProfileID
	m.Store(quest.ID, quest)
	storage.Repo.SaveQuest(quest.ToDatabase(m.ProfileID))
}

func (m *QuestMutex) DeleteQuest(id string) {
	m.Delete(id)
	storage.Repo.DeleteQuest(id)
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
	ProfileType string
	ProfileID	 string
}

func NewAttributeMutex(profile *storage.DB_Profile) *AttributeMutex {
	return &AttributeMutex{
		ProfileID:	 profile.ID,
	}
}

func (m *AttributeMutex) AddAttribute(attribute *Attribute) {
	attribute.ProfileID = m.ProfileID
	m.Store(attribute.ID, attribute)
	storage.Repo.SaveAttribute(attribute.ToDatabase(m.ProfileID))
}

func (m *AttributeMutex) DeleteAttribute(id string) {
	m.Delete(id)
	storage.Repo.DeleteAttribute(id)
}

func (m *AttributeMutex) GetAttribute(id string) *Attribute {
	value, ok := m.Load(id)
	if !ok {
		return nil
	}

	return value.(*Attribute)
}

func (m *AttributeMutex) GetAttributeByKey(key string) *Attribute {
	var found *Attribute

	m.RangeAttributes(func(id string, attribute *Attribute) bool {
		if attribute.Key == key {
			found = attribute
			return false
		}

		return true
	})

	return found
}

func (m *AttributeMutex) RangeAttributes(f func(id string, attribute *Attribute) bool) {
	m.Range(func(key, value interface{}) bool {
		return f(key.(string), value.(*Attribute))
	})
}

func (m *AttributeMutex) Count() int {
	count := 0
	m.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}