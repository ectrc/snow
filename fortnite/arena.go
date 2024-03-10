package fortnite

import (
	"time"

	"github.com/ectrc/snow/aid"
)

type ArenaScoringRule struct {
	StatName string
	MatchRule string
	RewardTiers []struct{
		Value int
		Points int
		Multiply bool
	}
}

func NewScoringRule(stat, rule string) *ArenaScoringRule {
	return &ArenaScoringRule{
		StatName: stat,
		MatchRule: rule,
		RewardTiers: new(ArenaScoringRule).RewardTiers,
	}
}

func (sr *ArenaScoringRule) AddTier(value, points int, multiply bool) *ArenaScoringRule {
	sr.RewardTiers = append(sr.RewardTiers, struct{
		Value int
		Points int
		Multiply bool
	}{
		Value: value,
		Points: points,
		Multiply: multiply,
	})

	return sr
}

func (sr *ArenaScoringRule) GenerateFortniteScoringRule() aid.JSON {
	tiers := make([]aid.JSON, 0)

	for _, tier := range sr.RewardTiers {
		tiers = append(tiers, aid.JSON{
			"keyValue": tier.Value,
			"pointsEarned": tier.Points,
			"multiplicative": tier.Multiply,
		})
	}

	return aid.JSON{
		"trackedStat": sr.StatName,
		"matchRule": sr.MatchRule,
		"rewardTiers": tiers,
	}
}

type ArenaEventTemplate struct {
	ID string
	MatchLimit int
	PlaylistID string
	ScoringRules []*ArenaScoringRule
}

func NewEventTemplate(id string, limit int) *ArenaEventTemplate {
	return &ArenaEventTemplate{
		ID: id,
		MatchLimit: limit,
		ScoringRules: make([]*ArenaScoringRule, 0),
	}
}

func (et *ArenaEventTemplate) AddScoringRule(rule ...*ArenaScoringRule) {
	et.ScoringRules = append(et.ScoringRules, rule...)
}

func (et *ArenaEventTemplate) GenerateFortniteEventTemplate() aid.JSON {
	rules := make([]aid.JSON, 0)

	for _, rule := range et.ScoringRules {
		rules = append(rules, rule.GenerateFortniteScoringRule())
	}

	return aid.JSON{
		"gameId": "Fortnite",
		"eventTemplateId": et.ID,
		"playlistId": et.PlaylistID,
		"persistentScoreId": "Hype",
		"matchCap": et.MatchLimit,
		"scoringRules": rules,
	}
}

type ArenaEventWindow struct {
	ID string
	ParentEvent *Event
	Template *ArenaEventTemplate
	Round int
	ToBeDetermined bool
	CanLiveSpectate bool
	Meta struct {
		DivisionRank int
		ThresholdToAdvanceDivision int
	}
}

func NewEventWindow(id string, template *ArenaEventTemplate) *ArenaEventWindow {
	return &ArenaEventWindow{
		ID: id,
		Meta: new(ArenaEventWindow).Meta,
		Template: template,
	}
}

func (ew *ArenaEventWindow) GenerateFortniteEventWindow() aid.JSON {
	meta := aid.JSON{
		"divisionRank": ew.Meta.DivisionRank,
		"ThresholdToAdvanceDivision": ew.Meta.ThresholdToAdvanceDivision,
		"RoundType": "Arena",
	}

	allTokens := []string{
		"ARENA_S8_Division1",
		"ARENA_S8_Division2",
		"ARENA_S8_Division3",
		"ARENA_S8_Division4",
		"ARENA_S8_Division5",
		"ARENA_S8_Division6",
		"ARENA_S8_Division7",
	}
	requireAll := []string{}
	requireNone := []string{}

	for index, token := range allTokens {
		if index == ew.Meta.DivisionRank {
			requireAll = append(requireAll, token)
			continue
		}

		requireNone = append(requireNone, token)
	}

	return aid.JSON{
		"eventWindowId": ew.ID,
		"eventTemplateId": ew.Template.ID,
		"countdownBeginTime": "2023-06-15T15:00:00.000Z",
		"beginTime": time.Now().Add(time.Hour * -24).Format(time.RFC3339),
		"endTime": "9999-12-31T23:59:59.000Z",
		"payoutDelay": 30,
		"round": ew.Round,
		"isTBD": ew.ToBeDetermined,
		"canLiveSpectate": ew.CanLiveSpectate,
		"visibility": "public",
		"scoreLocations": []aid.JSON{},
		"blackoutPeriods": []string{},
		"requireAnyTokens": []string{},
		"requireAllTokens": requireAll,
		"requireAllTokensCaller": []string{},
		"requireNoneTokensCaller": requireNone,
		"requireAnyTokensCaller": []string{},
		"additionalRequirements": []string{},
		"teammateEligibility": "any",
		"metadata": meta,
	}
}

type Event struct {
	ID string
	DisplayID string
	Windows []*ArenaEventWindow
}

func NewEvent(id string, displayId string) *Event {
	return &Event{
		ID: id,
		DisplayID: displayId,
		Windows: make([]*ArenaEventWindow, 0),
	}
}

func (e *Event) AddWindow(window *ArenaEventWindow) {
	window.ParentEvent = e
	e.Windows = append(e.Windows, window)
}

func (e *Event) GenerateFortniteEvent() aid.JSON {
	eventWindows := make([]aid.JSON, 0)

	for _, window := range e.Windows {
		eventWindows = append(eventWindows, window.GenerateFortniteEventWindow())
	}

	return aid.JSON{
		"gameId": "Fortnite",
		"eventId": e.ID,
		"eventGroup": "",
		"regions": []string{ "NAE", "ME", "NAW", "OCE", "ASIA", "EU", "BR", },
		"regionMappings": aid.JSON{},
		"platforms": []string{ "PS4", "XboxOne", "Switch", "Android", "IOS", "Windows", },
		"platformMappings": aid.JSON{},
		"displayDataId": e.DisplayID,
		"eventWindows": eventWindows,
		"appId": nil,
		"link": nil,
		"metadata": aid.JSON{
			"minimumAccountLevel": 1,
			"TrackedStats": []string{
				"PLACEMENT_STAT_INDEX",
				"TEAM_ELIMS_STAT_INDEX",
				"MATCH_PLAYED_STAT",
			},
		},
		"environment": nil,
		"announcementTime": time.Now().Format(time.RFC3339),
		"beginTime": time.Now().Add(time.Hour * -24).Format(time.RFC3339),
		"endTime": "9999-12-31T23:59:59.000Z",
	}
}

var (
	ArenaEvents = make([]*Event, 0)
)

func PreloadEvents() {
	if aid.Config.Fortnite.Season < 8 {
		return
	}

	ArenaEvents = []*Event{
		createDuoEvent(),
		createSoloEvent(),
	}
}

func createSoloEvent() *Event {
	ArenaSolo := NewEvent("epicgames_Arena_S8_Solo", "SnowArenaSolo")
	
	defaultPlacement := NewScoringRule("PLACEMENT_STAT_INDEX", "lte")
	defaultPlacement.AddTier(1, 3, false)
	defaultPlacement.AddTier(5, 2, false)
	defaultPlacement.AddTier(15, 2, false)
	defaultPlacement.AddTier(25, 3, false)
	defaultEliminations := NewScoringRule("TEAM_ELIMS_STAT_INDEX", "gte")
	defaultEliminations.AddTier(1, 1, true)

	soloOpen1T := NewEventTemplate("eventTemplate_Arena_S8_Division1_Solo", 100)
	soloOpen1T.PlaylistID = "Playlist_ShowdownAlt_Solo"
	soloOpen1T.AddScoringRule(defaultPlacement, defaultEliminations)
	soloOpen1W := NewEventWindow("Arena_S8_Division1_Solo", soloOpen1T)
	soloOpen1W.ToBeDetermined = false
	soloOpen1W.CanLiveSpectate = false
	soloOpen1W.Round = 0
	soloOpen1W.Meta.DivisionRank = 0
	soloOpen1W.Meta.ThresholdToAdvanceDivision = 25
	ArenaSolo.AddWindow(soloOpen1W)

	soloOpen2T := NewEventTemplate("eventTemplate_Arena_S8_Division2_Solo", 100)
	soloOpen2T.PlaylistID = "Playlist_ShowdownAlt_Solo"
	soloOpen2T.AddScoringRule(defaultPlacement, defaultEliminations)
	soloOpen2W := NewEventWindow("Arena_S8_Division2_Solo", soloOpen2T)
	soloOpen2W.ToBeDetermined = false
	soloOpen2W.CanLiveSpectate = false
	soloOpen2W.Round = 1
	soloOpen2W.Meta.DivisionRank = 1
	soloOpen2W.Meta.ThresholdToAdvanceDivision = 75
	ArenaSolo.AddWindow(soloOpen2W)

	soloOpen3T := NewEventTemplate("eventTemplate_Arena_S8_Division3_Solo", 100)
	soloOpen3T.PlaylistID = "Playlist_ShowdownAlt_Solo"
	soloOpen3T.AddScoringRule(defaultPlacement, defaultEliminations)
	soloOpen3W := NewEventWindow("Arena_S8_Division3_Solo", soloOpen3T)
	soloOpen3W.Round = 2
	soloOpen3W.ToBeDetermined = false
	soloOpen3W.CanLiveSpectate = false
	soloOpen3W.Meta.DivisionRank = 2
	soloOpen3W.Meta.ThresholdToAdvanceDivision = 125
	ArenaSolo.AddWindow(soloOpen3W)

	soloContender4T := NewEventTemplate("eventTemplate_Arena_S8_Division4_Solo", 100)
	soloContender4T.PlaylistID = "Playlist_ShowdownAlt_Solo"
	soloContender4T.AddScoringRule(defaultPlacement, NewScoringRule("MATCH_PLAYED_STAT", "gtw").AddTier(1, -2, false), defaultEliminations)
	soloConteder4W := NewEventWindow("Arena_S8_Division4_Solo", soloContender4T)
	soloConteder4W.Round = 3
	soloConteder4W.ToBeDetermined = false
	soloConteder4W.CanLiveSpectate = false
	soloConteder4W.Meta.DivisionRank = 3
	soloConteder4W.Meta.ThresholdToAdvanceDivision = 175
	ArenaSolo.AddWindow(soloConteder4W)

	soloContender5T := NewEventTemplate("eventTemplate_Arena_S8_Division5_Solo", 100)
	soloContender5T.PlaylistID = "Playlist_ShowdownAlt_Solo"
	soloContender5T.AddScoringRule(defaultPlacement, NewScoringRule("MATCH_PLAYED_STAT", "gtw").AddTier(1, -4, false), defaultEliminations)
	soloConteder5W := NewEventWindow("Arena_S8_Division5_Solo", soloContender5T)
	soloConteder5W.Round = 4
	soloConteder5W.ToBeDetermined = false
	soloConteder5W.CanLiveSpectate = false
	soloConteder5W.Meta.DivisionRank = 4
	soloConteder5W.Meta.ThresholdToAdvanceDivision = 225
	ArenaSolo.AddWindow(soloConteder5W)

	soloContender6T := NewEventTemplate("eventTemplate_Arena_S8_Division6_Solo", 100)
	soloContender6T.PlaylistID = "Playlist_ShowdownAlt_Solo"
	soloContender6T.AddScoringRule(defaultPlacement, NewScoringRule("MATCH_PLAYED_STAT", "gtw").AddTier(1, -6, false), defaultEliminations)
	soloConteder6W := NewEventWindow("Arena_S8_Division6_Solo", soloContender6T)
	soloConteder6W.Round = 5
	soloConteder6W.ToBeDetermined = false
	soloConteder6W.CanLiveSpectate = false
	soloConteder6W.Meta.DivisionRank = 5
	soloConteder6W.Meta.ThresholdToAdvanceDivision = 300
	ArenaSolo.AddWindow(soloConteder6W)

	soloChampions7T := NewEventTemplate("eventTemplate_Arena_S8_Division7_Solo", 100)
	soloChampions7T.PlaylistID = "Playlist_ShowdownAlt_Solo"
	soloChampions7T.AddScoringRule(defaultPlacement, NewScoringRule("MATCH_PLAYED_STAT", "gtw").AddTier(1, -8, false), defaultEliminations)
	soloChampions7W := NewEventWindow("Arena_S8_Division7_Solo", soloChampions7T)
	soloChampions7W.Round = 6
	soloChampions7W.ToBeDetermined = true
	soloChampions7W.CanLiveSpectate = false
	soloChampions7W.Meta.DivisionRank = 6
	soloChampions7W.Meta.ThresholdToAdvanceDivision = 9999999999
	ArenaSolo.AddWindow(soloChampions7W)

	return ArenaSolo
}

func createDuoEvent() *Event {
	ArenaDuo := NewEvent("epicgames_Arena_S8_Duos", "SnowArenaDuos")

	defaultPlacement := NewScoringRule("PLACEMENT_STAT_INDEX", "lte")
	defaultPlacement.AddTier(1, 3, false)
	defaultPlacement.AddTier(3, 2, false)
	defaultPlacement.AddTier(7, 2, false)
	defaultPlacement.AddTier(12, 3, false)
	defaultEliminations := NewScoringRule("TEAM_ELIMS_STAT_INDEX", "gte")
	defaultEliminations.AddTier(1, 1, true)

	duoOpen1T := NewEventTemplate("eventTemplate_Arena_S8_Division1_Duos", 100)
	duoOpen1T.PlaylistID = "Playlist_ShowdownAlt_Duos"
	duoOpen1T.AddScoringRule(defaultPlacement, defaultEliminations)
	duoOpen1W := NewEventWindow("Arena_S8_Division1_Duos", duoOpen1T)
	duoOpen1W.ToBeDetermined = false
	duoOpen1W.CanLiveSpectate = false
	duoOpen1W.Round = 0
	duoOpen1W.Meta.DivisionRank = 0
	duoOpen1W.Meta.ThresholdToAdvanceDivision = 25
	ArenaDuo.AddWindow(duoOpen1W)

	duoOpen2T := NewEventTemplate("eventTemplate_Arena_S8_Division2_Duos", 100)
	duoOpen2T.PlaylistID = "Playlist_ShowdownAlt_Duos"
	duoOpen2T.AddScoringRule(defaultPlacement, defaultEliminations)
	duoOpen2W := NewEventWindow("Arena_S8_Division2_Duos", duoOpen2T)
	duoOpen2W.ToBeDetermined = false
	duoOpen2W.CanLiveSpectate = false
	duoOpen2W.Round = 1
	duoOpen2W.Meta.DivisionRank = 1
	duoOpen2W.Meta.ThresholdToAdvanceDivision = 75
	ArenaDuo.AddWindow(duoOpen2W)

	duoOpen3T := NewEventTemplate("eventTemplate_Arena_S8_Division3_Duos", 100)
	duoOpen3T.PlaylistID = "Playlist_ShowdownAlt_Duos"
	duoOpen3T.AddScoringRule(defaultPlacement, defaultEliminations)
	duoOpen3W := NewEventWindow("Arena_S8_Division3_Duos", duoOpen3T)
	duoOpen3W.Round = 2
	duoOpen3W.ToBeDetermined = false
	duoOpen3W.CanLiveSpectate = false
	duoOpen3W.Meta.DivisionRank = 2
	duoOpen3W.Meta.ThresholdToAdvanceDivision = 125
	ArenaDuo.AddWindow(duoOpen3W)

	duoContender4T := NewEventTemplate("eventTemplate_Arena_S8_Division4_Duos", 100)
	duoContender4T.PlaylistID = "Playlist_ShowdownAlt_Duos"
	duoContender4T.AddScoringRule(defaultPlacement, NewScoringRule("MATCH_PLAYED_STAT", "gtw").AddTier(1, -2, false), defaultEliminations)
	duoConteder4W := NewEventWindow("Arena_S8_Division4_Duos", duoContender4T)
	duoConteder4W.Round = 3
	duoConteder4W.ToBeDetermined = false
	duoConteder4W.CanLiveSpectate = false
	duoConteder4W.Meta.DivisionRank = 3
	duoConteder4W.Meta.ThresholdToAdvanceDivision = 175
	ArenaDuo.AddWindow(duoConteder4W)

	duoContender5T := NewEventTemplate("eventTemplate_Arena_S8_Division5_Duos", 100)
	duoContender5T.PlaylistID = "Playlist_ShowdownAlt_Duos"
	duoContender5T.AddScoringRule(defaultPlacement, NewScoringRule("MATCH_PLAYED_STAT", "gtw").AddTier(1, -4, false), defaultEliminations)
	duoConteder5W := NewEventWindow("Arena_S8_Division5_Duos", duoContender5T)
	duoConteder5W.Round = 4
	duoConteder5W.ToBeDetermined = false
	duoConteder5W.CanLiveSpectate = false
	duoConteder5W.Meta.DivisionRank = 4
	duoConteder5W.Meta.ThresholdToAdvanceDivision = 225
	ArenaDuo.AddWindow(duoConteder5W)

	duoContender6T := NewEventTemplate("eventTemplate_Arena_S8_Division6_Duos", 100)
	duoContender6T.PlaylistID = "Playlist_ShowdownAlt_Duos"
	duoContender6T.AddScoringRule(defaultPlacement, NewScoringRule("MATCH_PLAYED_STAT", "gtw").AddTier(1, -6, false), defaultEliminations)
	duoConteder6W := NewEventWindow("Arena_S8_Division6_Duos", duoContender6T)
	duoConteder6W.Round = 5
	duoConteder6W.ToBeDetermined = false
	duoConteder6W.CanLiveSpectate = false
	duoConteder6W.Meta.DivisionRank = 5
	duoConteder6W.Meta.ThresholdToAdvanceDivision = 300
	ArenaDuo.AddWindow(duoConteder6W)

	duoChampions7T := NewEventTemplate("eventTemplate_Arena_S8_Division7_Duos", 100)
	duoChampions7T.PlaylistID = "Playlist_ShowdownAlt_Duos"
	duoChampions7T.AddScoringRule(defaultPlacement, NewScoringRule("MATCH_PLAYED_STAT", "gtw").AddTier(1, -8, false), defaultEliminations)
	duoChampions7W := NewEventWindow("Arena_S8_Division7_Duos", duoChampions7T)
	duoChampions7W.Round = 6
	duoChampions7W.ToBeDetermined = true
	duoChampions7W.CanLiveSpectate = false
	duoChampions7W.Meta.DivisionRank = 6
	duoChampions7W.Meta.ThresholdToAdvanceDivision = 9999999999
	ArenaDuo.AddWindow(duoChampions7W)

	return ArenaDuo
}