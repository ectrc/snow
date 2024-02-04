package person

import (
	"time"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/storage"
	"github.com/google/uuid"
)

type Purchase struct {
	ID string
	ProfileID string
	PersonID string
	Loot []*Item
	OfferID string
	PurchaseDate time.Time
	FreeRefundExpiry time.Time
	RefundExpiry time.Time
	RefundedAt time.Time
	TotalPaid int
}

func NewPurchase(offerID string, price int) *Purchase {
	return &Purchase{
		ID: uuid.New().String(),
		OfferID: offerID,
		PurchaseDate: time.Now(),
		RefundedAt: time.Unix(0, 0),
		FreeRefundExpiry: time.Now().Add(time.Hour * 24),
		RefundExpiry: time.Now().Add(time.Hour * 24 * 30),
		Loot: []*Item{},
		TotalPaid: price,
	}
}

func FromDatabasePurchase(purchase *storage.DB_Purchase) *Purchase {
	loot := []*Item{}

	for _, item := range purchase.Loot {
		loot = append(loot, FromDatabasePurchaseLoot(&item))
	}

	return &Purchase{
		ID: purchase.ID,
		ProfileID: purchase.ProfileID,
		Loot: loot,
		OfferID: purchase.OfferID,
		PurchaseDate: time.Unix(purchase.PurchaseDate, 0),
		FreeRefundExpiry: time.Unix(purchase.FreeRefundExpiry, 0),
		RefundExpiry: time.Unix(purchase.RefundExpiry, 0),
		RefundedAt: time.Unix(purchase.RefundedAt, 0),
	}
}

func (p *Purchase) GenerateFortnitePurchaseEntry() aid.JSON {
	json := aid.JSON{
		"offerId": p.OfferID,
		"freeRefundEligible": time.Now().Before(p.FreeRefundExpiry),
		"purchaseId": p.ID,
		"purchaseDate": p.PurchaseDate.Format(time.RFC3339),
		"undoTimeout": p.FreeRefundExpiry.Format(time.RFC3339),
		"totalMtxPaid": p.TotalPaid,
		"lootResult": []aid.JSON{},
		"gameContext": "",
		"metadata": aid.JSON{},
		"fulfillments": []aid.JSON{},
	}

	for _, item := range p.Loot {
		json["lootResult"] = append(json["lootResult"].([]aid.JSON), aid.JSON{
			"itemGuid": item.ID,
			"itemType": item.TemplateID,
			"itemProfile": item.ProfileType,
			"quantity": item.Quantity,
		})
	}

	if p.RefundedAt.Unix() > 0 {
		json["refundDate"] = p.RefundedAt.Format(time.RFC3339)
	}

	return json
}

func (p *Purchase) AddLoot(item *Item) {
	p.Loot = append(p.Loot, item)
}

func (p *Purchase) Delete() {
	storage.Repo.DeletePurchase(p.ID)
}

func (p *Purchase) ToDatabase(profileId string) *storage.DB_Purchase {
	loot := []storage.DB_PurchaseLoot{}

	for _, item := range p.Loot {
		loot = append(loot, *item.ToPurchaseLootDatabase(p.ID))
	}

	return &storage.DB_Purchase{
		ID: p.ID,
		ProfileID: profileId,
		Loot: loot,
		OfferID: p.OfferID,
		PurchaseDate: p.PurchaseDate.Unix(),
		FreeRefundExpiry: p.FreeRefundExpiry.Unix(),
		RefundExpiry: p.RefundExpiry.Unix(),
		RefundedAt: p.RefundedAt.Unix(),
		TotalPaid: p.TotalPaid,
	}
}

func (p *Purchase) Save() {
	storage.Repo.SavePurchase(p.ToDatabase(p.ProfileID))
	Find(p.PersonID).SetPurchaseHistoryAttribute()
}