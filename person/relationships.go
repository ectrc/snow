package person

import (
	"fmt"
	"time"

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
		"created": time.Now().Add(-time.Hour * 24 * 3).Format(time.RFC3339),
		"favorite": false,
	}

	switch t {
	case GenerateTypeFromPerson:
		result["direction"] = RelationshipOutboundDirection
		result["accountId"] = r.Towards.ID
	case GenerateTypeTowardsPerson:
		result["direction"] = RelationshipInboundDirection
		result["accountId"] = r.From.ID
	}

	return result
}

func (r *Relationship) GenerateFortniteFriendRemovalEntry(t RelationshipGenerateType) aid.JSON {
	result := aid.JSON{
		"reason": "DELETED",
	}

	switch t {
	case GenerateTypeFromPerson:
		result["accountId"] = r.Towards.ID
	case GenerateTypeTowardsPerson:
		result["accountId"] = r.From.ID
	}

	return result
}


func (r *Relationship) GenerateFortniteSummaryEntry(t RelationshipGenerateType) aid.JSON {
	result := aid.JSON{
		"created": time.Now().Add(-time.Hour * 24 * 3).Format(time.RFC3339),
		"favorite": false,
		"groups": []string{},
		"mutual": 0,
		"note": "",
	}

	switch t {
	case GenerateTypeFromPerson:
		result["accountId"] = r.Towards.ID
	case GenerateTypeTowardsPerson:
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
			From: FindShallow(entry.FromPersonID),
			Towards: p,
			Status: entry.Status,
			Direction: RelationshipInboundDirection,
		}

		p.Relationships.Set(entry.FromPersonID, relationship)
	}

	outgoing := storage.Repo.Storage.GetOutgoingRelationships(p.ID)
	for _, entry := range outgoing {
		relationship := &Relationship{
			From: p,
			Towards: FindShallow(entry.TowardsPersonID),
			Status: entry.Status,
			Direction: RelationshipOutboundDirection,
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
		Towards: FindShallow(towards),
		Status: "PENDING",
		Direction: RelationshipOutboundDirection,
	}
	return relationship.Save()
}

func (p *Person) createAcceptInboundRelationship(towards string) (*Relationship, error) {
	relationship := &Relationship{
		From: FindShallow(towards),
		Towards: p,
		Status: "ACCEPTED",
		Direction: RelationshipInboundDirection,
	}
	return relationship.Save()
}