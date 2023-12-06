package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/fortnite"
	p "github.com/ectrc/snow/person"

	"github.com/gofiber/fiber/v2"
)

var (
	profileActions = map[string]func(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
		"QueryProfile": PostQueryProfileAction,
		"ClientQuestLogin": PostQueryProfileAction,
		"MarkItemSeen": PostMarkItemSeenAction,
		"SetItemFavoriteStatusBatch": PostSetItemFavoriteStatusBatchAction,
		"EquipBattleRoyaleCustomization": PostEquipBattleRoyaleCustomizationAction,
		"SetBattleRoyaleBanner": PostSetBattleRoyaleBannerAction,
		"SetCosmeticLockerSlot": PostSetCosmeticLockerSlotAction,
		"SetCosmeticLockerBanner": PostSetCosmeticLockerBannerAction,
		"PurchaseCatalogEntry": PostPurchaseCatalogEntryAction,
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
	defer profile.ClearProfileChanges()

	profileSnapshots := map[string]*p.ProfileSnapshot{
		"athena": nil,
		"common_core": nil,
		"common_public": nil,
	}

	for key := range profileSnapshots {
		profileSnapshots[key] = person.GetProfileFromType(key).Snapshot()
	}

	notifications := []aid.JSON{}

	action, ok := profileActions[c.Params("action")];
	if ok && profile != nil {
		if err := action(c, person, profile, &notifications); err != nil {
			return c.Status(400).JSON(aid.ErrorBadRequest(err.Error()))
		}
	}

	for key, profileSnapshot := range profileSnapshots {
		profile := person.GetProfileFromType(key)
		if profile == nil {
			continue
		}

		if profileSnapshot == nil {
			continue
		}

		profile.Diff(profileSnapshot)
	}
	
	revision, _ := strconv.Atoi(c.Query("rvn"))
	if revision == -1 {
		revision = profile.Revision
	}
	revision++
	profile.Revision = revision
	profile.Save()
	delete(profileSnapshots, profile.Type)

	multiUpdate := []aid.JSON{}
	for key := range profileSnapshots {
		profile := person.GetProfileFromType(key)
		if profile == nil {
			continue
		}
		profile.Revision++

		if len(profile.Changes) == 0 {
			continue
		}

		multiUpdate = append(multiUpdate, aid.JSON{
			"profileId": profile.Type,
			"profileRevision": profile.Revision,
			"profileCommandRevision": profile.Revision,
			"profileChangesBaseRevision": profile.Revision - 1,
			"profileChanges": profile.Changes,
		})

		profile.ClearProfileChanges()
		profile.Save()
	}

	return c.Status(200).JSON(aid.JSON{
		"profileId": c.Query("profileId"),
		"profileRevision": revision,
		"profileCommandRevision": revision,
		"profileChangesBaseRevision": revision - 1,
		"profileChanges": profile.Changes,
		"multiUpdate": multiUpdate,
		"notifications": notifications,
		"responseVersion": 1,
		"serverTime": time.Now().Format("2006-01-02T15:04:05.999Z"),
	})
}

func PostQueryProfileAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	profile.CreateFullProfileUpdateChange()
	return nil
}

func PostMarkItemSeenAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	var body struct {
		ItemIds []string `json:"itemIds"`
	}

	err := c.BodyParser(&body)
	if err != nil {
		return fmt.Errorf("invalid Body")
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

func PostEquipBattleRoyaleCustomizationAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	var body struct {
		SlotName string `json:"slotName" binding:"required"`
		ItemToSlot string `json:"itemToSlot"`
		IndexWithinSlot int `json:"indexWithinSlot"`
	}

	err := c.BodyParser(&body)
	if err != nil {
		return fmt.Errorf("invalid Body")
	}

	item := profile.Items.GetItem(body.ItemToSlot)
	if item == nil {
		if body.ItemToSlot != "" && !strings.Contains(strings.ToLower(body.ItemToSlot), "random") {
			return fmt.Errorf("item not found")
		}

		item = &p.Item{
			ID: body.ItemToSlot,
		}
	}

	attr := profile.Attributes.GetAttributeByKey("favorite_" + strings.ToLower(body.SlotName))
	if attr == nil {
		return fmt.Errorf("attribute not found")
	}

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

	go attr.Save()
	return nil
}

func PostSetBattleRoyaleBannerAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	var body struct {
		HomebaseBannerColorID string `json:"homebaseBannerColorId" binding:"required"`
		HomebaseBannerIconID string `json:"homebaseBannerIconId" binding:"required"`
	}
	
	err := c.BodyParser(&body)
	if err != nil {
		return fmt.Errorf("invalid Body")
	}

	colorItem := person.CommonCoreProfile.Items.GetItemByTemplateID("HomebaseBannerColor:"+body.HomebaseBannerColorID)
	if colorItem == nil {
		return fmt.Errorf("color item not found")
	}

	iconItem := person.CommonCoreProfile.Items.GetItemByTemplateID("HomebaseBannerIcon:"+body.HomebaseBannerIconID)
	if iconItem == nil {
		return fmt.Errorf("icon item not found")
	}

	iconAttr := profile.Attributes.GetAttributeByKey("banner_icon")
	if iconAttr == nil {
		return fmt.Errorf("icon attribute not found")
	}

	colorAttr := profile.Attributes.GetAttributeByKey("banner_color")
	if colorAttr == nil {
		return fmt.Errorf("color attribute not found")
	}

	iconAttr.ValueJSON = aid.JSONStringify(strings.Split(iconItem.TemplateID, ":")[1])
	colorAttr.ValueJSON = aid.JSONStringify(strings.Split(colorItem.TemplateID, ":")[1])

	go func() {
		iconAttr.Save()
		colorAttr.Save()
	}()
	return nil
}

func PostSetItemFavoriteStatusBatchAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	var body struct {
		ItemIds []string `json:"itemIds" binding:"required"`
		Favorite []bool `json:"itemFavStatus" binding:"required"`
	}

	err := c.BodyParser(&body)
	if err != nil {
		return fmt.Errorf("invalid Body")
	}

	for i, itemId := range body.ItemIds {
		item := profile.Items.GetItem(itemId)
		if item == nil {
			continue
		}

		item.Favorite = body.Favorite[i]
		go item.Save()
	}

	return nil
}

func PostSetCosmeticLockerSlotAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error { 
	var body struct {
		Category string `json:"category" binding:"required"` // item type e.g. Character
		ItemToSlot string `json:"itemToSlot" binding:"required"` // template id
		LockerItem string `json:"lockerItem" binding:"required"` // locker id
		SlotIndex int `json:"slotIndex" binding:"required"` // index of slot
		VariantUpdates []aid.JSON `json:"variantUpdates" binding:"required"` // variant updates
	}

	err := c.BodyParser(&body)
	if err != nil {
		return fmt.Errorf("invalid Body")
	}

	item := profile.Items.GetItemByTemplateID(body.ItemToSlot)
	if item == nil {
		if body.ItemToSlot != "" && !strings.Contains(strings.ToLower(body.ItemToSlot), "random") {
			return fmt.Errorf("item not found")
		} 

		item = &p.Item{
			ID: body.ItemToSlot,
		}
	}

	currentLocker := profile.Loadouts.GetLoadout(body.LockerItem)
	if currentLocker == nil {
		return fmt.Errorf("current locker not found")
	}

	switch body.Category {
	case "Character":
		currentLocker.CharacterID = item.ID
	case "Backpack":
		currentLocker.BackpackID = item.ID
	case "Pickaxe":
		currentLocker.PickaxeID = item.ID
	case "Glider":
		currentLocker.GliderID = item.ID
	case "ItemWrap":
		defer profile.CreateLoadoutChangedChange(currentLocker, "ItemWrapID")
		if body.SlotIndex == -1 {
			for i := range currentLocker.ItemWrapID {
				currentLocker.ItemWrapID[i] = item.ID
			}
			break
		}
		currentLocker.ItemWrapID[body.SlotIndex] = item.ID
	case "Dance":
		defer profile.CreateLoadoutChangedChange(currentLocker, "DanceID")
		if body.SlotIndex == -1 {
			for i := range currentLocker.DanceID {
				currentLocker.DanceID[i] = item.ID
			}
			break
		}
		currentLocker.DanceID[body.SlotIndex] = item.ID
	case "SkyDiveContrail":
		currentLocker.ContrailID = item.ID
	case "LoadingScreen":
		currentLocker.LoadingScreenID = item.ID
	case "MusicPack":
		currentLocker.MusicPackID = item.ID
	}

	go currentLocker.Save()	
	return nil
}

func PostSetCosmeticLockerBannerAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error { 
	var body struct {
		LockerItem string `json:"lockerItem" binding:"required"` // locker id
		BannerColorTemplateName string `json:"bannerColorTemplateName" binding:"required"` // template id
		BannerIconTemplateName string `json:"bannerIconTemplateName" binding:"required"` // template id
	}

	err := c.BodyParser(&body)
	if err != nil {
		return fmt.Errorf("invalid Body")
	}

	color := person.CommonCoreProfile.Items.GetItemByTemplateID("HomebaseBannerColor:" + body.BannerColorTemplateName)
	if color == nil {
		return fmt.Errorf("color item not found")
	}

	icon := profile.Items.GetItemByTemplateID("HomebaseBannerIcon:" + body.BannerIconTemplateName)
	if icon == nil {
		icon = &p.Item{
			ID: body.BannerIconTemplateName,
		}
	}

	currentLocker := profile.Loadouts.GetLoadout(body.LockerItem)
	if currentLocker == nil {
		return fmt.Errorf("current locker not found")
	}

	currentLocker.BannerColorID = color.ID
	currentLocker.BannerID = icon.ID

	go currentLocker.Save()
	return nil
}

func PostPurchaseCatalogEntryAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	var body struct{
		OfferID string `json:"offerId" binding:"required"`
		PurchaseQuantity int `json:"purchaseQuantity" binding:"required"`
		ExpectedTotalPrice int `json:"expectedTotalPrice" binding:"required"`
	}

	err := c.BodyParser(&body)
	if err != nil {
		return fmt.Errorf("invalid Body")
	}

	offer := fortnite.StaticCatalog.GetOfferById(body.OfferID)
	if offer == nil {
		return fmt.Errorf("offer not found")
	}

	if offer.Price != body.ExpectedTotalPrice {
		return fmt.Errorf("invalid price")
	}

	vbucks := profile.Items.GetItemByTemplateID("Currency:MtxPurchased")
	if vbucks == nil {
		return fmt.Errorf("vbucks not found")
	}

	profile0Vbucks := person.Profile0Profile.Items.GetItemByTemplateID("Currency:MtxPurchased")
	if profile0Vbucks == nil {
		return fmt.Errorf("profile0vbucks not found")
	}

	if vbucks.Quantity < body.ExpectedTotalPrice {
		return fmt.Errorf("not enough vbucks")
	}

	vbucks.Quantity -= body.ExpectedTotalPrice

	go func() {
		profile0Vbucks.Quantity = vbucks.Quantity // for season 2 and lower
		vbucks.Save()
		profile0Vbucks.Save()
	}()

	if offer.ProfileType != "athena" {
		return fmt.Errorf("save the world not implemeted yet")
	}

	loot := []aid.JSON{}
	for i := 0; i < body.PurchaseQuantity; i++ {
		for _, grant := range offer.Grants {
			if profile.Items.GetItemByTemplateID(grant) != nil {
				item := profile.Items.GetItemByTemplateID(grant)
				item.Quantity++
				go item.Save()

				continue
			}

			item := p.NewItem(grant, 1)
			person.AthenaProfile.Items.AddItem(item)

			loot = append(loot, aid.JSON{
				"itemType": item.TemplateID,
				"itemGuid": item.ID,
				"quantity": item.Quantity,
				"itemProfile": offer.ProfileType,
			})
		}
	}

	*notifications = append(*notifications, aid.JSON{
		"type": "CatalogPurchase",
		"lootResult": aid.JSON{
			"items": loot,
		},
		"primary": true,
	})

	return nil
}