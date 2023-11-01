package person

import (
	"encoding/json"

	"github.com/ectrc/snow/storage"
)

type Attribute struct {
	Key   string
	Value interface{}
	Type  string
}

func NewAttribute(key string, value interface{}, attributeType string) *Attribute {
	return &Attribute{
		Key:   key,
		Value: value,
		Type:  attributeType,
	}
}

func FromDatabaseAttribute(db *storage.DB_PAttribute) *Attribute {
	var value interface{}
	err := json.Unmarshal([]byte(db.ValueJSON), &value)
	if err != nil {
		return nil
	}

	return &Attribute{
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
		ProfileID: profileId,
		Key:      a.Key,
		ValueJSON: string(value),
		Type:     a.Type,
	}
}