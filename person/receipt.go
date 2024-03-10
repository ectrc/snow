package person

import (
	"time"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/storage"
	"github.com/google/uuid"
)

type Receipt struct {
	ID string
	PersonID string
	OfferID string
	PurchaseDate int64
	TotalPaid int
	State string
	Loot []*Item
}

func NewReceipt(offerID string, totalPaid int) *Receipt {
	return &Receipt{
		ID: uuid.New().String(),
		OfferID: offerID,
		PurchaseDate: time.Now().Unix(),
		TotalPaid: totalPaid,
		Loot: []*Item{},
		State: "PENDING",
	}
}

func FromDatabaseReceipt(receipt *storage.DB_Receipt) *Receipt {
	loot := []*Item{}

	for _, item := range receipt.Loot {
		loot = append(loot, FromDatabaseReceiptLoot(&item))
	}

	return &Receipt{
		ID: receipt.ID,
		PersonID: receipt.PersonID,
		OfferID: receipt.OfferID,
		PurchaseDate: receipt.PurchaseDate,
		TotalPaid: receipt.TotalPaid,
		State: receipt.State,
		Loot: loot,
	}
}

func (r *Receipt) GenerateUnrealReceiptEntry() aid.JSON {
	return aid.JSON{
		"TransactionId": r.ID,
		"TransactionState": string(r.State),
		"Offers": []aid.JSON{{
			"OfferNamespace": "fn",
			"OfferId": r.OfferID,
			"Items": []aid.JSON{{
				"EntitlementId": r.ID,
				"EntitlementName": "",
				"ItemId": r.OfferID,
				"ItemNamespace": "fn",
			}},
		}},
		"grantedVoucher": aid.JSON{},
	}
}

func (r *Receipt) GenerateFortniteReceiptEntry() aid.JSON {
	return aid.JSON{
		"receiptId": r.ID,
		"appStoreId": r.OfferID,
		"receiptInfo": r.State,
	}
}

func (r *Receipt) AddLoot(item *Item) {
	r.Loot = append(r.Loot, item)
}

func (r *Receipt) SetState(state string) {
	r.State = state
}

func (r *Receipt) Delete() {
	for _, item := range r.Loot {
		item.DeleteReceiptLoot()
	}

	storage.Repo.DeleteReceipt(r.ID)
}

func (r *Receipt) Save() {
	for _, item := range r.Loot {
		storage.Repo.SaveReceiptLoot(item.ToReceiptLootDatabase(r.ID))
	}
	storage.Repo.SaveReceipt(r.ToDatabase())
}

func (r *Receipt) ToDatabase() *storage.DB_Receipt {
	loot := []storage.DB_ReceiptLoot{}

	for _, item := range r.Loot {
		loot = append(loot, *item.ToReceiptLootDatabase(r.ID))
	}

	return &storage.DB_Receipt{
		ID: r.ID,
		PersonID: r.PersonID,
		OfferID: r.OfferID,
		PurchaseDate: r.PurchaseDate,
		TotalPaid: r.TotalPaid,
		State: r.State,
		Loot: loot,
	}
}