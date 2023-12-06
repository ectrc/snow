package fortnite

import (
	"strconv"
	"strings"

	"github.com/ectrc/snow/aid"
	p "github.com/ectrc/snow/person"
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

func NewFortnitePerson(displayName string, key string) *p.Person {
	person := p.NewPerson()
	person.DisplayName = displayName
	person.AccessKey = key

	for _, item := range defaultAthenaItems {
		item := p.NewItem(item, 1)
		item.HasSeen = true
		person.AthenaProfile.Items.AddItem(item)
	}

	for _, item := range defaultCommonCoreItems {
		if item == "HomebaseBannerIcon:StandardBanner" {
			for i := 1; i < 32; i++ {
				item := p.NewItem(item+strconv.Itoa(i), 1)
				item.HasSeen = true
				person.CommonCoreProfile.Items.AddItem(item).Save()
			}
			continue
		}

		if item == "HomebaseBannerColor:DefaultColor" {
			for i := 1; i < 22; i++ {
				item := p.NewItem(item+strconv.Itoa(i), 1)
				item.HasSeen = true
				person.CommonCoreProfile.Items.AddItem(item).Save()
			}
			continue
		}

		if item == "Currency:MtxPurchased" {
			person.CommonCoreProfile.Items.AddItem(p.NewItem(item, 0)).Save()
			person.Profile0Profile.Items.AddItem(p.NewItem(item, 0)).Save()
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
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("season_update", 0)).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("season_num", aid.Config.Fortnite.Season)).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("permissions", []aid.JSON{})).Save()

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

	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("favorite_character", person.AthenaProfile.Items.GetItemByTemplateID("AthenaCharacter:CID_001_Athena_Commando_F_Default").ID)).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("favorite_backpack", "")).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("favorite_pickaxe", person.AthenaProfile.Items.GetItemByTemplateID("AthenaPickaxe:DefaultPickaxe").ID)).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("favorite_glider", person.AthenaProfile.Items.GetItemByTemplateID("AthenaGlider:DefaultGlider").ID)).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("favorite_skydivecontrail", "")).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("favorite_dance", make([]string, 6))).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("favorite_itemwraps", make([]string, 7))).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("favorite_loadingscreen", "")).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("favorite_musicpack", "")).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("banner_icon", "")).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("banner_color", "")).Save()

	person.CommonCoreProfile.Attributes.AddAttribute(p.NewAttribute("mfa_enabled", true)).Save()
	person.CommonCoreProfile.Attributes.AddAttribute(p.NewAttribute("mtx_affiliate", "")).Save()
	person.CommonCoreProfile.Attributes.AddAttribute(p.NewAttribute("mtx_purchase_history", aid.JSON{
		"refundsUsed": 0,
		"refundCredits": 3,
		"purchases": []any{},
	})).Save()
	person.CommonCoreProfile.Attributes.AddAttribute(p.NewAttribute("current_mtx_platform", "EpicPC")).Save()
	person.CommonCoreProfile.Attributes.AddAttribute(p.NewAttribute("allowed_to_receive_gifts", true)).Save()
	person.CommonCoreProfile.Attributes.AddAttribute(p.NewAttribute("allowed_to_send_gifts", true)).Save()
	person.CommonCoreProfile.Attributes.AddAttribute(p.NewAttribute("gift_history", aid.JSON{})).Save()

	loadout := p.NewLoadout("sandbox_loadout", person.AthenaProfile)
	person.AthenaProfile.Loadouts.AddLoadout(loadout).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("loadouts", []string{loadout.ID})).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("last_applied_loadout", loadout.ID)).Save()
	person.AthenaProfile.Attributes.AddAttribute(p.NewAttribute("active_loadout_index", 0)).Save()

	if aid.Config.Fortnite.Everything {
		for _, item := range Cosmetics.Items {
			if strings.Contains(strings.ToLower(item.ID), "random") {
				continue
			}

			item := p.NewItem(item.Type.BackendValue + ":" + item.ID, 1)
			item.HasSeen = true
			person.AthenaProfile.Items.AddItem(item).Save()
		}
	}
	
	person.Save()

	return person
}