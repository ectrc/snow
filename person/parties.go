package person

import (
	"sync"
	"time"

	"github.com/ectrc/snow/aid"
	"github.com/google/uuid"
)

type PartyInvite struct{
	Party *Party `json:"-"`
	Inviter *Person `json:"-"`
	Towards *Person `json:"-"`
	Meta map[string]interface{}
	Status string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}

func NewPartyInvite(party *Party, inviter *Person, towards *Person, meta map[string]interface{}) *PartyInvite {
	return &PartyInvite{
		Party: party,
		Inviter: inviter,
		Towards: towards,
		Meta: meta,
		Status: "SENT",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Minute * 60),
	}
}

func (pi *PartyInvite) GenerateFortnitePartyInvite() aid.JSON {
	return aid.JSON{
		"party_id": pi.Party.ID,
		"sent_by": pi.Inviter.ID,
		"sent_to": pi.Towards.ID,
		"sent_at": pi.CreatedAt.Format(time.RFC3339),
		"status": pi.Status,
		"meta": pi.Meta,
		"inviter_pl": "win",
		"inviter_pl_dn": pi.Inviter.DisplayName,
		"expires_at": pi.ExpiresAt.Format(time.RFC3339),
		"updated_at": pi.UpdatedAt.Format(time.RFC3339),
	}
}

type PartyMember struct{
	Person *Person
	ConnectionID string
	Meta map[string]interface{}
	Connection aid.JSON
	Role string
	JoinedAt time.Time
	UpdatedAt time.Time
}

func (pm *PartyMember) GenerateFortnitePartyMember() aid.JSON {
	conn := pm.Connection
	conn["yield_leadership"] = true

	return aid.JSON{
		"account_id": pm.Person.ID,
		"role": pm.Role,
		"meta": pm.Meta,
		"joined_at": pm.JoinedAt.Format(time.RFC3339),
		"connections": []aid.JSON{conn},
		"revision": 0,
	}
}

type Party struct{
	ID string
	Captain *PartyMember
	Members map[string]*PartyMember
	Invites []*PartyInvite
	Config map[string]interface{}
	Meta map[string]interface{}
	CreatedAt time.Time
	m sync.Mutex
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
		Invites: []*PartyInvite{},
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

func (p *Party) ChangeNewCaptain() {
	p.m.Lock()
	defer p.m.Unlock()

	for _, member := range p.Members {
		p.Captain = member
		p.Captain.Role = "CAPTAIN"
		p.Captain.UpdatedAt = time.Now()
		break
	}
}

func (p *Party) AddMember(person *Person) {
	p.m.Lock()
	defer p.m.Unlock()

	partyMember := &PartyMember{
		Person: person,
		Meta: make(map[string]interface{}),
		Role: "MEMBER",
		JoinedAt: time.Now(),
	}

	p.Members[person.ID] = partyMember
	person.Parties.Set(p.ID, p)
}

func (p *Party) RemoveMember(person *Person) {
	p.m.Lock()
	defer p.m.Unlock()

	person.Parties.Delete(p.ID)

	delete(p.Members, person.ID)
	if len(p.Members) == 0 {
		Parties.Delete(p.ID)
	}
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

func (p *Party) UpdateMemberConnection(person *Person, m aid.JSON) {
	p.m.Lock()
	defer p.m.Unlock()

	m["connected_at"] = time.Now().Format(time.RFC3339)
	m["updated_at"] = time.Now().Format(time.RFC3339)

	member, ok := p.Members[person.ID]
	if !ok {
		return
	}

	member.Connection = m
}

func (p *Party) UpdateConfig(m map[string]interface{}) {
	p.m.Lock()
	defer p.m.Unlock()

	for key, value := range m {
		p.Config[key] = value
	}
}

func (p *Party) AddInvite(invite *PartyInvite) {
	p.m.Lock()
	defer p.m.Unlock()

	p.Invites = append(p.Invites, invite)
}

func (p *Party) RemoveInvite(invite *PartyInvite) {
	p.m.Lock()
	defer p.m.Unlock()

	for i, v := range p.Invites {
		if v == invite {
			p.Invites = append(p.Invites[:i], p.Invites[i+1:]...)
			break
		}
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

	for _, invite := range p.Invites {
		party["invites"] = append(party["invites"].([]aid.JSON), invite.GenerateFortnitePartyInvite())
	}

	return party
}