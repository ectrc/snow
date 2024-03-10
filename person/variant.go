package person

import (
	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/storage"
	"github.com/google/uuid"
)

type VariantToken struct {
	ID string
	ProfileID string
	TemplateID string
	Name string
	AutoEquipOnGrant bool
	CreateGiftboxOnGrant  bool
	MarkItemUnseenOnGrant bool
	VariantGrants []*VariantTokenGrant
}

func NewVariantToken(profileID, templateID string) *VariantToken {
	return &VariantToken{
		ID: uuid.New().String(),
		ProfileID:  profileID,
		TemplateID: templateID,
	}
}

func (v *VariantToken) AddVariantGrant(channel, value string) {
	vtGrant := &VariantTokenGrant{
		ID: uuid.New().String(),
		VariantTokenID: v.ID,
		Channel: channel,
		Value: value,
	}

	v.VariantGrants = append(v.VariantGrants, vtGrant)
}

func FromDatabaseVariantToken(token *storage.DB_VariantToken) *VariantToken {
	variantGrants := []*VariantTokenGrant{}

	for _, grant := range token.VariantGrants {
		variantGrants = append(variantGrants, FromDatabaseVariantTokenGrant(&grant))
	}

	return &VariantToken{
		ID: token.ID,
		ProfileID: token.ProfileID,
		TemplateID: token.TemplateID,
		Name: token.Name,
		AutoEquipOnGrant: token.AutoEquipOnGrant,
		CreateGiftboxOnGrant: token.CreateGiftboxOnGrant,
		MarkItemUnseenOnGrant: token.MarkItemUnseenOnGrant,
		VariantGrants: variantGrants,
	}
}

func (v *VariantToken) GenerateFortniteVariantTokenEntry() aid.JSON {
	return aid.JSON{
		"templateId": v.TemplateID,
		"attributes": aid.JSON{
			"auto_equip_variant": v.AutoEquipOnGrant,
			"create_giftbox": v.CreateGiftboxOnGrant,
			"mark_item_unseen": v.MarkItemUnseenOnGrant,
			"variant_name": v.Name,
		},
		"quantity": 1,
	}
}

func (v *VariantToken) ToDatabase(profileID string) *storage.DB_VariantToken {
	variantGrants := []storage.DB_VariantTokenGrant{}

	for _, grant := range v.VariantGrants {
		variantGrants = append(variantGrants, *grant.ToDatabase())
	}

	return &storage.DB_VariantToken{
		ID: v.ID,
		ProfileID: profileID,
		TemplateID: v.TemplateID,
		Name: v.Name,
		AutoEquipOnGrant: v.AutoEquipOnGrant,
		CreateGiftboxOnGrant: v.CreateGiftboxOnGrant,
		MarkItemUnseenOnGrant: v.MarkItemUnseenOnGrant,
		VariantGrants: variantGrants,
	}
}

func (v *VariantToken) Save() {
	storage.Repo.SaveVariantToken(v.ToDatabase(v.ProfileID))
}

func (v *VariantToken) Delete() {
	storage.Repo.DeleteVariantToken(v.ID)
}

type VariantTokenGrant struct {
	ID string
	VariantTokenID string
	Channel string
	Value string
}

func NewVariantTokenGrant(vtID, channel, value string) *VariantTokenGrant {
	return &VariantTokenGrant{
		ID: uuid.New().String(),
		VariantTokenID: vtID,
		Channel: channel,
		Value: value,
	}
}

func FromDatabaseVariantTokenGrant(grant *storage.DB_VariantTokenGrant) *VariantTokenGrant {
	return &VariantTokenGrant{
		ID: grant.ID,
		VariantTokenID: grant.VariantTokenID,
		Channel: grant.Channel,
		Value: grant.Value,
	}
}

func (v *VariantTokenGrant) ToDatabase() *storage.DB_VariantTokenGrant {
	return &storage.DB_VariantTokenGrant{
		ID: v.ID,
		VariantTokenID: v.VariantTokenID,
		Channel: v.Channel,
		Value: v.Value,
	}
}