package handlers

import (
	"time"

	"github.com/ectrc/snow/aid"
	p "github.com/ectrc/snow/person"

	"github.com/gofiber/fiber/v2"
)

var (
	profileActions = map[string]func(c *fiber.Ctx, person *p.Person, profile *p.Profile) error {
		"QueryProfile": PostQueryProfileAction,
		"ClientQuestLogin": PostQueryProfileAction,
	}
)

func PostProfileAction(c *fiber.Ctx) error {
	person := p.Find(c.Params("accountId"))
	if person == nil {
		return c.Status(404).JSON(aid.ErrorBadRequest("No Account Found"))
	}

	profile := person.GetProfileFromType(c.Query("profileId"))
	if profile == nil {
		return c.Status(404).JSON(aid.ErrorBadRequest("No Profile Found"))
	}

	snapshot := profile.Snapshot()
	if action, ok := profileActions[c.Params("action")]; ok {
		err := action(c, person, profile)
		if err != nil {
			return err
		}
	}
	profile.Diff(snapshot)
	profile.Revision++

	return c.Status(200).JSON(aid.JSON{
		"profileId": profile.Type,
		"profileRevision": profile.Revision,
		"profileCommandRevision": profile.Revision,
		"profileChangesBaseRevision": profile.Revision - 1,
		"profileChanges": profile.Changes,
		"multiUpdate": []aid.JSON{},
		"notifications": []aid.JSON{},
		"responseVersion": 1,
		"serverTime": time.Now().Format("2006-01-02T15:04:05.999Z"),
	})
}

func PostQueryProfileAction(c *fiber.Ctx, person *p.Person, profile *p.Profile) error {
	profile.Changes = []interface{}{}
	profile.CreateFullProfileUpdateChange()
	return nil
}