package fortnite

import (
	"strconv"

	"github.com/ectrc/snow/aid"
	p "github.com/ectrc/snow/person"
	"github.com/google/uuid"
)

var (
	defaultAthenaItems = []string{
		"AthenaCharacter:CID_001_Athena_Commando_F_Default",
		"AthenaPickaxe:DefaultPickaxe",
		"AthenaGlider:DefaultGlider",
		"AthenaDance:EID_DanceMoves",
	}
	defaultCommonCoreItems = []string{
		"Currency:MtxPurchased",
		"HomebaseBannerIcon:StandardBanner",
		"HomebaseBannerColor:DefaultColor",
	}
)

func NewFortnitePerson(displayName string, everything bool) *p.Person {
	return NewFortnitePersonWithId(uuid.New().String(), displayName, everything)
}

func GiveEverything(person *p.Person) {
	for _, item := range DataClient.FortniteItems {
		GrantToPerson(person, NewItemGrant(item.Type.BackendValue+":"+item.ID, 1))
	}

	for key := range DataClient.SnowVariantTokens {
		GrantToPerson(person, NewItemGrant("CosmeticVariantToken:"+key, 1))
	}

	person.Save()
}

func NewFortnitePersonWithId(id string, displayName string, everything bool) *p.Person {
	person := p.NewPersonWithCustomID(id)
	person.DisplayName = displayName

	for _, item := range defaultAthenaItems {
		item := p.NewItem(item, 1)
		item.HasSeen = true
		person.AthenaProfile.Items.AddItem(item)
	}

	for _, item := range defaultCommonCoreItems {
		if item == "HomebaseBannerIcon:StandardBanner" {
			for i := 1; i < 32; i++ {
				GrantToPerson(person, NewItemGrant(item+strconv.Itoa(i), 1))
			}
			continue
		}

		if item == "HomebaseBannerColor:DefaultColor" {
			for i := 1; i < 22; i++ {
				GrantToPerson(person, NewItemGrant(item+strconv.Itoa(i), 1))
			}
			continue
		}

		if item == "Currency:MtxPurchased" {
			person.CommonCoreProfile.Items.AddItem(p.NewItem(item, 9999999)).Save()
			person.Profile0Profile.Items.AddItem(p.NewItem(item, 99999999)).Save()
			continue
		}

		person.CommonCoreProfile.Items.AddItem(p.NewItem(item, 1)).Save()
	}

	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("mfa_reward_claimed", true)).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("rested_xp_overflow", 0)).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("lifetime_wins", 0)).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("party_assist_quest", "")).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("quest_manager", aid.JSON{})).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("inventory_limit_bonus", 0)).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("daily_rewards", []aid.JSON{})).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("competitive_identity", aid.JSON{})).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("permissions", []aid.JSON{})).Save()
	
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("season_update", 0)).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("season_num", aid.Config.Fortnite.Season)).Save()
	
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("accountLevel", 1)).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("level", 1)).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("xp", 0)).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("xp_overflow", 0)).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("rested_xp", 0)).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("rested_xp_mult", 0)).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("rested_xp_exchange", 0)).Save()
	
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("book_purchased", false)).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("book_level", 1)).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("book_xp", 0)).Save()

	seasonStats := p.NewSeasonStats(aid.Config.Fortnite.Season)
	seasonStats.PersonID = person.ID
	seasonStats.Save()

	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("favorite_character", person.AthenaProfile.Items.GetItemByTemplateID("AthenaCharacter:CID_001_Athena_Commando_F_Default").ID)).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("favorite_backpack", "")).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("favorite_pickaxe", person.AthenaProfile.Items.GetItemByTemplateID("AthenaPickaxe:DefaultPickaxe").ID)).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("favorite_glider", person.AthenaProfile.Items.GetItemByTemplateID("AthenaGlider:DefaultGlider").ID)).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("favorite_skydivecontrail", "")).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("favorite_dance", make([]string, 6))).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("favorite_itemwraps", make([]string, 7))).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("favorite_loadingscreen", "")).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("favorite_musicpack", "")).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("banner_icon", "StandardBanner1")).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("banner_color", "DefaultColor1")).Save()

	person.CommonCoreProfile.Attributes.AddAttribute(p.NewAttribute("mfa_enabled", true)).Save()
	person.CommonCoreProfile.Attributes.AddAttribute(p.NewAttribute("mtx_affiliate", "")).Save()
	person.CommonCoreProfile.Attributes.AddAttribute(p.NewAttribute("mtx_affiliate_set_time", 0)).Save()
	person.CommonCoreProfile.Attributes.AddAttribute(p.NewAttribute("mtx_purchase_history", aid.JSON{
		"refundsUsed": 0,
		"refundCredits": 3,
		"purchases": []any{},
	})).Save()
	person.CommonCoreProfile.Attributes.AddAttribute(p.NewAttribute("current_mtx_platform", "EpicPC")).Save()
	person.CommonCoreProfile.Attributes.AddAttribute(p.NewAttribute("allowed_to_receive_gifts", true)).Save()
	person.CommonCoreProfile.Attributes.AddAttribute(p.NewAttribute("allowed_to_send_gifts", true)).Save()
	person.CommonCoreProfile.Attributes.AddAttribute(p.NewAttribute("gift_history", aid.JSON{})).Save()
	person.CommonCoreProfile.Attributes.AddAttribute(p.NewAttribute("in_app_purchases", aid.JSON{
		"receipts": []string{},
		"ignoredReceipts": []string{},
		"fulfillmentCounts": map[string]int{},
		"refreshTimers": aid.JSON{},
	})).Save()

	person.CommonCoreProfile.Attributes.AddAttribute(p.NewAttribute("party.recieveIntents", "ALL")).Save()
	person.CommonCoreProfile.Attributes.AddAttribute(p.NewAttribute("party.recieveInvites", "ALL")).Save()

	person.CommonCoreProfile.Attributes.AddAttribute(p.NewAttribute("season.bookFreeClaimedUpTo", 0)).Save()
	person.CommonCoreProfile.Attributes.AddAttribute(p.NewAttribute("season.bookPaidClaimedUpTo", 0)).Save()
	person.CommonCoreProfile.Attributes.AddAttribute(p.NewAttribute("season.levelClaimedUpTo", 0)).Save()

	loadout := p.NewLoadout("PRESET 1", person.AthenaProfile)
	person.AthenaProfile.Loadouts.AddLoadout(loadout).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("loadouts", []string{loadout.ID})).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("last_applied_loadout", loadout.ID)).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("active_loadout_index", 0)).Save()

	if everything {
		GiveEverything(person)
	}
	
	person.Save()

	return person
}