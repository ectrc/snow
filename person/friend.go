package person

import (
	"time"

	"github.com/ectrc/snow/aid"
)

type Friend struct {
	Person    *Person
	Status    string
	Direction string
}

func (f *Friend) GenerateSummaryResponse() aid.JSON {
	return aid.JSON{
		"accountId": f.Person.ID,
		"groups": []string{},
		"mutual": 0,
		"alias": "",
		"note": "",
		"favorite": false,
		"created": time.Now().Add(-time.Hour * 24 * 7).Format(time.RFC3339),
	}
}

func (f *Friend) GenerateFriendResponse() aid.JSON {
	return aid.JSON{
		"accountId": f.Person.ID,
		"status": f.Status,
		"direction": f.Direction,
		"created": time.Now().Add(-time.Hour * 24 * 7).Format(time.RFC3339),
		"favourite": false,
	}
}