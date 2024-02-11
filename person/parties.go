package person

import (
	"sync"
	"time"

	"github.com/ectrc/snow/aid"
	"github.com/google/uuid"
)

type PartyMember struct{
	Person *Person
	ConnectionID string
	Meta map[string]interface{}
	Connections map[string]aid.JSON
	Role string
}

type Party struct{
	ID string
	Members []*PartyMember
	Config map[string]interface{}
	Meta map[string]interface{}
	m sync.Mutex
}

var (
	Parties = aid.GenericSyncMap[Party]{}
)

func NewParty() *Party {
	party := &Party{
		ID: uuid.New().String(),
		Members: []*PartyMember{},
		Config: make(map[string]interface{}),
		Meta: make(map[string]interface{}),
	}

	Parties.Set(party.ID, party)
	return party
}

func (p *Party) AddMember(person *Person) {
	p.m.Lock()
	defer p.m.Unlock()

	partyMember := &PartyMember{
		Person: person,
		Meta: make(map[string]interface{}),
		Connections: make(map[string]aid.JSON),
		Role: "MEMBER",
	}

	p.Members = append(p.Members, partyMember)
	person.Parties.Set(p.ID, p)
	// xmpp to person and rest of party to say new member!
}

func (p *Party) RemoveMember(person *Person) {
	p.m.Lock()
	defer p.m.Unlock()

	for i, member := range p.Members {
		if member.Person == person {
			p.Members = append(p.Members[:i], p.Members[i+1:]...)
			break
		}
	}

	if len(p.Members) == 0 {
		Parties.Delete(p.ID)
	}

	person.Parties.Delete(p.ID)
	// xmpp to person and rest of party to say member left!
}

func (p *Party) UpdateMeta(key string, value interface{}) {
	p.m.Lock()
	defer p.m.Unlock()

	p.Meta[key] = value
	// xmpp to rest of party to say meta updated!
}

func (p *Party) DeleteMeta(key string) {
	p.m.Lock()
	defer p.m.Unlock()

	delete(p.Meta, key)
	// xmpp to rest of party to say meta deleted!
}

func (p *Party) UpdateMemberMeta(person *Person, key string, value interface{}) {
	p.m.Lock()
	defer p.m.Unlock()

	for _, member := range p.Members {
		if member.Person == person {
			member.Meta[key] = value
			// xmpp to person and rest of party to say member meta updated!
			break
		}
	}
}

func (p *Party) DeleteMemberMeta(person *Person, key string) {
	p.m.Lock()
	defer p.m.Unlock()

	for _, member := range p.Members {
		if member.Person == person {
			delete(member.Meta, key)
			// xmpp to person and rest of party to say member meta deleted!
			break
		}
	}
}

func (p *Party) UpdateConfig(key string, value interface{}) {
	p.m.Lock()
	defer p.m.Unlock()

	p.Config[key] = value
	// xmpp to rest of party to say config updated!
}

func (p *Party) DeleteConfig(key string) {
	p.m.Lock()
	defer p.m.Unlock()

	delete(p.Config, key)
	// xmpp to rest of party to say config deleted!
}

func (p *Party) GenerateFortniteParty() aid.JSON {
	p.m.Lock()
	defer p.m.Unlock()

	party := aid.JSON{
		"id": p.ID,
		"members": aid.JSON{},
		"config": p.Config,
		"meta": p.Meta,
		"created_at": "0000-00-00T00:00:00Z",
		"updated_at": time.Now().Format(time.RFC3339),
		"revision": 0,
	}

	for _, member := range p.Members {
		party["members"].(aid.JSON)[member.Person.ID] = aid.JSON{
			"account_id": member.Person.ID,
			"role": member.Role,
			"meta": member.Meta,
			"connections": member.Connections,
		}
	}

	return party
}