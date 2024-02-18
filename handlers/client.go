package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/fortnite"
	p "github.com/ectrc/snow/person"
	"github.com/ectrc/snow/socket"

	"github.com/gofiber/fiber/v2"
)

var (
	clientActions = map[string]func(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
		"QueryProfile": clientQueryProfileAction,
		"ClientQuestLogin": clientClientQuestLoginAction,
		"MarkItemSeen": clientMarkItemSeenAction,
		"SetItemFavoriteStatusBatch": clientSetItemFavoriteStatusBatchAction,
		"EquipBattleRoyaleCustomization": clientEquipBattleRoyaleCustomizationAction,
		"SetBattleRoyaleBanner": clientSetBattleRoyaleBannerAction,
		"SetCosmeticLockerSlot": clientSetCosmeticLockerSlotAction,
		"SetCosmeticLockerBanner": clientSetCosmeticLockerBannerAction,
		"SetCosmeticLockerName": clientSetCosmeticLockerNameAction,
		"CopyCosmeticLoadout": clientCopyCosmeticLoadoutAction,
		"DeleteCosmeticLoadout": clientDeleteCosmeticLoadoutAction,
		"PurchaseCatalogEntry": clientPurchaseCatalogEntryAction,
		"RefundMtxPurchase": clientRefundMtxPurchaseAction,
		"GiftCatalogEntry": clientGiftCatalogEntryAction,
		"RemoveGiftBox": clientRemoveGiftBoxAction,
		"SetAffiliateName": clientSetAffiliateNameAction,
		"SetReceiveGiftsEnabled": clientSetReceiveGiftsEnabledAction,
	}
)

func PostClientProfileAction(c *fiber.Ctx) error {
	person := c.Locals("person").(*p.Person)
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

	action, ok := clientActions[c.Params("action")];
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
	go profile.Save()
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
		go profile.Save()
	}

	return c.Status(200).JSON(aid.JSON{
		"profileId": c.Query("profileId"),
		"profileRevision": profile.Revision,
		"profileCommandRevision": profile.Revision,
		"profileChangesBaseRevision": profile.Revision - 1,
		"profileChanges": profile.Changes,
		"multiUpdate": multiUpdate,
		"notifications": notifications,
		"responseVersion": 1,
		"serverTime": time.Now().Format("2006-01-02T15:04:05.999Z"),
	})
}

func clientQueryProfileAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	profile.CreateFullProfileUpdateChange()
	return nil
}

func clientClientQuestLoginAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	return nil
}

func clientMarkItemSeenAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	var body struct {
		ItemIds []string `json:"itemIds"`
	}

	if err := c.BodyParser(&body); err != nil {
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

func clientEquipBattleRoyaleCustomizationAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	var body struct {
		SlotName string `json:"slotName" binding:"required"`
		ItemToSlot string `json:"itemToSlot"`
		IndexWithinSlot int `json:"indexWithinSlot"`
		VariantUpdates []struct{
			Active string `json:"active"`
			Channel string `json:"channel"`
		} `json:"variantUpdates"`
	}

	if err := c.BodyParser(&body); err != nil {
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

	for _, update := range body.VariantUpdates {
		channel := item.GetChannel(update.Channel)
		if channel == nil {
			continue
		}

		channel.Active = update.Active
		go channel.Save()
	}

	attr := profile.Attributes.GetAttributeByKey("favorite_" + strings.ReplaceAll(strings.ToLower(body.SlotName), "wrap", "wraps"))
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
		if body.IndexWithinSlot == -1 {
			attr.ValueJSON = aid.JSONStringify([]any{item.ID,item.ID,item.ID,item.ID,item.ID,item.ID,item.ID})
			break
		}
		value.([]any)[body.IndexWithinSlot] = item.ID
		attr.ValueJSON = aid.JSONStringify(value)
	default:
		attr.ValueJSON = aid.JSONStringify(item.ID)
	}

	go attr.Save()
	return nil
}

func clientSetBattleRoyaleBannerAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	var body struct {
		HomebaseBannerColorID string `json:"homebaseBannerColorId" binding:"required"`
		HomebaseBannerIconID string `json:"homebaseBannerIconId" binding:"required"`
	}
	
	if err := c.BodyParser(&body); err != nil {
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
	iconAttr.Save()
	colorAttr.Save()

	return nil
}

func clientSetItemFavoriteStatusBatchAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	var body struct {
		ItemIds []string `json:"itemIds" binding:"required"`
		Favorite []bool `json:"itemFavStatus" binding:"required"`
	}

	if err := c.BodyParser(&body); err != nil {
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

func clientSetCosmeticLockerSlotAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error { 
	var body struct {
		Category string `json:"category" binding:"required"` // item type e.g. Character
		ItemToSlot string `json:"itemToSlot" binding:"required"` // template id
		LockerItem string `json:"lockerItem" binding:"required"` // locker id
		SlotIndex int `json:"slotIndex" binding:"required"` // index of slot
		VariantUpdates []struct{
			Active string `json:"active"`
			Channel string `json:"channel"`
		} `json:"variantUpdates" binding:"required"` // variant updates
	}

	if err := c.BodyParser(&body); err != nil {
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

	for _, update := range body.VariantUpdates {
		channel := item.GetChannel(update.Channel)
		if channel == nil {
			continue
		}

		channel.Active = update.Active
		go channel.Save()
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

func clientSetCosmeticLockerBannerAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error { 
	var body struct {
		LockerItem string `json:"lockerItem" binding:"required"` // locker id
		BannerColorTemplateName string `json:"bannerColorTemplateName" binding:"required"` // template id
		BannerIconTemplateName string `json:"bannerIconTemplateName" binding:"required"` // template id
	}

	if err := c.BodyParser(&body); err != nil {
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

func clientSetCosmeticLockerNameAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	var body struct {
		LockerItem string `json:"lockerItem" binding:"required"`
		Name string `json:"name" binding:"required"`
	}

	if err := c.BodyParser(&body); err != nil {
		return fmt.Errorf("invalid Body")
	}

	loadoutsAttribute := profile.Attributes.GetAttributeByKey("loadouts")
	if loadoutsAttribute == nil {
		return fmt.Errorf("loadouts not found")
	}
	loadouts := p.AttributeConvertToSlice[string](loadoutsAttribute)

	currentLocker := profile.Loadouts.GetLoadout(body.LockerItem)
	if currentLocker == nil {
		return fmt.Errorf("current locker not found")
	}

	if loadouts[0] == currentLocker.ID {
		return fmt.Errorf("cannot rename default locker")
	}

	currentLocker.LockerName = body.Name
	go currentLocker.Save()

	return nil
}

func clientCopyCosmeticLoadoutAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	var body struct {
		OptNewNameForTarget string `json:"optNewNameForTarget" binding:"required"`
		SourceIndex int `json:"sourceIndex" binding:"required"`
		TargetIndex int `json:"targetIndex" binding:"required"`
	}

	if err := c.BodyParser(&body); err != nil {
		return fmt.Errorf("invalid Body")
	}

	lastAppliedLoadoutAttribute := profile.Attributes.GetAttributeByKey("last_applied_loadout")
	if lastAppliedLoadoutAttribute == nil {
		return fmt.Errorf("last_applied_loadout not found")
	}

	activeLoadoutIndexAttribute := profile.Attributes.GetAttributeByKey("active_loadout_index")
	if activeLoadoutIndexAttribute == nil {
		return fmt.Errorf("active_loadout_index not found")
	}

	loadoutsAttribute := profile.Attributes.GetAttributeByKey("loadouts")
	if loadoutsAttribute == nil {
		return fmt.Errorf("loadouts not found")
	}
	loadouts := p.AttributeConvertToSlice[string](loadoutsAttribute)

	if body.SourceIndex >= len(loadouts) {
		return fmt.Errorf("source index out of range")
	}

	sandboxLoadout := profile.Loadouts.GetLoadout(loadouts[0])
	if sandboxLoadout == nil {
		return fmt.Errorf("sandbox loadout not found")
	}

	lastAppliedLoadout := profile.Loadouts.GetLoadout(p.AttributeConvert[string](lastAppliedLoadoutAttribute))
	if lastAppliedLoadout == nil {
		return fmt.Errorf("last applied loadout not found")
	}

	if body.TargetIndex >= len(loadouts) {
		newLoadout := p.NewLoadout(body.OptNewNameForTarget, profile)
		newLoadout.CopyFrom(lastAppliedLoadout)
		profile.Loadouts.AddLoadout(newLoadout)
		go newLoadout.Save()

		lastAppliedLoadout.CopyFrom(sandboxLoadout)
		go lastAppliedLoadout.Save()

		lastAppliedLoadoutAttribute.ValueJSON = aid.JSONStringify(newLoadout.ID)
		activeLoadoutIndexAttribute.ValueJSON = aid.JSONStringify(body.TargetIndex)
		go lastAppliedLoadoutAttribute.Save()
		go activeLoadoutIndexAttribute.Save()

		loadouts = append(loadouts, newLoadout.ID)
		loadoutsAttribute.ValueJSON = aid.JSONStringify(loadouts)
		go loadoutsAttribute.Save()

		sandboxLoadout.CopyFrom(newLoadout)
		go sandboxLoadout.Save()

		if len(profile.Changes) == 0 {
			profile.CreateLoadoutChangedChange(sandboxLoadout, "DanceID")
		}

		return nil
	}

	if body.SourceIndex > 0  {
		sourceLoadout := profile.Loadouts.GetLoadout(loadouts[body.SourceIndex])
		if sourceLoadout == nil {
			return fmt.Errorf("target loadout not found")
		}
	
		sandboxLoadout.CopyFrom(sourceLoadout)
		go sandboxLoadout.Save()

		lastAppliedLoadoutAttribute.ValueJSON = aid.JSONStringify(sourceLoadout.ID)
		activeLoadoutIndexAttribute.ValueJSON = aid.JSONStringify(body.SourceIndex)

		go lastAppliedLoadoutAttribute.Save()
		go activeLoadoutIndexAttribute.Save()

		if len(profile.Changes) == 0{
			profile.CreateLoadoutChangedChange(sandboxLoadout, "DanceID")
			profile.CreateLoadoutChangedChange(sourceLoadout, "DanceID")
		}

		return nil
	}

	targetLoadout := profile.Loadouts.GetLoadout(loadouts[body.TargetIndex])
	if targetLoadout == nil {
		return fmt.Errorf("target loadout not found")
	}

	sandboxLoadout.CopyFrom(targetLoadout)
	go sandboxLoadout.Save()

	if len(profile.Changes) == 0{
		profile.CreateLoadoutChangedChange(sandboxLoadout, "DanceID")
		profile.CreateLoadoutChangedChange(targetLoadout, "DanceID")
	}

	return nil
}

func clientDeleteCosmeticLoadoutAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	var body struct {
		FallbackLoadoutIndex int `json:"fallbackLoadoutIndex" binding:"required"`
		LoadoutIndex int `json:"index" binding:"required"`
	}

	if err := c.BodyParser(&body); err != nil {
		return fmt.Errorf("invalid Body")
	}

	lastAppliedLoadoutAttribute := profile.Attributes.GetAttributeByKey("last_applied_loadout")
	if lastAppliedLoadoutAttribute == nil {
		return fmt.Errorf("last_applied_loadout not found")
	}

	activeLoadoutIndexAttribute := profile.Attributes.GetAttributeByKey("active_loadout_index")
	if activeLoadoutIndexAttribute == nil {
		return fmt.Errorf("active_loadout_index not found")
	}

	loadoutsAttribute := profile.Attributes.GetAttributeByKey("loadouts")
	if loadoutsAttribute == nil {
		return fmt.Errorf("loadouts not found")
	}
	loadouts := p.AttributeConvertToSlice[string](loadoutsAttribute)

	if body.LoadoutIndex >= len(loadouts) {
		return fmt.Errorf("loadout index out of range")
	}

	if body.LoadoutIndex == 0 {
		return fmt.Errorf("cannot delete default loadout")
	}

	if body.FallbackLoadoutIndex == -1 {
		body.FallbackLoadoutIndex = 0
	}

	fallbackLoadout := profile.Loadouts.GetLoadout(loadouts[body.FallbackLoadoutIndex])
	if fallbackLoadout == nil {
		return fmt.Errorf("fallback loadout not found")
	}

	lastAppliedLoadoutAttribute.ValueJSON = aid.JSONStringify(fallbackLoadout.ID)
	activeLoadoutIndexAttribute.ValueJSON = aid.JSONStringify(body.FallbackLoadoutIndex)
	lastAppliedLoadoutAttribute.Save()
	activeLoadoutIndexAttribute.Save()

	profile.Loadouts.DeleteLoadout(loadouts[body.LoadoutIndex])
	loadouts = append(loadouts[:body.LoadoutIndex], loadouts[body.LoadoutIndex+1:]...)
	loadoutsAttribute.ValueJSON = aid.JSONStringify(loadouts)
	loadoutsAttribute.Save()

	return nil
}

func clientPurchaseCatalogEntryAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	var body struct {
		OfferID string `json:"offerId" binding:"required"`
		PurchaseQuantity int `json:"purchaseQuantity" binding:"required"`
		ExpectedTotalPrice int `json:"expectedTotalPrice" binding:"required"`
	}

	if err := c.BodyParser(&body); err != nil {
		return fmt.Errorf("invalid Body")
	}

	offer := fortnite.GetOfferByOfferId(body.OfferID)
	if offer == nil {
		return fmt.Errorf("offer not found")
	}

	if offer.TotalPrice != body.ExpectedTotalPrice {
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
	profile0Vbucks.Quantity = vbucks.Quantity
	vbucks.Save()
	profile0Vbucks.Save()

	if offer.Meta.ProfileId != "athena" {
		return fmt.Errorf("save the world not implemeted yet")
	}

	loot := []aid.JSON{}
	purchase := p.NewPurchase(body.OfferID, body.ExpectedTotalPrice)
	for i := 0; i < body.PurchaseQuantity; i++ {
		for _, grant := range offer.Grants {
			templateId := grant.Type.BackendValue + ":" + grant.ID
			if profile.Items.GetItemByTemplateID(templateId) != nil {
				item := profile.Items.GetItemByTemplateID(templateId)
				item.Quantity++
				go item.Save()

				continue
			}

			item := p.NewItem(templateId, 1)
			person.AthenaProfile.Items.AddItem(item)
			purchase.AddLoot(item)

			loot = append(loot, aid.JSON{
				"itemType": item.TemplateID,
				"itemGuid": item.ID,
				"quantity": item.Quantity,
				"itemProfile": offer.Meta.ProfileId,
			})
		}
	}

	person.AthenaProfile.Purchases.AddPurchase(purchase).Save()

	*notifications = append(*notifications, aid.JSON{
		"type": "CatalogPurchase",
		"lootResult": aid.JSON{
			"items": loot,
		},
		"primary": true,
	})

	affiliate := person.CommonCoreProfile.Attributes.GetAttributeByKey("mtx_affiliate")
	if affiliate == nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("Invalid affiliate attribute"))
	}

	creator := p.Find(p.AttributeConvert[string](affiliate))
	if creator != nil {
		creator.CommonCoreProfile.Items.GetItemByTemplateID("Currency:MtxPurchased").Quantity += body.ExpectedTotalPrice
		creator.Profile0Profile.Items.GetItemByTemplateID("Currency:MtxPurchased").Quantity += body.ExpectedTotalPrice
	}

	return nil
}

func clientRefundMtxPurchaseAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	var body struct {
		PurchaseID string `json:"purchaseId" binding:"required"`
	}

	if err := c.BodyParser(&body); err != nil {
		return fmt.Errorf("invalid Body")
	}

	purchase := person.AthenaProfile.Purchases.GetPurchase(body.PurchaseID)
	if purchase == nil {
		return fmt.Errorf("purchase not found")
	}

	if person.RefundTickets <= 0 {
		return fmt.Errorf("not enough refund tickets")
	}

	if time.Now().After(purchase.FreeRefundExpiry) {
		person.RefundTickets--
	}

	for _, item := range purchase.Loot {
		profile.Items.DeleteItem(item.ID)
	}
	
	purchase.RefundedAt = time.Now()
	purchase.Save()

	vbucks := profile.Items.GetItemByTemplateID("Currency:MtxPurchased")
	if vbucks == nil {
		return fmt.Errorf("vbucks not found")
	}

	vbucks.Quantity += purchase.TotalPaid
	vbucks.Save()

	profile0Vbucks := person.Profile0Profile.Items.GetItemByTemplateID("Currency:MtxPurchased")
	if profile0Vbucks == nil {
		return fmt.Errorf("profile0vbucks not found")
	}

	profile0Vbucks.Quantity = vbucks.Quantity
	profile0Vbucks.Save()

	return nil
}

func clientGiftCatalogEntryAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	var body struct {
		Currency string `json:"currency" binding:"required"`
		CurrencySubType string `json:"currencySubType" binding:"required"`
		ExpectedTotalPrice int `json:"expectedTotalPrice" binding:"required"`
		GameContext string `json:"gameContext" binding:"required"`
		GiftWrapTemplateId string `json:"giftWrapTemplateId" binding:"required"`
		PersonalMessage string `json:"personalMessage" binding:"required"`
		ReceiverAccountIds []string `json:"receiverAccountIds" binding:"required"`
		OfferId string `json:"offerId" binding:"required"`
	}

	if err := c.BodyParser(&body); err != nil {
		return fmt.Errorf("invalid Body")
	}

	offer := fortnite.GetOfferByOfferId(body.OfferId)
	if offer == nil {
		return fmt.Errorf("offer not found")
	}

	if offer.TotalPrice != body.ExpectedTotalPrice {
		return fmt.Errorf("invalid price")
	}

	for _, receiverAccountId := range body.ReceiverAccountIds {
		receiverPerson := p.Find(receiverAccountId)
		if receiverPerson == nil {
			return fmt.Errorf("one or more receivers not found")
		}

		for _, grant := range offer.Grants {
			if receiverPerson.AthenaProfile.Items.GetItemByTemplateID(grant.Type.BackendValue + ":" + grant.ID) != nil {
				return fmt.Errorf("one or more receivers has one of the items")
			}
		}
	}

	price := offer.TotalPrice * len(body.ReceiverAccountIds)

	vbucks := profile.Items.GetItemByTemplateID("Currency:MtxPurchased")
	if vbucks == nil {
		return fmt.Errorf("vbucks not found")
	}

	profile0Vbucks := person.Profile0Profile.Items.GetItemByTemplateID("Currency:MtxPurchased")
	if profile0Vbucks == nil {
		return fmt.Errorf("profile0vbucks not found")
	}

	if vbucks.Quantity < price {
		return fmt.Errorf("not enough vbucks")
	}

	vbucks.Quantity -= price
	profile0Vbucks.Quantity = price
	vbucks.Save()
	profile0Vbucks.Save()

	for _, receiverAccountId := range body.ReceiverAccountIds {
		receiverPerson := p.Find(receiverAccountId)
		gift := p.NewGift(body.GiftWrapTemplateId, 1, person.ID, body.PersonalMessage)
		for _, grant := range offer.Grants {
			item := p.NewItem(grant.Type.BackendValue + ":" + grant.ID, 1)
			item.ProfileType = offer.Meta.ProfileId
			gift.AddLoot(item)
		}
		
		receiverPerson.CommonCoreProfile.Gifts.AddGift(gift).Save()

		socket, ok := socket.JabberSockets.Get(receiverPerson.ID)
		if ok {
			socket.JabberSendMessageToPerson(aid.JSON{
				"payload": aid.JSON{},
				"type": "com.epicgames.gift.received",
				"timestamp": time.Now().Format("2006-01-02T15:04:05.999Z"),
			})
		}
	}

	return nil
}

func clientRemoveGiftBoxAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	var body struct {
		GiftBoxItemId string `json:"giftBoxItemId" binding:"required"`
	}

	if err := c.BodyParser(&body); err != nil {
		return fmt.Errorf("invalid Body")
	}

	gift := profile.Gifts.GetGift(body.GiftBoxItemId)
	if gift == nil {
		return fmt.Errorf("gift not found")
	}

	for _, item := range gift.Loot {
		person.GetProfileFromType(item.ProfileType).Items.AddItem(item).Save()
	}

	profile.Gifts.DeleteGift(gift.ID)

	return nil
}

func clientSetAffiliateNameAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	var body struct {
		AffiliateName string `json:"affiliateName" binding:"required"`
	}

	if err := c.BodyParser(&body); err != nil {
		return fmt.Errorf("invalid Body")
	}

	affiliate := person.CommonCoreProfile.Attributes.GetAttributeByKey("mtx_affiliate")
	if affiliate == nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("Invalid affiliate attribute"))
	}

	affiliate.ValueJSON = aid.JSONStringify(body.AffiliateName)
	affiliate.Save()

	setTime := person.CommonCoreProfile.Attributes.GetAttributeByKey("mtx_affiliate_set_time")
	if setTime == nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("Invalid affiliate set time attribute"))
	}

	setTime.ValueJSON = aid.JSONStringify(time.Now().Format("2006-01-02T15:04:05.999Z"))
	setTime.Save()

	return nil
}

func clientSetReceiveGiftsEnabledAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	var body struct {
		ReceiveGifts bool `json:"bReceiveGifts" binding:"required"`
	}

	if err := c.BodyParser(&body); err != nil {
		return fmt.Errorf("invalid Body")
	}

	attribute := profile.Attributes.GetAttributeByKey("allowed_to_receive_gifts")
	if attribute == nil {
		return fmt.Errorf("attribute not found")
	}

	attribute.ValueJSON = aid.JSONStringify(body.ReceiveGifts)
	go attribute.Save()

	return nil
}