package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/fortnite"
	"github.com/ectrc/snow/person"
	"github.com/ectrc/snow/storage"
)

func informationHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	looker := person.FindByDiscord(i.Member.User.ID)
	if looker == nil {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	if !looker.HasPermission(person.PermissionInformation) {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	playerCount := storage.Repo.GetPersonsCount()
	totalVbucks := storage.Repo.TotalVBucks()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				NewEmbedBuilder().
					SetTitle("Information").
					SetColor(0x2b2d31).
					AddField("Players Registered", aid.FormatNumber(playerCount), true).
					AddField("Players Online", aid.FormatNumber(0), true).
					AddField("VBucks in Circulation", aid.FormatNumber(totalVbucks), false).
					Build(),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
}

func whoHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	looker := person.FindByDiscord(i.Member.User.ID)
	if looker == nil {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	if !looker.HasPermission(person.PermissionLookup) {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	player := getPersonFromOptions(i.ApplicationCommandData().Options, s)
	if player == nil {
		s.InteractionRespond(i.Interaction, &ErrorInvalidDisplayOrDiscord)
		return
	}

	playerVbucks := player.CommonCoreProfile.Items.GetItemByTemplateID("Currency:MtxPurchased")
	if playerVbucks == nil {
		return
	}

	activeCharacter := func() string {
		if player.AthenaProfile == nil {
			return "None"
		}

		characterId := ""
		player.AthenaProfile.Loadouts.RangeLoadouts(func(key string, value *person.Loadout) bool {
			characterId = value.CharacterID
			return false
		})

		if characterId == "" {
			return "None"
		}

		character := player.AthenaProfile.Items.GetItem(characterId)
		if character == nil {
			return "None"
		}

		return character.TemplateID
	}()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				NewEmbedBuilder().
					SetTitle("Player Lookup").
					SetColor(0x2b2d31).
					AddField("Display Name", player.DisplayName, true).
					AddField("VBucks", aid.FormatNumber(playerVbucks.Quantity), true).
					AddField("Discord Account", "<@"+player.Discord.ID+">", true).
					AddField("ID", player.ID, true).
					SetThumbnail("https://fortnite-api.com/images/cosmetics/br/" + strings.Split(activeCharacter, ":")[1] + "/icon.png").
					Build(),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
}

func bansHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	looker := person.FindByDiscord(i.Member.User.ID)
	if looker == nil {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	if !looker.HasPermission(person.PermissionBansControl) {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	if len(i.ApplicationCommandData().Options) <= 0 {
		s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
		return
	}

	subCommand := i.ApplicationCommandData().Options[0]
	if len(subCommand.Options) <= 0 {
		s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
		return
	}

	player := getPersonFromOptions(subCommand.Options, s)
	if player == nil {
		s.InteractionRespond(i.Interaction, &ErrorInvalidDisplayOrDiscord)
		return
	}

	lookup := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, looker *person.Person, player *person.Person, options []*discordgo.ApplicationCommandInteractionDataOption){
		"add": addBanHandler,
		"clear": clearBansHandler,
		"list": listBansHandler,
	}

	if handler, ok := lookup[subCommand.Name]; ok {
		handler(s, i, looker, player, subCommand.Options)
		return
	}

	s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
}

func addBanHandler(s *discordgo.Session, i *discordgo.InteractionCreate, looker *person.Person, player *person.Person, options []*discordgo.ApplicationCommandInteractionDataOption) {
	reason := options[0].StringValue()
	if reason == "" {
		s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
		return
	}

	var expiry string
	for _, option := range options {
		if option.Name == "expires" {
			expiry = option.StringValue()
			break
		}
	}

	player.AddBan(reason, looker.ID, expiry)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: player.DisplayName + " has been banned for `" + reason + "`.",
		},
	})
}

func clearBansHandler(s *discordgo.Session, i *discordgo.InteractionCreate, looker *person.Person, player *person.Person, options []*discordgo.ApplicationCommandInteractionDataOption) {
	player.ClearBans()
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: player.DisplayName + " has had all bans cleared.",
		},
	})
}

func listBansHandler(s *discordgo.Session, i *discordgo.InteractionCreate, looker *person.Person, player *person.Person, options []*discordgo.ApplicationCommandInteractionDataOption) {
	embed := NewEmbedBuilder().
		SetTitle("Ban History").
		SetColor(0x2b2d31)

	player.BanHistory.Range(func(key string, ban *storage.DB_BanStatus) bool {
		banIssuer := person.Find(ban.IssuedBy)
		if banIssuer == nil {
			banIssuer = &person.Person{Discord: &storage.DB_DiscordPerson{ID: "0"}}
		}
		
		embed.AddField(ban.Reason, "Banned by <@"+banIssuer.Discord.ID+">, expires <t:"+fmt.Sprintf("%d", ban.Expiry.Unix())+":D>", false)
		return true
	})

	if player.BanHistory.Len() <= 0 {
		embed.SetDescription("No bans found.")
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.Build()},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
}

func itemsHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	looker := person.FindByDiscord(i.Member.User.ID)
	if looker == nil {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	if !looker.HasPermission(person.PermissionItemControl) {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	if len(i.ApplicationCommandData().Options) <= 0 {
		s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
		return
	}

	subCommand := i.ApplicationCommandData().Options[0]
	if len(subCommand.Options) <= 0 {
		s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
		return
	}

	player := getPersonFromOptions(subCommand.Options, s)
	if player == nil {
		s.InteractionRespond(i.Interaction, &ErrorInvalidDisplayOrDiscord)
		return
	}

	lookup := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, looker *person.Person, player *person.Person, options []*discordgo.ApplicationCommandInteractionDataOption){
		"add": addItemHandler,
		"remove": removeItemHandler,
		"fill": fillItemsHandler,
	}

	if handler, ok := lookup[subCommand.Name]; ok {
		handler(s, i, looker, player, subCommand.Options)
		return
	}

	s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
}

func addItemHandler(s *discordgo.Session, i *discordgo.InteractionCreate, looker *person.Person, player *person.Person, options []*discordgo.ApplicationCommandInteractionDataOption) {
	if !looker.HasPermission(person.PermissionItemControl) {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	item := options[0].StringValue()
	if item == "" {
		s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
		return
	}

	qty := options[1].IntValue()
	if qty <= 0 {
		s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
		return
	}

	profile := options[2].StringValue()
	if profile == "" {
		s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
		return
	}

	if player.GetProfileFromType(profile) == nil {
		s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
		return
	}

	snapshot := player.GetProfileFromType(profile).Snapshot()
	foundItem := player.GetProfileFromType(profile).Items.GetItemByTemplateID(item)
	switch (foundItem) {
	case nil:
		foundItem = person.NewItem(item, int(qty))
		player.GetProfileFromType(profile).Items.AddItem(foundItem)
	default:
		foundItem.Quantity += int(qty)
	}
	foundItem.Save()
	player.GetProfileFromType(profile).Diff(snapshot)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: player.DisplayName + " has been given or updated `" + item + "` in `" + profile + "`.",
		},
	})
}

func removeItemHandler(s *discordgo.Session, i *discordgo.InteractionCreate, looker *person.Person, player *person.Person, options []*discordgo.ApplicationCommandInteractionDataOption) {
	if !looker.HasPermission(person.PermissionItemControl) {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	item := options[0].StringValue()
	if item == "" {
		s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
		return
	}

	qty := options[1].IntValue()
	if qty <= 0 {
		s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
		return
	}

	profile := options[2].StringValue()
	if profile == "" {
		s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
		return
	}

	if player.GetProfileFromType(profile) == nil {
		s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
		return
	}

	snapshot := player.GetProfileFromType(profile).Snapshot()
	foundItem := player.GetProfileFromType(profile).Items.GetItemByTemplateID(item)
	remove := false
	switch (foundItem) {
	case nil:
		s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
		return
	default:
		foundItem.Quantity -= int(qty)
		foundItem.Save()
		if foundItem.Quantity <= 0 {
			player.GetProfileFromType(profile).Items.DeleteItem(foundItem.ID)
			remove = true
		}
	}
	player.GetProfileFromType(profile).Diff(snapshot)

	str := player.DisplayName + " has had `" + aid.FormatNumber(int(qty)) + "` of `" + item + "` removed from `" + profile + "`."
	if remove {
		str = player.DisplayName + " has had `" + item + "` removed from `" + profile + "`."
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: str,
		},
	})
}

func fillItemsHandler(s *discordgo.Session, i *discordgo.InteractionCreate, looker *person.Person, player *person.Person, options []*discordgo.ApplicationCommandInteractionDataOption) {
	if !looker.HasPermission(person.PermissionItemControl) || !looker.HasPermission(person.PermissionLockerControl) {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	fortnite.GiveEverything(player)
	str := player.DisplayName + "has been granted all items."
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &str,
	})
}

func permissionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	looker := person.FindByDiscord(i.Member.User.ID)
	if looker == nil {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	if !looker.HasPermission(person.PermissionPermissionControl) {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	if len(i.ApplicationCommandData().Options) <= 0 {
		s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
		return
	}

	subCommand := i.ApplicationCommandData().Options[0]
	if len(subCommand.Options) <= 0 {
		s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
		return
	}

	player := getPersonFromOptions(subCommand.Options, s)
	if player == nil {
		s.InteractionRespond(i.Interaction, &ErrorInvalidDisplayOrDiscord)
		return
	}

	permission := person.IntToPermission(subCommand.Options[0].IntValue())
	if permission == 0 {
		s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
		return
	}
	
	if permission == person.PermissionAll && !looker.HasPermission(person.PermissionOwner) {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	if player.HasPermission(person.PermissionOwner) && !looker.HasPermission(person.PermissionOwner) {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	if player.HasPermission(person.PermissionAll) && !looker.HasPermission(person.PermissionOwner) {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	switch subCommand.Name {
	case "add":
		player.AddPermission(permission)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: player.DisplayName + " has been given permission `" + permission.GetName() + "`.",
			},
		})
	case "remove":
		player.RemovePermission(permission)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: player.DisplayName + " has had permission `" + permission.GetName() + "` removed.",
			},
		})
	default:
		s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
		return
	}

	s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
}