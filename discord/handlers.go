package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ectrc/snow/person"
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
			Name: "code",
			Description: "Generate a one-time use code to link your account.",
		},
		Handler: codeHandler,
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

	addCommand(&DiscordCommand{
		Command: &discordgo.ApplicationCommand{
			Name: "everything",
			Description: "Give a player full locker",
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
		Handler: giveEverythingHandler,
		AdminOnly: true,
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