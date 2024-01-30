package person

import (
	"fmt"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/storage"
)

type RelationshipDirection string
const RelationshipInboundDirection RelationshipDirection = "INBOUND"
const RelationshipOutboundDirection RelationshipDirection = "OUTBOUND"

type RelationshipGenerateType string
const GenerateTypeFromPerson RelationshipGenerateType = "FROM_PERSON"
const GenerateTypeTowardsPerson RelationshipGenerateType = "TOWARDS_PERSON"

type Relationship struct {
	From *Person
	Towards *Person
	Status string
	Direction RelationshipDirection
}

func (r *Relationship) ToDatabase() *storage.DB_Relationship {
	return &storage.DB_Relationship{
		FromPersonID: r.From.ID,
		TowardsPersonID: r.Towards.ID,
		Status: r.Status,
	}
}

func (r *Relationship) GenerateFortniteFriendEntry(t RelationshipGenerateType) aid.JSON {
	result := aid.JSON{
		"status": r.Status,
		"created": "0000-00-00T00:00:00.000Z",
		"favorite": false,
	}

	switch t {
	case GenerateTypeFromPerson:
		result["direction"] = "OUTBOUND"
		result["accountId"] = r.Towards.ID
	case GenerateTypeTowardsPerson:
		result["direction"] = "INBOUND"
		result["accountId"] = r.From.ID
	}

	return result
}

func (r *Relationship) Save() (*Relationship, error) {
	storage.Repo.Storage.SaveRelationship(r.ToDatabase())
	r.From.Relationships.Set(r.Towards.ID, r)
	r.Towards.Relationships.Set(r.From.ID, r)
	return r, nil
}

func (r *Relationship) Delete() error {
	storage.Repo.Storage.DeleteRelationship(r.ToDatabase())
	return nil
}

func (p *Person) LoadRelationships() {
	incoming := storage.Repo.Storage.GetIncomingRelationships(p.ID)
	for _, entry := range incoming {
		relationship := &Relationship{
			From: Find(entry.FromPersonID),
			Towards: p,
			Status: entry.Status,
			Direction: RelationshipInboundDirection,
		}

		p.Relationships.Set(entry.FromPersonID, relationship)
	}
}

func (p *Person) CreateRelationship(personId string) (*Relationship, error) {
	exists, okay := p.Relationships.Get(personId)
	if !okay {
		return p.createOutboundRelationship(personId)
	}

	if exists.Status != "PENDING" {
		return nil, fmt.Errorf("relationship already exists")
	}

	if exists.Towards.ID == p.ID {
		return p.createAcceptInboundRelationship(personId)
	}

	return nil, fmt.Errorf("relationship already exists")
}

func (p *Person) createOutboundRelationship(towards string) (*Relationship, error) {
	relationship := &Relationship{
		From: p,
		Towards: Find(towards),
		Status: "PENDING",
		Direction: RelationshipOutboundDirection,
	}
	return relationship.Save()
}

func (p *Person) createAcceptInboundRelationship(towards string) (*Relationship, error) {
	relationship := &Relationship{
		From: Find(towards),
		Towards: p,
		Status: "ACCEPTED",
		Direction: RelationshipInboundDirection,
	}
	return relationship.Save()
}