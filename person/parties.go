package person

import (
	"sync"
	"time"

	"github.com/ectrc/snow/aid"
	"github.com/google/uuid"
)

type PartyPing struct{
	Person *Person
	Party *Party
}

type PartyMember struct{
	Person *Person
	ConnectionID string
	Meta map[string]interface{}
	Connections map[string]aid.JSON
	Role string
	JoinedAt time.Time
	UpdatedAt time.Time
}

func (pm *PartyMember) GenerateFortnitePartyMember() aid.JSON {
	connections := []aid.JSON{}
	for _, connection := range pm.Connections {
		connections = append(connections, connection)
	}

	return aid.JSON{
		"account_id": pm.Person.ID,
		"role": pm.Role,
		"meta": pm.Meta,
		"joined_at": pm.JoinedAt.Format(time.RFC3339),
		"connections": connections,
		"revision": 0,
	}
}

type Party struct{
	ID string
	Members map[string]*PartyMember
	Config map[string]interface{}
	Meta map[string]interface{}
	m sync.Mutex
	CreatedAt time.Time
}

var (
	Parties = aid.GenericSyncMap[Party]{}
)

func NewParty() *Party {
	party := &Party{
		ID: uuid.New().String(),
		Members: make(map[string]*PartyMember),
		Config: map[string]interface{}{
			"type": "DEFAULT",
			"sub_type": "default",
			"intention_ttl:": 60,
			"invite_ttl:": 60,
		},
		Meta: make(map[string]interface{}),
		CreatedAt: time.Now(),
	}

	Parties.Set(party.ID, party)
	return party
}

func (p *Party) GetMember(person *Person) *PartyMember {
	p.m.Lock()
	defer p.m.Unlock()

	return p.Members[person.ID]
}

func (p *Party) AddMember(person *Person, role string) {
	p.m.Lock()
	defer p.m.Unlock()

	partyMember := &PartyMember{
		Person: person,
		Meta: make(map[string]interface{}),
		Connections: make(map[string]aid.JSON),
		Role: role,
		JoinedAt: time.Now(),
	}

	p.Members[person.ID] = partyMember
	person.Parties.Set(p.ID, p)
	// xmpp to person and rest of party to say new member!
}

func (p *Party) RemoveMember(person *Person) {
	p.m.Lock()
	defer p.m.Unlock()

	delete(p.Members, person.ID)
	if len(p.Members) == 0 {
		Parties.Delete(p.ID)
	}

	person.Parties.Delete(p.ID)
}

func (p *Party) UpdateMeta(m map[string]interface{}) {
	p.m.Lock()
	defer p.m.Unlock()

	for key, value := range m {
		p.Meta[key] = value
	}
}

func (p *Party) DeleteMeta(keys []string) {
	p.m.Lock()
	defer p.m.Unlock()

	for _, key := range keys {
		delete(p.Meta, key)
	}
}

func (p *Party) UpdateMemberMeta(person *Person, m map[string]interface{}) {
	p.m.Lock()
	defer p.m.Unlock()

	member, ok := p.Members[person.ID]
	if !ok {
		return
	}

	for key, value := range m {
		member.Meta[key] = value
	}
}

func (p *Party) UpdateMemberRevision(person *Person, revision int) {
	p.m.Lock()
	defer p.m.Unlock()

	member, ok := p.Members[person.ID]
	if !ok {
		return
	}

	member.Meta["revision"] = revision
}

func (p *Party) DeleteMemberMeta(person *Person, keys []string) {
	p.m.Lock()
	defer p.m.Unlock()

	member, ok := p.Members[person.ID]
	if !ok {
		return
	}

	for _, key := range keys {
		delete(member.Meta, key)
	}
}

func (p *Party) UpdateMemberConnections(person *Person, m aid.JSON) {
	p.m.Lock()
	defer p.m.Unlock()

	member, ok := p.Members[person.ID]
	if !ok {
		return
	}

	member.Connections[m["id"].(string)] = m
}

func (p *Party) UpdateConfig(m map[string]interface{}) {
	p.m.Lock()
	defer p.m.Unlock()

	for key, value := range m {
		p.Config[key] = value
	}
}

func (p *Party) GenerateFortniteParty() aid.JSON {
	p.m.Lock()
	defer p.m.Unlock()

	party := aid.JSON{
		"id": p.ID,
		"config": p.Config,
		"meta": p.Meta,
		"applicants": []aid.JSON{},
		"members": []aid.JSON{},
		"invites": []aid.JSON{},
		"intentions": []aid.JSON{},
		"created_at": p.CreatedAt.Format(time.RFC3339),
		"updated_at": time.Now().Format(time.RFC3339),
		"revision": 0,
	}

	for _, member := range p.Members {
		party["members"] = append(party["members"].([]aid.JSON), member.GenerateFortnitePartyMember())
	}

	return party
}