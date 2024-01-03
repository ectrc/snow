package person

import (
	"reflect"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/storage"
	"github.com/google/uuid"
)

type Attribute struct {
	ID string
	ProfileID string
	Key string
	ValueJSON string
	Type  string
}

func NewAttribute(key string, value interface{}) *Attribute {
	return &Attribute{
		ID: uuid.New().String(),
		Key: key,
		ValueJSON: aid.JSONStringify(value),
		Type: reflect.TypeOf(value).String(),
	}
}

func FromDatabaseAttribute(db *storage.DB_PAttribute) *Attribute {
	return &Attribute{
		ID: db.ID,
		ProfileID: db.ProfileID,
		Key: db.Key,
		ValueJSON: db.ValueJSON,
		Type: db.Type,
	}
}

func (a *Attribute) ToDatabase(profileId string) *storage.DB_PAttribute {
	return &storage.DB_PAttribute{
		ID: a.ID,
		ProfileID: profileId,
		Key: a.Key,
		ValueJSON: a.ValueJSON,
		Type: a.Type,
	}
}

func (a *Attribute) Delete() {
	storage.Repo.DeleteAttribute(a.ID)
}

func (a *Attribute) Save() {
	if a.ProfileID == "" {
		return
	}
	storage.Repo.SaveAttribute(a.ToDatabase(a.ProfileID))
}