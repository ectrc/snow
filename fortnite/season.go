package fortnite

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/person"
)

type SeasonObjectReference struct {
	ObjectName string `json:"ObjectName"`
	ObjectPath string `json:"ObjectPath"`
}

type SeasonAssetReference struct {
	AssetPathName string `json:"AssetPathName"`
	SubPathString string `json:"SubPathString"`
}

type SeasonObjectTag struct {
	TagName string `json:"TagName"`
}

type SeasonRewardGiftBox struct {
	GiftBoxToUse      SeasonAssetReference `json:"GiftBoxToUse"`
	GiftBoxFormatData []any                `json:"GiftBoxFormatData"`
}

type ObjectSeasonReward struct {
	ItemDefinition SeasonAssetReference `json:"ItemDefinition"`
	TemplateID     string               `json:"TemplateId"`
	Quantity       int                  `json:"Quantity"`
	RewardGiftBox  SeasonRewardGiftBox  `json:"RewardGiftBox"`
	IsChaseReward  bool                 `json:"IsChaseReward"`
	RewardType     string               `json:"RewardType"`
}

type SeasonScheduleLevel struct {
	Rewards []ObjectSeasonReward `json:"Rewards"`
}

type SeasonSchedule struct {
	Levels []SeasonScheduleLevel `json:"Levels"`
}

type UnrealSeasonProperties struct {
	SeasonNumber                int                   `json:"SeasonNumber"`
	NumSeasonLevels             int                   `json:"NumSeasonLevels"`
	NumBookLevels               int                   `json:"NumBookLevels"`
	ChallengesVisibility        string                `json:"ChallengesVisibility"`
	SeasonXpCurve               SeasonObjectReference `json:"SeasonXpCurve"`
	BookXpCurve                 SeasonObjectReference `json:"BookXpCurve"`
	SeasonStorefront            string                `json:"SeasonStorefront"`
	FreeTokenItemPrimaryAssetId struct {
		PrimaryAssetType struct {
			Name string `json:"Name"`
		} `json:"PrimaryAssetType"`
		PrimaryAssetName string `json:"PrimaryAssetName"`
	} `json:"FreeTokenItemPrimaryAssetId"`

	BattlePassOfferId       string `json:"BattlePassOfferId"`
	BattlePassBundleOfferId string `json:"BattlePassBundleOfferId"`
	BattlePassLevelOfferID  string `json:"BattlePassLevelOfferId"`

	ChallengeSchedulesAlwaysShown        []SeasonObjectReference `json:"ChallengeSchedulesAlwaysShown"`
	FreeLevelsThatNavigateToBattlePass   []int                   `json:"FreeLevelsThatNavigateToBattlePass"`
	FreeLevelsThatAutoOpenTheAboutScreen []int                   `json:"FreeLevelsThatAutoOpenTheAboutScreen"`
	FreeSeasonItemContentTag             SeasonObjectTag         `json:"FreeSeasonItemContentTag"`
	SeasonFirstWinItemContentTag         SeasonObjectTag         `json:"SeasonFirstWinItemContentTag"`
	SeasonGrantsToEveryoneItemContentTag SeasonObjectTag         `json:"SeasonGrantsToEveryoneItemContentTag"`
	BattlePassPaidItemContentTag         SeasonObjectTag         `json:"BattlePassPaidItemContentTag"`
	BattlePassFreeItemContentTag         SeasonObjectTag         `json:"BattlePassFreeItemContentTag"`

	SeasonXpScheduleFree				 SeasonSchedule         `json:"SeasonXpScheduleFree"`
	BookXpScheduleFree           SeasonSchedule         `json:"BookXpScheduleFree"`
	BookXpSchedulePaid           SeasonSchedule         `json:"BookXpSchedulePaid"`
	SeasonGrantsToEveryone       SeasonScheduleLevel    `json:"SeasonGrantsToEveryone"`
	SeasonFirstWinRewards        SeasonScheduleLevel    `json:"SeasonFirstWinRewards"`
	BattleStarSubstitutionReward SeasonScheduleLevel    `json:"BattleStarSubstitutionReward"`
	ExpiringRewardTypes          []SeasonAssetReference `json:"ExpiringRewardTypes"`
	TokensToRemoveAtSeasonEnd    []SeasonAssetReference `json:"TokensToRemoveAtSeasonEnd"`

	DisplayName struct {
		Key             string `json:"Key"`
		SourceString    string `json:"SourceString"`
		LocalizedString string `json:"LocalizedString"`
	}
}

type UnrealProgressionProperties struct {
	RowStruct SeasonObjectReference `json:"RowStruct"`
}

type UnrealEngineObjectProperties interface {
	UnrealSeasonProperties | UnrealProgressionProperties
}

type UnrealNoRows struct{}
type UnrealProgressionRows struct {
	Level         int `json:"Level"`
	XpToNextLevel int `json:"XpToNextLevel"`
	XpTotal       int `json:"XpTotal"`
}

type UnrealEngineObjectRows interface {
	UnrealNoRows | UnrealProgressionRows
}

type UnrealEngineObject[T UnrealEngineObjectProperties, K UnrealEngineObjectRows] struct {
	Type       string        `json:"Type"`
	Name       string        `json:"Name"`
	Class      string        `json:"Class"`
	Properties T             `json:"Properties"`
	Rows       *map[string]K `json:"Rows"`
}

type SnowSeasonDefinition struct {
	DefaultOfferID     string
	BundleOfferID      string
	TierOfferID        string
	LevelProgression   []*UnrealProgressionRows
	BookProgression    []*UnrealProgressionRows
	TierRewardsPremium [][]*ItemGrant
	TierRewardsFree    [][]*ItemGrant
	LevelRewards			 [][]*ItemGrant
	VictoryRewards     []*ItemGrant
	SeasonTokenRemoval []*ItemGrant
	BookXPReplacements []*ItemGrant
}

func NewSeasonDefinition() *SnowSeasonDefinition {
	return &SnowSeasonDefinition{
		LevelProgression:   make([]*UnrealProgressionRows, 101),
		BookProgression:    make([]*UnrealProgressionRows, 101),
		TierRewardsPremium: make([][]*ItemGrant, 101),
		TierRewardsFree:    make([][]*ItemGrant, 101),
		LevelRewards:       make([][]*ItemGrant, 101),
		SeasonTokenRemoval: make([]*ItemGrant, 0),
		VictoryRewards:     make([]*ItemGrant, 0),
		BookXPReplacements: make([]*ItemGrant, 0),
	}
}

func convertAssetPathToTemplateId(assetPath string) string {
	templateIdParts := make([]string, 2)
	regex := regexp.MustCompile(`\.(.*)`)
	assetPathParts := regex.FindStringSubmatch(assetPath)
	if len(assetPathParts) <= 1 {
		return ""
	}
	templateIdParts[1] = assetPathParts[1]

	switch {
	case strings.Contains(assetPath, "Game/Items/PersistentResource"):
		templateIdParts[0] = "AccountResource"
	case strings.Contains(assetPath, "Game/Items/Currency"):
		templateIdParts[0] = "Currency"
	case strings.Contains(assetPath, "Game/Items/Tokens"):
		templateIdParts[0] = "Token"
	case strings.Contains(assetPath, "Game/Athena/Items/ChallengeBundleSchedules"):
		templateIdParts[0] = "ChallengeBundleSchedule"
	case strings.Contains(assetPath, "Game/Athena/Items/Cosmetics/Pickaxes"):
		templateIdParts[0] = "AthenaPickaxe"
	case strings.Contains(assetPath, "Game/Athena/Items/Cosmetics/Dances"):
		templateIdParts[0] = "AthenaDance"
	case strings.Contains(assetPath, "Game/Athena/Items/Cosmetics/Sprays"):
		templateIdParts[0] = "AthenaDance"
	case strings.Contains(assetPath, "Game/Athena/Items/Cosmetics/Backpacks"):
		templateIdParts[0] = "AthenaBackpack"
	case strings.Contains(assetPath, "Game/Athena/Items/Cosmetics/PetCarriers"):
		templateIdParts[0] = "AthenaBackpack"
	case strings.Contains(assetPath, "Game/Athena/Items/Cosmetics/Pets"):
		templateIdParts[0] = "AthenaBackpack"
	case strings.Contains(assetPath, "Game/Athena/Items/Cosmetics/MusicPacks"):
		templateIdParts[0] = "AthenaMusicPack"
	case strings.Contains(assetPath, "Game/Athena/Items/Cosmetics/Characters"):
		templateIdParts[0] = "AthenaCharacter"
	case strings.Contains(assetPath, "Game/Athena/Items/Cosmetics/Gliders"):
		templateIdParts[0] = "AthenaGlider"
	case strings.Contains(assetPath, "Game/Athena/Items/Cosmetics/LoadingScreens"):
		templateIdParts[0] = "AthenaLoadingScreen"
	case strings.Contains(assetPath, "Game/Athena/Items/Cosmetics/Toys"):
		templateIdParts[0] = "AthenaDance"
	case strings.Contains(assetPath, "Game/Athena/Items/Cosmetics/ItemWraps"):
		templateIdParts[0] = "AthenaItemWrap"
	case strings.Contains(assetPath, "Game/Athena/Items/Cosmetics/Contrails"):
		templateIdParts[0] = "AthenaSkydiveContrail"
	case strings.Contains(assetPath, "Game/Athena/Items/CosmeticVariantTokens"):
		templateIdParts[0] = "CosmeticVariantToken"
	default:
		aid.Print("Unknown asset path:", assetPath)
	}

	return strings.Join(templateIdParts, ":")
}

func (s *SnowSeasonDefinition) GetSeasonLevel(stats *person.SeasonStats) int {
	level := 0

	for i, data := range s.LevelProgression {
		if i == 0 {
			continue
		}

		if stats.SeasonXP < data.XpTotal {
			break
		}

		level = i
	}

	return level
}

func (s *SnowSeasonDefinition) GetRelativeSeasonXP(stats *person.SeasonStats) int {
	level := s.GetSeasonLevel(stats)
	if level == 0 {
		return 0
	}

	return stats.SeasonXP - s.LevelProgression[level].XpTotal
}

func (s *SnowSeasonDefinition) GetBookLevel(stats *person.SeasonStats) int {
	level := 0

	for i, data := range s.BookProgression {
		if i == 0 {
			continue
		}

		level = i
		if stats.BookXP - s.BookProgression[i - 1].XpTotal < data.XpToNextLevel {
			break
		}
	}

	return level
}

func (s *SnowSeasonDefinition) GetRelativeBookXP(stats *person.SeasonStats) int {
	level := s.GetBookLevel(stats)
	if level == 0 {
		return 0
	}

	return stats.BookXP - s.BookProgression[level - 1].XpTotal
}

func (s *SnowSeasonDefinition) GrantUnredeemedBookRewards(p *person.Person, giftBoxId string) bool {
	changed := false
	gift := person.NewGift(fmt.Sprintf("GiftBox:%s", giftBoxId), 1, "", "")

	grantUpTo := s.GetBookLevel(p.CurrentSeasonStats)
	freeClaimedUpTo := aid.JSONParseG[int](p.CommonCoreProfile.Attributes.GetAttributeByKey("season.bookFreeClaimedUpTo").ValueJSON)
	paidClaimedUpTo := aid.JSONParseG[int](p.CommonCoreProfile.Attributes.GetAttributeByKey("season.bookPaidClaimedUpTo").ValueJSON)

	if freeClaimedUpTo >= len(s.TierRewardsFree) - 1 {
		return changed
	}

	if freeClaimedUpTo > grantUpTo {
		freeClaimedUpTo = grantUpTo
	}

	if paidClaimedUpTo > grantUpTo {
		paidClaimedUpTo = grantUpTo
	}

	freeRewards := aid.Flatten[*ItemGrant](s.TierRewardsFree[freeClaimedUpTo+1:grantUpTo+1])
	paidRewards := aid.Flatten[*ItemGrant](s.TierRewardsPremium[paidClaimedUpTo+1:grantUpTo+1])

	rewards := []*ItemGrant{}
	rewards = append(rewards, freeRewards...)
	rewards = append(rewards, aid.Ternary[[]*ItemGrant](p.CurrentSeasonStats.BookPurchased, paidRewards, []*ItemGrant{})...)

	for _, reward := range rewards {
		gift.AddLoot(person.NewItem(reward.TemplateID, reward.Quantity))
	}

	if p.CurrentSeasonStats.BookPurchased {
		p.CommonCoreProfile.Attributes.GetAttributeByKey("season.bookPaidClaimedUpTo").SetValue(grantUpTo).Save()
	}
	p.CommonCoreProfile.Attributes.GetAttributeByKey("season.bookFreeClaimedUpTo").SetValue(grantUpTo).Save()

	if len(gift.Loot) > 0 {
		p.CommonCoreProfile.Gifts.AddGift(gift).Save()
		changed = true
	}

	return changed
}

func (s *SnowSeasonDefinition) GrantUnredeemedLevelRewards(p *person.Person) bool {
	changed := false
	grantUpTo := s.GetSeasonLevel(p.CurrentSeasonStats)
	bookLevel := s.GetBookLevel(p.CurrentSeasonStats)

	levelClaimedUpTo := aid.JSONParseG[int](p.CommonCoreProfile.Attributes.GetAttributeByKey("season.levelClaimedUpTo").ValueJSON)
	if levelClaimedUpTo > grantUpTo {
		levelClaimedUpTo = grantUpTo
	}

	if levelClaimedUpTo >= len(s.LevelRewards) - 1 {
		return changed
	}

	wantedRewards := aid.Flatten[*ItemGrant](s.LevelRewards[levelClaimedUpTo+1:grantUpTo+1])
	replacementRewards := []*ItemGrant{}

	for _, reward := range wantedRewards {
		for _, replacement := range s.BookXPReplacements {
			replacementRewards = append(replacementRewards, aid.Ternary[*ItemGrant](reward.TemplateID == "AccountResource:AthenaBattleStar", replacement, reward))
		}
	}
	
	realRewards := aid.Ternary[[]*ItemGrant](bookLevel > 100, replacementRewards, wantedRewards)
	for _, reward := range realRewards {
		GrantToPerson(p, reward)
	}
	
	p.CommonCoreProfile.Attributes.GetAttributeByKey("season.levelClaimedUpTo").SetValue(grantUpTo).Save()
	
	if len(realRewards) > 0 {
		changed = true
	}

	return changed
}