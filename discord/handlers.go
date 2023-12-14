package discord

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/fortnite"
	"github.com/ectrc/snow/person"
	"github.com/ectrc/snow/storage"
)

func addCommand(command *DiscordCommand) {
	StaticClient.Commands[command.Command.Name] = command
}

func addModal(modal *DiscordModal) {
	StaticClient.Modals[modal.ID] = modal
}

func addCommands() {
	if StaticClient == nil {
		panic("StaticClient is nil")
	}

	addCommand(&DiscordCommand{
		Command: &discordgo.ApplicationCommand{
			Name: "create",
			Description: "Create an account with the bot.",
		},
		Handler: createHandler,
	})

	addModal(&DiscordModal{
		ID: "create",
		Handler: createModalHandler,
	})

	addCommand(&DiscordCommand{
		Command: &discordgo.ApplicationCommand{
			Name: "information",
			Description: "Useful information about this server's activity! Admin Only.",
		},
		Handler: informationHandler,
		AdminOnly: true,
	})

	addCommand(&DiscordCommand{
		Command: &discordgo.ApplicationCommand{
			Name: "who",
			Description: "Lookup a player's information.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "display",
					Description: "The display name of the player.",
					Required: false,
				},
				{
					Type: discordgo.ApplicationCommandOptionUser,
					Name: "discord",
					Description: "The discord account of the player.",
					Required: false,
				},
			},
		},
		Handler: whoHandler,
		AdminOnly: true,
	})

	addCommand(&DiscordCommand{
		Command: &discordgo.ApplicationCommand{
			Name: "me",
			Description: "Lookup your own information.",
		},
		Handler: meHandler,
	})

	addCommand(&DiscordCommand{
		Command: &discordgo.ApplicationCommand{
			Name: "delete",
			Description: "Delete your account with the bot.",
		},
		Handler: deleteHandler,
	})

	addCommand(&DiscordCommand{
		Command: &discordgo.ApplicationCommand{
			Name: "ban",
			Description: "Ban a player from using the bot.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionUser,
					Name: "discord",
					Description: "The discord account of the player.",
					Required: false,
				},
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "display",
					Description: "The display name of the player.",
					Required: false,
				},
			},
		},
		Handler: banHandler,
		AdminOnly: true,
	})

	addCommand(&DiscordCommand{
		Command: &discordgo.ApplicationCommand{
			Name: "unban",
			Description: "Unban a player from using the bot.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionUser,
					Name: "discord",
					Description: "The discord account of the player.",
					Required: false,
				},
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "display",
					Description: "The display name of the player.",
					Required: false,
				},
			},
		},
		Handler: unbanHandler,
		AdminOnly: true,
	})

	addCommand(&DiscordCommand{
		Command: &discordgo.ApplicationCommand{
			Name: "give",
			Description: "Grant a player an item in the game.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "template_id",
					Description: "The item id of the cosmetic to give.",
					Required: true,
				},
				{
					Type: discordgo.ApplicationCommandOptionInteger,
					Name: "quantity",
					Description: "The amount of the item to give.",
					Required: true,
				},
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "profile",
					Description: "common_core, athena, common_public, profile0, collections, creative",
					Required: true,
				},
				{
					Type: discordgo.ApplicationCommandOptionUser,
					Name: "discord",
					Description: "The discord account of the player.",
					Required: false,
				},
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "display",
					Description: "The display name of the player.",
					Required: false,
				},
			},
		},
		Handler: giveItemHandler,
		AdminOnly: true,
	})

	addCommand(&DiscordCommand{
		Command: &discordgo.ApplicationCommand{
			Name: "take",
			Description: "Take an item from a player in the game.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "template_id",
					Description: "The item id of the cosmetic to take.",
					Required: true,
				},
				{
					Type: discordgo.ApplicationCommandOptionInteger,
					Name: "quantity",
					Description: "The amount of the item to take.",
					Required: true,
				},
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "profile",
					Description: "common_core, athena, common_public, profile0, collections, creative",
					Required: true,
				},
				{
					Type: discordgo.ApplicationCommandOptionUser,
					Name: "discord",
					Description: "The discord account of the player.",
					Required: false,
				},
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "display",
					Description: "The display name of the player.",
					Required: false,
				},
			},
		},
		Handler: takeItemHandler,
		AdminOnly: true,
	})
}

func createHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	modal := &discordgo.InteractionResponseData{
		CustomID: "create://" + i.Member.User.ID,
		Title: "Create an account",
		Components: []discordgo.MessageComponent{
			&discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID: "display",
						Label: "DISPLAY NAME",
						Style: discordgo.TextInputShort,
						Placeholder: "Enter your crazy display name here!",
						Required: true,
						MaxLength: 20,
						MinLength: 2,
					},
				},
			},
		},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: modal,
	})
}

func createModalHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()
	if len(data.Components) <= 0 {
		aid.Print("No components found")
		return
	}

	components, ok := data.Components[0].(*discordgo.ActionsRow)
	if !ok {
		aid.Print("Failed to assert TextInput")
		return
	}

	display, ok := components.Components[0].(*discordgo.TextInput)
	if !ok {
		aid.Print("Failed to assert TextInput")
		return
	}

	found := person.FindByDiscord(i.Member.User.ID)
	if found != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You already have an account with the display name: `"+ found.DisplayName +"`",
			},
		})
		return
	}
	
	found = person.FindByDisplay(display.Value)
	if found != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Someone already has an account with the display name: `"+ found.DisplayName +"`, please choose another one.",
			},
		})
		return
	}

	account := fortnite.NewFortnitePerson(display.Value, false) // or aid.Config.Fortnite.Everything
	discord := &storage.DB_DiscordPerson{
		ID: i.Member.User.ID,
		PersonID: account.ID,
		Username: i.Member.User.Username,
	}
	storage.Repo.SaveDiscordPerson(discord)
	account.Discord = discord
	account.Save()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Your account has been created with the display name: `"+ account.DisplayName +"`",
		},
	})
}

func deleteHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	found := person.FindByDiscord(i.Member.User.ID)
	if found == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You do not have an account with the bot.",
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	storage.Repo.DeleteDiscordPerson(found.Discord.ID)
	found.Delete()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Your account has been deleted.",
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
}

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

	player := getPersonFromOptions(i.ApplicationCommandData(), s)
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
					AddField("ID",  player.ID, true).
					SetThumbnail("https://fortnite-api.com/images/cosmetics/br/"+ strings.Split(activeCharacter, ":")[1] +"/icon.png").
					Build(),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
}

func meHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	player := person.FindByDiscord(i.Member.User.ID)
	if player == nil {
		s.InteractionRespond(i.Interaction, &ErrorNoAccount)
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
					AddField("ID",  player.ID, true).
					SetThumbnail("https://fortnite-api.com/images/cosmetics/br/"+ strings.Split(activeCharacter, ":")[1] +"/icon.png").
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

	player := getPersonFromOptions(i.ApplicationCommandData(), s)
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

	player := getPersonFromOptions(i.ApplicationCommandData(), s)
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

	if !looker.HasPermission(person.PermissionGiveItem) {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	player := getPersonFromOptions(i.ApplicationCommandData(), s)
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

	if !looker.HasPermission(person.PermissionTakeItem) {
		s.InteractionRespond(i.Interaction, &ErrorNoPermission)
		return
	}

	player := getPersonFromOptions(i.ApplicationCommandData(), s)
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

func getPersonFromOptions(data discordgo.ApplicationCommandInteractionData, s *discordgo.Session) *person.Person {
	options := data.Options

	if len(options) <= 0 {
		return nil
	}

	for _, option := range options {
		switch option.Type {
		case discordgo.ApplicationCommandOptionUser: 
			if option.Name != "discord" {
				continue
			} 
			return person.FindByDiscord(option.UserValue(s).ID)
		case discordgo.ApplicationCommandOptionString:
			if option.Name != "display" {
				continue
			}
			return person.FindByDisplay(option.StringValue())
		}
	}

	return nil
}