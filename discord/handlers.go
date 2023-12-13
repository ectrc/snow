package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/fortnite"
	"github.com/ectrc/snow/person"
	"github.com/ectrc/snow/storage"
)

var (
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
			Name: "delete",
			Description: "Delete your account with the bot.",
		},
		Handler: deleteHandler,
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
				Embeds: []*discordgo.MessageEmbed{
					NewEmbedBuilder().SetTitle("Account already exists").SetDescription("You already have an account with the display name: `"+ found.DisplayName +"`").SetColor(0xda373c).Build(),
				},
			},
		})
		return
	}
	
	found = person.FindByDisplay(display.Value)
	if found != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					NewEmbedBuilder().SetTitle("Account already exists").SetDescription("An account with that display name already exists, please try a different name.").SetColor(0xda373c).Build(),
				},
			},
		})
		return
	}

	account := fortnite.NewFortnitePerson(display.Value, aid.RandomString(10), false) // or aid.Config.Fortnite.Everything
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
			Embeds: []*discordgo.MessageEmbed{
				NewEmbedBuilder().SetTitle("Account created").SetDescription("Your account has been created with the display name: `"+ account.DisplayName +"`").SetColor(0x2093dc).Build(),
			},
		},
	})
}

func deleteHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	found := person.FindByDiscord(i.Member.User.ID)
	if found == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					NewEmbedBuilder().SetTitle("Account not found").SetDescription("You don't have an account with the bot.").SetColor(0xda373c).Build(),
				},
			},
		})
		return
	}

	storage.Repo.DeleteDiscordPerson(found.Discord.ID)
	found.Delete()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				NewEmbedBuilder().SetTitle("Account deleted").SetDescription("Your account has been deleted.").SetColor(0x2093dc).Build(),
			},
		},
	})
}

func informationHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	playerCount := storage.Repo.GetPersonsCount()
	totalVbucks := storage.Repo.TotalVBucks()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				NewEmbedBuilder().
					SetTitle("Information").
					SetDescription("Useful information about this server's activity!").
					SetColor(0x2093dc).
					AddField("Players Registered", aid.FormatNumber(playerCount), true).
					AddField("Players Online", aid.FormatNumber(0), true).
					AddField("VBucks in Circulation", aid.FormatNumber(totalVbucks), false).
					Build(),
			},
		},
	})
}