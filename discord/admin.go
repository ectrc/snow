package discord

import (
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

func banHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	looker := person.FindByDiscord(i.Member.User.ID)
	if looker == nil {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	if !looker.HasPermission(person.PermissionBan) {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	player := getPersonFromOptions(i.ApplicationCommandData().Options, s)
	if player == nil {
		s.InteractionRespond(i.Interaction, &ErrorInvalidDisplayOrDiscord)
		return
	}

	player.Ban()
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: player.DisplayName + " has been banned.",
		},
	})
}

func unbanHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	looker := person.FindByDiscord(i.Member.User.ID)
	if looker == nil {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	if !looker.HasPermission(person.PermissionBan) {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	player := getPersonFromOptions(i.ApplicationCommandData().Options, s)
	if player == nil {
		s.InteractionRespond(i.Interaction, &ErrorInvalidDisplayOrDiscord)
		return
	}

	player.Unban()
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: player.DisplayName + " has been unbanned.",
		},
	})
}

func giveItemHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	looker := person.FindByDiscord(i.Member.User.ID)
	if looker == nil {
		s.InteractionRespond(i.Interaction, &ErrorNoAccount)
		return
	}

	if !looker.HasPermission(person.PermissionItemControl) {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	player := getPersonFromOptions(i.ApplicationCommandData().Options, s)
	if player == nil {
		s.InteractionRespond(i.Interaction, &ErrorInvalidDisplayOrDiscord)
		return
	}

	item := i.ApplicationCommandData().Options[0].StringValue()
	if item == "" {
		s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
		return
	}

	qty := i.ApplicationCommandData().Options[1].IntValue()
	if qty <= 0 {
		s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
		return
	}

	profile := i.ApplicationCommandData().Options[2].StringValue()
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

func takeItemHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	looker := person.FindByDiscord(i.Member.User.ID)
	if looker == nil {
		s.InteractionRespond(i.Interaction, &ErrorNoAccount)
		return
	}

	if !looker.HasPermission(person.PermissionItemControl) {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	player := getPersonFromOptions(i.ApplicationCommandData().Options, s)
	if player == nil {
		s.InteractionRespond(i.Interaction, &ErrorInvalidDisplayOrDiscord)
		return
	}

	item := i.ApplicationCommandData().Options[0].StringValue()
	if item == "" {
		s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
		return
	}

	qty := i.ApplicationCommandData().Options[1].IntValue()
	if qty <= 0 {
		s.InteractionRespond(i.Interaction, &ErrorInvalidArguments)
		return
	}

	profile := i.ApplicationCommandData().Options[2].StringValue()
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

func giveEverythingHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	looker := person.FindByDiscord(i.Member.User.ID)
	if looker == nil {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	if !looker.HasPermission(person.PermissionItemControl) || !looker.HasPermission(person.PermissionLockerControl) {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	player := getPersonFromOptions(i.ApplicationCommandData().Options, s)
	if player == nil {
		s.InteractionRespond(i.Interaction, &ErrorInvalidDisplayOrDiscord)
		return
	}
	
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	fortnite.GiveEverything(player)

	str := player.DisplayName + "has been granted everything." 
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