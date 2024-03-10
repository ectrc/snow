package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/fortnite"
	p "github.com/ectrc/snow/person"
	"github.com/ectrc/snow/shop"
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
		"VerifyRealMoneyPurchase": clientVerifyRealMoneyPurchaseAction,
	}

	repeatingActions = []func(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error{
		clientCalculateTierAndLevel,
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

	for _, action := range repeatingActions {
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

	profile.Revision = aid.Ternary[int](c.QueryInt("rvn") == -1, profile.Revision, c.QueryInt("rvn"))+1
	go profile.Save()
	delete(profileSnapshots, profile.Type)

	multiUpdate := []aid.JSON{}
	for key := range profileSnapshots {
		profile := person.GetProfileFromType(key)
		if profile == nil {
			continue
		}
	
		if len(profile.Changes) == 0 {
			continue
		}
		profile.Revision++
		
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
	accountLevel := 0
	person.AllSeasonsStats.Range(func(key string, value *p.SeasonStats) bool {
		accountLevel += fortnite.DataClient.SnowSeason.GetSeasonLevel(value)
		return true
	})

	profile.CreateFullProfileUpdateChange()
	return nil
}

func clientClientQuestLoginAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	// _g := p.NewGift("GiftBox:GB_FortnitemaresChallenges", 1, "", "")
	// _g.AddLoot(p.NewItemWithType("Token:FoundersPackDailyRewardToken", 1, "common_core"))
	// _g.AddLoot(p.NewItemWithType("Token:MysteryToken", 1, "common_core"))
	// _g.AddLoot(p.NewItemWithType("Token:AccountInventoryBonus", 1, "common_core"))
	// _g.AddLoot(p.NewItemWithType("Token:DeniedOrDisabledCosmeticPlaceholderToken", 1, "athena"))
	// _g.AddLoot(p.NewItemWithType("Token:WorldInventoryBonus", 23, "athena"))
	// _g.AddLoot(p.NewItemWithType("Token:CTF_Dom_Key", 1, "common_core"))
	// _g.AddLoot(p.NewItemWithType("FortIngredient:Ingredient_Crystal_ShadowShard", 32, "common_core"))
	// _g.AddLoot(p.NewItemWithType("Ingredient:Ingredient_Crystal_ShadowShard", 4232, "common_core"))
	// _g.AddLoot(p.NewItemWithType("AccountResource:AthenaBattleStar", 3, "common_core"))
	// _g.AddLoot(p.NewItemWithType("AccountResource:AthenaSeasonalXP", 2, "common_core"))
	// _g.AddLoot(p.NewItemWithType("HomebaseNode:QuestReward_BuildingUpgradeLevel2", 5, "common_core"))
	// _g.AddLoot(p.NewItemWithType("FortHomebaseNode:QuestReward_BuildingUpgradeLevel2", 8, "common_core"))
	// _g.AddLoot(p.NewItemWithType("Token:NeighborhoodCurrency", 2, "common_core"))
	// person.CommonCoreProfile.Gifts.AddGift(_g).Save()
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
	
	storefront := shop.GetShop()
	offerRaw, type_ := storefront.GetOfferByID(body.OfferID)
	if offerRaw == nil {
		return fmt.Errorf("offer not found")
	}

	switch type_ {
	case shop.StorefrontCatalogOfferEnumItem:
		offer := offerRaw.(*shop.StorefrontCatalogOfferTypeItem)
		if (offer.Price.FinalPrice * body.PurchaseQuantity) != body.ExpectedTotalPrice {
			return fmt.Errorf("invalid price")
		}
	case shop.StorefrontCatalogOfferEnumBattlePass:
		offer := offerRaw.(*shop.StorefrontCatalogOfferTypeBattlePass)
		if (offer.Price.FinalPrice * body.PurchaseQuantity) != body.ExpectedTotalPrice {
			return fmt.Errorf("invalid price")
		}
	default:
		return fmt.Errorf("invalid offer type")
	}

	purchaseLookup := map[shop.StorefrontCatalogOfferEnum]func(quantity int, offerRaw interface{}, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error{
		shop.StorefrontCatalogOfferEnumItem: clientPurchaseCatalogItemEntryAction,
		shop.StorefrontCatalogOfferEnumBattlePass: clientPurchaseCatalogBattlePassEntryAction,
	}

	if purchaseFunc, ok := purchaseLookup[type_]; ok {
		return purchaseFunc(body.PurchaseQuantity, offerRaw, person, profile, notifications)
	}

	return nil
}

func clientPurchaseCatalogItemEntryAction(quantity int, offerRaw interface{}, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	offer := offerRaw.(*shop.StorefrontCatalogOfferTypeItem)
	for _, grant := range offer.Rewards {
		if grant.ProfileType != shop.ShopGrantProfileTypeAthena {
			return fmt.Errorf("save the world not implemeted yet")
		}
	}
	person.TakeAndSyncVbucks(offer.Price.FinalPrice * quantity)

	loot := []aid.JSON{}
	purchase := p.NewPurchase(offer.OfferID, offer.Price.FinalPrice)

	groupedRewards := map[string]int{}
	for i := 0; i < quantity; i++ {
		for _, grant := range offer.Rewards {
			groupedRewards[grant.TemplateID] += grant.Quantity
		}
	}

	for templateID, quantity := range groupedRewards {
		r, err := fortnite.GrantToPerson(person, fortnite.NewItemGrant(templateID, quantity))
		if err != nil {
			continue
		}

		loot = append(loot, r.GenerateFortniteLootResultEntry()...)
	}

	for _, item := range loot {
		purchaseItem := p.NewItem(item["itemType"].(string), 1)
		purchaseItem.ID = item["itemGuid"].(string)
		purchaseItem.ProfileType = item["itemProfile"].(string)
		purchase.AddLoot(purchaseItem)
	}

	*notifications = append(*notifications, aid.JSON{
		"type": "CatalogPurchase",
		"lootResult": aid.JSON{
			"items": loot,
		},
		"primary": true,
	})
	person.AthenaProfile.Purchases.AddPurchase(purchase).Save()

	affiliate := person.CommonCoreProfile.Attributes.GetAttributeByKey("mtx_affiliate")
	if affiliate == nil {
		return nil
	}

	creator := p.Find(p.AttributeConvert[string](affiliate))
	if creator != nil {
		creator.CommonCoreProfile.Items.GetItemByTemplateID("Currency:MtxPurchased").Quantity += int(float64(offer.Price.FinalPrice) * 0.10)
		creator.Profile0Profile.Items.GetItemByTemplateID("Currency:MtxPurchased").Quantity += int(float64(offer.Price.FinalPrice) * 0.10)
	}

	return nil
}

func clientPurchaseCatalogBattlePassEntryAction(quantity int, offerRaw interface{}, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	offer := offerRaw.(*shop.StorefrontCatalogOfferTypeBattlePass)
	person.TakeAndSyncVbucks(offer.Price.FinalPrice * quantity)

	groupedRewards := map[string]int{}
	for i := 0; i < quantity; i++ {
		for _, grant := range offer.Rewards {
			groupedRewards[grant.TemplateID] += grant.Quantity
		}
	}

	for templateID, quantity := range groupedRewards {
		_, err := fortnite.GrantToPerson(person, fortnite.NewItemGrant(templateID, quantity))
		if err != nil {
			continue
		}
	}

	receipt := p.NewReceipt(offer.GetOfferID(), 0)
	receipt.SetState("OK")
	person.Receipts.AddReceipt(receipt).Save()
	affiliate := person.CommonCoreProfile.Attributes.GetAttributeByKey("mtx_affiliate")
	if affiliate == nil {
		return nil
	}

	creator := p.Find(p.AttributeConvert[string](affiliate))
	if creator != nil {
		creator.CommonCoreProfile.Items.GetItemByTemplateID("Currency:MtxPurchased").Quantity += int(float64(offer.Price.FinalPrice) * 0.10)
		creator.Profile0Profile.Items.GetItemByTemplateID("Currency:MtxPurchased").Quantity += int(float64(offer.Price.FinalPrice) * 0.10)
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

	person.RefundTickets--
	for _, lootItem := range purchase.Loot {
		person.GetProfileFromType(lootItem.ProfileType).Items.DeleteItem(lootItem.ID)
		person.GetProfileFromType(lootItem.ProfileType).CreateItemRemovedChange(lootItem.ID)
	}
	person.GiveAndSyncVbucks(purchase.TotalPaid)

	return nil
}

func clientGiftCatalogEntryAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	var body struct {
		OfferID string `json:"offerId" binding:"required"`
		Currency string `json:"currency" binding:"required"`
		CurrencySubType string `json:"currencySubType" binding:"required"`
		ExpectedTotalPrice int `json:"expectedTotalPrice" binding:"required"`
		GameContext string `json:"gameContext" binding:"required"`
		GiftWrapTemplateId string `json:"giftWrapTemplateId" binding:"required"`
		PersonalMessage string `json:"personalMessage" binding:"required"`
		ReceiverAccountIds []string `json:"receiverAccountIds" binding:"required"`
	}

	if err := c.BodyParser(&body); err != nil {
		return fmt.Errorf("invalid Body")
	}

	storefront := shop.GetShop()
	offerRaw, type_ := storefront.GetOfferByID(body.OfferID)
	if offerRaw == nil {
		return fmt.Errorf("offer not found")
	}
	offer := offerRaw.(*shop.StorefrontCatalogOfferTypeItem)
	if type_ != shop.StorefrontCatalogOfferEnumItem {
		return fmt.Errorf("invalid offer type")
	}

	if offer.Price.FinalPrice != body.ExpectedTotalPrice {
		return fmt.Errorf("invalid price")
	}

	for _, receiverAccountId := range body.ReceiverAccountIds {
		receiverPerson := p.Find(receiverAccountId)
		if receiverPerson == nil {
			return fmt.Errorf("one or more receivers not found")
		}

		for _, grant := range offer.Rewards {
			if receiverPerson.AthenaProfile.Items.GetItemByTemplateID(grant.TemplateID) != nil {
				return fmt.Errorf("one or more receivers has one of the items")
			}
		}
	}

	price := offer.Price.FinalPrice * len(body.ReceiverAccountIds)
	person.TakeAndSyncVbucks(price)

	for _, receiverAccountId := range body.ReceiverAccountIds {
		receiverPerson := p.Find(receiverAccountId)
		gift := p.NewGift(body.GiftWrapTemplateId, 1, person.ID, body.PersonalMessage)
		for _, grant := range offer.Rewards {
			item := p.NewItem(grant.TemplateID, grant.Quantity)
			item.ProfileType = string(grant.ProfileType)
			gift.AddLoot(item)
		}
		
		receiverPerson.CommonCoreProfile.Gifts.AddGift(gift).Save()
		socket.EmitGiftReceived(receiverPerson)
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

	gift := person.CommonCoreProfile.Gifts.GetGift(body.GiftBoxItemId)
	if gift == nil {
		return fmt.Errorf("gift not found")
	}

	loot := []aid.JSON{}
	for _, item := range gift.Loot {
		result, err := fortnite.GrantToPerson(person, fortnite.NewItemGrant(item.TemplateID, item.Quantity))
		if err != nil {
			continue
		}

		loot = append(loot, result.GenerateFortniteLootResultEntry()...)
		item.DeleteLoot()
	}

	person.CommonCoreProfile.Gifts.DeleteGift(gift.ID)

	*notifications = append(*notifications, aid.JSON{
		"type": "CatalogPurchase",
		"lootResult": aid.JSON{
			"items": loot,
		},
		"primary": true,
	})

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

	profile.Attributes.GetAttributeByKey("allowed_to_receive_gifts").SetValue(body.ReceiveGifts).Save()
	return nil
}

func clientVerifyRealMoneyPurchaseAction(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	var body struct {
		AppStore string `json:"appStore" binding:"required"`
		AppStoreId string `json:"appStoreId" binding:"required"`
		PurchaseCorrelationId string `json:"purchaseCorrelationId" binding:"required"`
		ReceiptId string `json:"receiptId" binding:"required"`
		ReceiptInfo string `json:"receiptInfo" binding:"required"`
	}

	if err := c.BodyParser(&body); err != nil {
		return fmt.Errorf("invalid Body")
	}

	receipt := person.Receipts.GetReceipt(body.ReceiptId)
	if receipt == nil {
		return fmt.Errorf("receipt does not exist")
	}

	if receipt.OfferID != body.AppStoreId {
		return fmt.Errorf("receipt does not match offer")
	}

	gift := p.NewGift("GiftBox:GB_MakeGood", 1, "", "Thank you for your purchase!")
	for _, grant := range receipt.Loot {
		item := p.NewItem(grant.TemplateID, grant.Quantity)
		item.ProfileType = grant.ProfileType
		gift.AddLoot(item)
	}	
	
	person.CommonCoreProfile.Gifts.AddGift(gift).Save()
	person.SetInAppPurchasesAttribute()
	person.SyncVBucks("common_core")
	receipt.SetState("OK")
	receipt.Save()
	return nil
}

func clientCalculateTierAndLevel(c *fiber.Ctx, person *p.Person, profile *p.Profile, notifications *[]aid.JSON) error {
	for {
		tierChanged := fortnite.DataClient.SnowSeason.GrantUnredeemedBookRewards(person, "GB_BattlePass")
		levelChanged := fortnite.DataClient.SnowSeason.GrantUnredeemedLevelRewards(person)

		if !tierChanged && !levelChanged {
			break
		}
	}

	person.AthenaProfile.Attributes.GetAttributeByKey("season_num").SetValue(person.CurrentSeasonStats.Season).Save()
	person.AthenaProfile.Attributes.GetAttributeByKey("level").SetValue(fortnite.DataClient.SnowSeason.GetSeasonLevel(person.CurrentSeasonStats)).Save()
	person.AthenaProfile.Attributes.GetAttributeByKey("xp").SetValue(fortnite.DataClient.SnowSeason.GetRelativeSeasonXP(person.CurrentSeasonStats)).Save()
	person.AthenaProfile.Attributes.GetAttributeByKey("book_purchased").SetValue(person.CurrentSeasonStats.BookPurchased).Save()
	person.AthenaProfile.Attributes.GetAttributeByKey("book_level").SetValue(fortnite.DataClient.SnowSeason.GetBookLevel(person.CurrentSeasonStats)).Save()
	person.AthenaProfile.Attributes.GetAttributeByKey("book_xp").SetValue(fortnite.DataClient.SnowSeason.GetRelativeBookXP(person.CurrentSeasonStats)).Save()

	return nil
}