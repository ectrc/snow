package person

import (
	"fmt"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/storage"
)

type RelationshipDirection string

type RelationshipInboundDirection RelationshipDirection
const RelationshipInboundDirectionValue RelationshipInboundDirection = "INBOUND"

type RelationshipOutboundDirection RelationshipDirection
const RelationshipOutboundDirectionValue RelationshipOutboundDirection = "OUTBOUND"

type Relationship[T RelationshipInboundDirection | RelationshipOutboundDirection] struct {
	Me *Person
	Towards *Person
	Status string
	Direction T
}

func (r *Relationship[T]) ToDatabase() *storage.DB_Relationship {
	return &storage.DB_Relationship{
		IncomingPersonID: r.Me.ID,
		OutgoingPersonID: r.Towards.ID,
		Status: r.Status,
	}
}

func (r *Relationship[T]) GenerateFortniteFriendEntry() aid.JSON {
	return aid.JSON{
		"accountId": r.Towards.ID,
		"status": r.Status,
		"direction": string(r.Direction),
		"created": "0000-00-00T00:00:00.000Z",
		"favorite": false,
	}
}

func (r *Relationship[T]) Save() {
	storage.Repo.Storage.SaveRelationship(r.ToDatabase())
}

func (r *Relationship[T]) Delete() {
	storage.Repo.Storage.DeleteRelationship(r.ToDatabase())
}

func (p *Person) LoadRelationships() {
	incoming := storage.Repo.Storage.GetIncomingRelationships(p.ID)
	for _, entry := range incoming {
		relationship := &Relationship[RelationshipInboundDirection]{
			Status: entry.Status,
			Me: p,
			Towards: FindShallow(entry.OutgoingPersonID),
			Direction: RelationshipInboundDirectionValue,
		}

		p.IncomingRelationships.Set(entry.OutgoingPersonID, relationship)
	}

	outgoing := storage.Repo.Storage.GetOutgoingRelationships(p.ID)
	for _, entry := range outgoing {
		relationship := &Relationship[RelationshipOutboundDirection]{
			Status: entry.Status,
			Me: p,
			Towards: FindShallow(entry.IncomingPersonID),
			Direction: RelationshipOutboundDirectionValue,
		}

		p.OutgoingRelationships.Set(entry.IncomingPersonID, relationship)
	}
}

func (p *Person) CreateRelationship(personId string) (string, error) {
	if p.ID == personId {
		return "", fmt.Errorf("cannot create relationship with yourself")
	}

	if p.IncomingRelationships.Has(personId) {
		return "INBOUND", p.createAcceptInboundRelationship(personId)
	}

	return "OUTBOUND", p.createOutboundRelationship(personId)
}

func (p *Person) createOutboundRelationship(towards string) error {
	towardsPerson := Find(towards)
	if towardsPerson == nil {
		return fmt.Errorf("person not found")
	}

	relationship := &Relationship[RelationshipOutboundDirection]{
		Me: p,
		Towards: towardsPerson,
		Status: "PENDING",
		Direction: RelationshipOutboundDirectionValue,
	}
	relationship.Save()
	p.OutgoingRelationships.Set(towards, relationship)

	tempRelationship := &Relationship[RelationshipInboundDirection]{
		Me: towardsPerson,
		Towards: p,
		Status: "PENDING",
		Direction: RelationshipInboundDirectionValue,
	}
	tempRelationship.Save()
	towardsPerson.IncomingRelationships.Set(p.ID, tempRelationship)

	return nil
}

func (p *Person) createAcceptInboundRelationship(towards string) error {
	towardsPerson := Find(towards)
	if towardsPerson == nil {
		return fmt.Errorf("person not found")
	}

	relationship := &Relationship[RelationshipInboundDirection]{
		Me: p,
		Towards: towardsPerson,
		Status: "ACCEPTED",
		Direction: RelationshipInboundDirectionValue,
	}
	relationship.Save()
	p.IncomingRelationships.Set(towards, relationship)

	tempRelationship := &Relationship[RelationshipOutboundDirection]{
		Me: towardsPerson,
		Towards: p,
		Status: "ACCEPTED",
		Direction: RelationshipOutboundDirectionValue,
	}
	tempRelationship.Save()
	towardsPerson.OutgoingRelationships.Set(p.ID, tempRelationship)

	return nil
}