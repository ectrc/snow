package fortnite

import p "github.com/ectrc/snow/person"

func GrantRewardTwitchPrime1(person *p.Person) {
	gift := p.NewGift("GiftBox:GB_Twitch", 1, "", "")
	gift.AddLoot(p.NewItemWithType("AthenaCharacter:CID_089_Athena_Commando_M_RetroGrey", 1, "athena"))
	gift.AddLoot(p.NewItemWithType("AthenaCharacter:CID_085_Athena_Commando_M_Twitch", 1, "athena"))
	gift.AddLoot(p.NewItemWithType("AthenaBackpack:BID_029_RetroGrey", 1, "athena"))
	gift.AddLoot(p.NewItemWithType("AthenaGlider:Glider_ID_018_Twitch", 1, "athena"))
	gift.AddLoot(p.NewItemWithType("AthenaPickaxe:Pickaxe_ID_039_TacticalBlack", 1, "athena"))
	gift.AddLoot(p.NewItemWithType("AthenaDance:Emoji_VictoryRoyale", 1, "athena"))
	gift.AddLoot(p.NewItemWithType("AthenaDance:Emoji_Wow", 1, "athena"))
	gift.AddLoot(p.NewItemWithType("AthenaDance:Emoji_Bush ", 1, "athena"))
	person.CommonCoreProfile.Gifts.AddGift(gift).Save()
}

func GrantRewardTwitchPrime2(person *p.Person) {
	gift := p.NewGift("GiftBox:GB_Twitch", 1, "", "")
	gift.AddLoot(p.NewItemWithType("AthenaCharacter:CID_114_Athena_Commando_F_TacticalWoodland", 1, "athena"))
	gift.AddLoot(p.NewItemWithType("AthenaBackpack:BID_049_TacticalWoodland", 1, "athena"))
	gift.AddLoot(p.NewItemWithType("AthenaPickaxe:Pickaxe_ID_044_TacticalUrbanHammer", 1, "athena"))
	gift.AddLoot(p.NewItemWithType("AthenaDance:EID_HipHop01", 1, "athena"))
	person.CommonCoreProfile.Gifts.AddGift(gift).Save()
}

func GrantRewardSamsungGalaxy(person *p.Person) {
	gift := p.NewGift("GiftBox:GB_SamsungPromo", 1, "", "")
	gift.AddLoot(p.NewItemWithType("AthenaCharacter:CID_175_Athena_Commando_M_Celestial", 1, "athena"))
	gift.AddLoot(p.NewItemWithType("AthenaBackpack:BID_138_Celestial", 1, "athena"))
	gift.AddLoot(p.NewItemWithType("AthenaGlider:Glider_ID_090_Celestial", 1, "athena"))
	gift.AddLoot(p.NewItemWithType("AthenaPickaxe:Pickaxe_ID_116_Celestial", 1, "athena"))
	person.CommonCoreProfile.Gifts.AddGift(gift).Save()
}

func GrantRewardSamsungIkonic(person *p.Person) {
	gift := p.NewGift("GiftBox:GB_SamsungPromo", 1, "", "")
	gift.AddLoot(p.NewItemWithType("AthenaCharacter:CID_313_Athena_Commando_M_KpopFashion", 1, "athena"))
	gift.AddLoot(p.NewItemWithType("AthenaDance:EID_KPopDance03", 1, "athena"))
	person.CommonCoreProfile.Gifts.AddGift(gift).Save()
}

func GrantRewardHonorGuard(person *p.Person) {
	gift := p.NewGift("GiftBox:GB_HonorPromo", 1, "", "")
	gift.AddLoot(p.NewItemWithType("AthenaCharacter:CID_342_Athena_Commando_M_StreetRacerMetallic", 1, "athena"))
	person.CommonCoreProfile.Gifts.AddGift(gift).Save()
}

func GrantRewardTwoFactor(person *p.Person) {
	gift := p.NewGift("GiftBox:GB_MfaReward", 1, "", "")
	gift.AddLoot(p.NewItemWithType("AthenaDance:EID_BoogieDown", 1, "athena"))
	person.CommonCoreProfile.Gifts.AddGift(gift).Save()
}