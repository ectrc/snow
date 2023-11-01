package person

import (
	"encoding/json"
	"reflect"

	"github.com/ectrc/snow/storage"
	"github.com/google/uuid"
)

type Attribute struct {
	ID 	  string
	Key   string
	Value interface{}
	Type  string
}

func NewAttribute(key string, value interface{}) *Attribute {
	return &Attribute{
		ID:    uuid.New().String(),
		Key:   key,
		Value: value,
		Type:  reflect.TypeOf(value).String(),
	}
}

func FromDatabaseAttribute(db *storage.DB_PAttribute) *Attribute {
	var value interface{}
	err := json.Unmarshal([]byte(db.ValueJSON), &value)
	if err != nil {
		return nil
	}

	return &Attribute{
		ID:    db.ID,
		Key:   db.Key,
		Value: value,
		Type:  db.Type,
	}
}

func (a *Attribute) ToDatabase(profileId string) *storage.DB_PAttribute {
	value, err := json.Marshal(a.Value)
	if err != nil {
		return nil
	}

	return &storage.DB_PAttribute{
		ID:        a.ID,
		ProfileID: profileId,
		Key:       a.Key,
		ValueJSON: string(value),
		Type:      a.Type,
	}
}

func (a *Attribute) Delete() {
	storage.Repo.DeleteAttribute(a.ID)
}