package person

import (
	"github.com/ectrc/snow/aid"
	"github.com/google/uuid"
)

type Party struct{
	ID string
	Members []*Person
}

var (
	Parties = aid.GenericSyncMap[Party]{}
)

func NewParty() *Party {
	party := &Party{
		ID: uuid.New().String(),
	}
	Parties.Set(party.ID, party)
	
	return party
}