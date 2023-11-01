package person

import (
	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/storage"
	"github.com/google/uuid"
)

type Item struct {
	ID          string
	TemplateID  string
	Quantity    int
	Favorite    bool
	HasSeen     bool
	Variants    []*VariantChannel
	ProfileType string
}

func NewItem(templateID string, quantity int) *Item {
	return &Item{
		ID:         uuid.New().String(),
		TemplateID: templateID,
		Quantity:   quantity,
		Favorite:   false,
		HasSeen:    false,
		Variants:   []*VariantChannel{},
	}
}

func NewItemWithType(templateID string, quantity int, profile string) *Item {
	return &Item{
		ID:         uuid.New().String(),
		TemplateID: templateID,
		Quantity:   quantity,
		Favorite:   false,
		HasSeen:    false,
		Variants:   []*VariantChannel{},
		ProfileType: profile,
	}
}

func FromDatabaseItem(item *storage.DB_Item, profileType *string) *Item {
	variants := []*VariantChannel{}

	for _, variant := range item.Variants {
		variants = append(variants, FromDatabaseVariant(&variant))
	}

	return &Item{
		ID:          item.ID,
		TemplateID:  item.TemplateID,
		Quantity:    item.Quantity,
		Favorite:    item.Favorite,
		HasSeen:     item.HasSeen,
		Variants:    variants,
		ProfileType: *profileType,
	}
}

func FromDatabaseLoot(item *storage.DB_Loot) *Item {
	return &Item{
		ID:	item.ID,
		TemplateID: item.TemplateID,
		Quantity: item.Quantity,
		Favorite: false,
		HasSeen: false,
		Variants: []*VariantChannel{},
		ProfileType: item.ProfileType,
	}
}

func (i *Item) GenerateFortniteItemEntry() aid.JSON {
	varaints := []aid.JSON{}

	for _, variant := range i.Variants {
		varaints = append(varaints, aid.JSON{
			"channel": variant.Channel,
			"owned": variant.Owned,
			"active": variant.Active,
		})
	}

	return aid.JSON{
		"templateId": i.TemplateID,
		"attributes": aid.JSON{
			"variants": varaints,
			"favorite": i.Favorite,
			"item_seen": i.HasSeen,
		},
		"quantity": i.Quantity,
	}
}

func (i *Item) GetAttribute(attribute string) interface{} {
	switch attribute {
	case "Favorite":
		return i.Favorite
	case "HasSeen":
		return i.HasSeen
	case "Variants":
		return i.Variants
	}

	return nil
}

func (i *Item) Delete() {
	//storage.Repo.DeleteItem(i.ID)
	i.Quantity = 0
}

func (i *Item) NewChannel(channel string, owned []string, active string) *VariantChannel {
	return &VariantChannel{
		ItemID:  i.ID,
		Channel: channel,
		Owned:   owned,
		Active:  active,
	}
}

func (i *Item) AddChannel(channel *VariantChannel) {
	i.Variants = append(i.Variants, channel)
	//storage.Repo.SaveItemVariant(i.ID, channel)
}

func (i *Item) RemoveChannel(channel *VariantChannel) {
	for index, c := range i.Variants {
		if c.Channel == channel.Channel {
			i.Variants = append(i.Variants[:index], i.Variants[index+1:]...)
		}
	}
	//storage.Repo.DeleteItemVariant(i.ID, channel)
}

func (i *Item) GetChannel(channel string) *VariantChannel {
	for _, c := range i.Variants {
		if c.Channel == channel {
			return c
		}
	}

	return nil
}

func (i *Item) FillChannels(channels []*VariantChannel) {
	i.Variants = []*VariantChannel{}
	for _, channel := range channels {
		i.AddChannel(channel)
	}
}

func (i *Item) ToDatabase(profileId string) *storage.DB_Item {
	variants := []storage.DB_VariantChannel{}

	for _, variant := range i.Variants {
		variants = append(variants, *variant.ToDatabase())
	}

	return &storage.DB_Item{
		ProfileID:  profileId,
		ID:         i.ID,
		TemplateID: i.TemplateID,
		Quantity:   i.Quantity,
		Favorite:   i.Favorite,
		HasSeen:    i.HasSeen,
		Variants:   variants,
	}
}

func (i *Item) Save() {
	//storage.Repo.SaveItem(i.ToDatabase())
}

func (i *Item) ToLootDatabase(giftId string) *storage.DB_Loot {
	return &storage.DB_Loot{
		GiftID:			 giftId,
		ProfileType: i.ProfileType,
		ID:          i.ID,
		TemplateID:  i.TemplateID,
		Quantity:    i.Quantity,
	}
}

func (i *Item) SaveLoot(giftId string) {
	//storage.Repo.SaveLoot(i.ToLootDatabase(giftId))
}

func (i *Item) Snapshot() ItemSnapshot {
	variants := []VariantChannel{}

	for _, variant := range i.Variants {
		variants = append(variants, *variant)
	}

	return ItemSnapshot{
		ID:          i.ID,
		TemplateID:  i.TemplateID,
		Quantity:    i.Quantity,
		Favorite:    i.Favorite,
		HasSeen:     i.HasSeen,
		Variants:    variants,
		ProfileType: i.ProfileType,
	}
}

type VariantChannel struct {
	ItemID	string
	Channel string
	Owned	 	[]string
	Active	string
}

func FromDatabaseVariant(variant *storage.DB_VariantChannel) *VariantChannel {
	return &VariantChannel{
		ItemID:  variant.ItemID,
		Channel: variant.Channel,
		Owned:   variant.Owned,
		Active:  variant.Active,
	}
}

func (v *VariantChannel) ToDatabase() *storage.DB_VariantChannel {
	return &storage.DB_VariantChannel{
		ItemID:  v.ItemID,
		Channel: v.Channel,
		Owned:   v.Owned,
		Active:  v.Active,
	}
}

func (v *VariantChannel) Save() {
	//storage.Repo.SaveItemVariant(v.ToDatabase())
}