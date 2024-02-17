package handlers

import (
	"time"

	"github.com/ectrc/snow/aid"
	p "github.com/ectrc/snow/person"
	"github.com/gofiber/fiber/v2"
)

var (
	serverActions = map[string]func(c *fiber.Ctx, person *p.Person, profile *p.Profile, profileChanges, multiUpdate, notifications *[]aid.JSON) error{
		"QueryProfile": serverQueryProfileAction,
	}
)

func PostServerProfileAction(c *fiber.Ctx) error {
	person := p.Find(c.Params("accountId"))
	if person == nil {
		return c.Status(404).JSON(aid.ErrorBadRequest("No Account Found"))
	}

	profile := person.GetProfileFromType(c.Query("profileId"))
	if profile == nil {
		return c.Status(404).JSON(aid.ErrorBadRequest("No Profile Found"))
	}

	profileChanges := []aid.JSON{}
	multiUpdate := []aid.JSON{}
	notifications := []aid.JSON{}

	if action, ok := serverActions[c.Query("action")]; ok {
		if err := action(c, person, profile, &profileChanges, &multiUpdate, &notifications); err != nil {
			return c.Status(500).JSON(aid.ErrorBadRequest(err.Error()))
		}
	}

	return c.Status(200).JSON(aid.JSON{
		"profileId": c.Query("profileId"),
		"profileRevision": profile.Revision,
		"profileCommandRevision": profile.Revision,
		"profileChangesBaseRevision": profile.Revision - 1,
		"profileChanges": profileChanges,
		"multiUpdate": multiUpdate,
		"notifications": notifications,
		"responseVersion": 1,
		"serverTime": time.Now().Format("2006-01-02T15:04:05.999Z"),
	})
}

func serverQueryProfileAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, profileChanges, multiUpdate, notifications *[]aid.JSON) error {
	*profileChanges = append(*profileChanges, aid.JSON{
		"changeType": "fullProfileUpdate",
		"profile": profile.GenerateFortniteProfileEntry(),
	})

	return nil
}