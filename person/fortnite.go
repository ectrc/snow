package person

import (
	"strconv"

	"github.com/ectrc/snow/aid"
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

func NewFortnitePerson(displayName string, key string) *Person {
	person := NewPerson()
	person.DisplayName = displayName
	person.AccessKey = key

	person.Profile0Profile.Items.AddItem(NewItem("Currency:MtxPurchased", 0)) // for season 2 and bellow

	for _, item := range defaultAthenaItems {
		person.AthenaProfile.Items.AddItem(NewItem(item, 1))
	}

	for _, item := range defaultCommonCoreItems {
		if item == "HomebaseBannerIcon:StandardBanner" {
			for i := 1; i < 32; i++ {
				person.CommonCoreProfile.Items.AddItem(NewItem(item+strconv.Itoa(i), 1))
			}
			continue
		}

		if item == "HomebaseBannerColor:DefaultColor" {
			for i := 1; i < 22; i++ {
				person.CommonCoreProfile.Items.AddItem(NewItem(item+strconv.Itoa(i), 1))
			}
			continue
		}

		if item == "Currency:MtxPurchased" {
			person.CommonCoreProfile.Items.AddItem(NewItem(item, 0))
			continue
		}

		person.CommonCoreProfile.Items.AddItem(NewItem(item, 1))
	}

	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("mfa_reward_claimed", true))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("rested_xp_overflow", 0))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("lifetime_wins", 0))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("party_assist_quest", ""))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("quest_manager", aid.JSON{}))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("inventory_limit_bonus", 0))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("daily_rewards", []aid.JSON{}))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("competitive_identity", aid.JSON{}))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("season_update", 0))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("season_num", aid.Config.Fortnite.Season))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("permissions", []aid.JSON{}))

	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("loadouts", []aid.JSON{}))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("last_applied_loadout", ""))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("active_loadout_index", 0))

	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("accountLevel", 1))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("level", 1))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("xp", 0))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("xp_overflow", 0))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("rested_xp", 0))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("rested_xp_mult", 0))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("rested_xp_exchange", 0))

	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("book_purchased", false))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("book_level", 1))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("book_xp", 0))

	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("favorite_character", person.AthenaProfile.Items.GetItemByTemplateID("AthenaCharacter:CID_001_Athena_Commando_F_Default").ID))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("favorite_backpack", ""))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("favorite_pickaxe", person.AthenaProfile.Items.GetItemByTemplateID("AthenaPickaxe:DefaultPickaxe").ID))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("favorite_glider", person.AthenaProfile.Items.GetItemByTemplateID("AthenaGlider:DefaultGlider").ID))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("favorite_skydivecontrail", ""))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("favorite_dance", make([]string, 6)))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("favorite_itemwraps", make([]string, 7)))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("favorite_loadingscreen", ""))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("favorite_musicpack", ""))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("banner_icon", ""))
	person.AthenaProfile.Attributes.AddAttribute(NewAttribute("banner_color", ""))

	person.CommonCoreProfile.Attributes.AddAttribute(NewAttribute("mfa_enabled", true))
	person.CommonCoreProfile.Attributes.AddAttribute(NewAttribute("mtx_affiliate", ""))
	person.CommonCoreProfile.Attributes.AddAttribute(NewAttribute("mtx_purchase_history", aid.JSON{
		"refundsUsed": 0,
		"refundCredits": 3,
		"purchases": []any{},
	}))
	person.CommonCoreProfile.Attributes.AddAttribute(NewAttribute("current_mtx_platform", "EpicPC"))
	person.CommonCoreProfile.Attributes.AddAttribute(NewAttribute("allowed_to_receive_gifts", true))
	person.CommonCoreProfile.Attributes.AddAttribute(NewAttribute("allowed_to_send_gifts", true))
	person.CommonCoreProfile.Attributes.AddAttribute(NewAttribute("gift_history", aid.JSON{}))
	
	person.Save()

	return person
}