package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/ectrc/snow/aid"
	p "github.com/ectrc/snow/person"

	"github.com/gofiber/fiber/v2"
)

var (
	profileActions = map[string]func(c *fiber.Ctx, person *p.Person, profile *p.Profile) error {
		"QueryProfile": PostQueryProfileAction,
		"ClientQuestLogin": PostQueryProfileAction,
		"MarkItemSeen": PostMarkItemSeenAction,
		"EquipBattleRoyaleCustomization": PostEquipBattleRoyaleCustomizationAction,
	}
)

func PostProfileAction(c *fiber.Ctx) error {
	person := p.Find(c.Params("accountId"))
	if person == nil {
		return c.Status(404).JSON(aid.ErrorBadRequest("No Account Found"))
	}

	profile := person.GetProfileFromType(c.Query("profileId"))
	defer profile.ClearProfileChanges()

	before := profile.Snapshot()
	if action, ok := profileActions[c.Params("action")]; ok {
		if err := action(c, person, profile); err != nil {
			return err
		}
	}
	profile.Diff(before)
	
	revision, _ := strconv.Atoi(c.Query("rvn"))
	if revision == -1 {
		revision = profile.Revision
	}
	revision++

	return c.Status(200).JSON(aid.JSON{
		"profileId": profile.Type,
		"profileRevision": revision,
		"profileCommandRevision": revision,
		"profileChangesBaseRevision": revision - 1,
		"profileChanges": profile.Changes,
		"multiUpdate": []aid.JSON{},
		"notifications": []aid.JSON{},
		"responseVersion": 1,
		"serverTime": time.Now().Format("2006-01-02T15:04:05.999Z"),
	})
}

func PostQueryProfileAction(c *fiber.Ctx, person *p.Person, profile *p.Profile) error {
	profile.CreateFullProfileUpdateChange()
	return nil
}

func PostMarkItemSeenAction(c *fiber.Ctx, person *p.Person, profile *p.Profile) error {
	var body struct {
		ItemIds []string `json:"itemIds"`
	}

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("Invalid Body"))
	}

	for _, itemId := range body.ItemIds {
		item := profile.Items.GetItem(itemId)
		if item == nil {
			continue
		}
		
		item.HasSeen = true
		go item.Save()
	}

	return nil
}

func PostEquipBattleRoyaleCustomizationAction(c *fiber.Ctx, person *p.Person, profile *p.Profile) error {
	var body struct {
		SlotName string `json:"slotName"`
		ItemToSlot string `json:"itemToSlot"`
		IndexWithinSlot int `json:"indexWithinSlot"`
	}

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("Invalid Body"))
	}

	item := profile.Items.GetItem(body.ItemToSlot)
	if item == nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("Item not found"))
	}

	attr := profile.Attributes.GetAttributeByKey("favorite_" + strings.ToLower(body.SlotName))
	if attr == nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("Attribute not found"))
	}
	defer func() {
		go attr.Save()
	}()

	switch body.SlotName {
	case "Dance":
		value := aid.JSONParse(attr.ValueJSON)
		value.([]any)[body.IndexWithinSlot] = item.ID
		attr.ValueJSON = aid.JSONStringify(value)
	case "ItemWrap":
		value := aid.JSONParse(attr.ValueJSON)
		value.([]any)[body.IndexWithinSlot] = item.ID
		attr.ValueJSON = aid.JSONStringify(value)
	default:
		attr.ValueJSON = aid.JSONStringify(item.ID)
	}

	return nil
}