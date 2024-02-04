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

func FromDatabaseAttribute(db *storage.DB_Attribute) *Attribute {
	return &Attribute{
		ID: db.ID,
		ProfileID: db.ProfileID,
		Key: db.Key,
		ValueJSON: db.ValueJSON,
		Type: db.Type,
	}
}

func (a *Attribute) ToDatabase(profileId string) *storage.DB_Attribute {
	return &storage.DB_Attribute{
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

func AttributeConvertToSlice[T any](attribute *Attribute) []T {
	valuesRaw := aid.JSONParse(attribute.ValueJSON).([]interface{})
	values := make([]T, len(valuesRaw))
	for i, value := range valuesRaw {
		values[i] = value.(T)
	}

	return values
}

func AttributeConvert[T any](attribute *Attribute) T {
	return aid.JSONParse(attribute.ValueJSON).(T)
}