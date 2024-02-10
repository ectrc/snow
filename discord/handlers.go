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

	personOptions := []*discordgo.ApplicationCommandOption{
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
			Description: "Useful information about this server's activity!",
		},
		Handler: informationHandler,
		AdminOnly: true,
	})

	addCommand(&DiscordCommand{
		Command: &discordgo.ApplicationCommand{
			Name: "who",
			Description: "Lookup a player's information.",
			Options: personOptions,
		},
		Handler: whoHandler,
		AdminOnly: true,
	})

	addCommand(&DiscordCommand{
		Command: &discordgo.ApplicationCommand{
			Name: "ban",
			Description: "Ban a player from using the bot.",
			Options: personOptions,
		},
		Handler: banHandler,
		AdminOnly: true,
	})

	addCommand(&DiscordCommand{
		Command: &discordgo.ApplicationCommand{
			Name: "unban",
			Description: "Unban a player from using the bot.",
			Options: personOptions,
		},
		Handler: unbanHandler,
		AdminOnly: true,
	})

	grantOptions := append([]*discordgo.ApplicationCommandOption{
		{
			Type: discordgo.ApplicationCommandOptionString,
			Name: "template_id",
			Description: "The item id of the cosmetic to give/take.",
			Required: true,
		},
		{
			Type: discordgo.ApplicationCommandOptionInteger,
			Name: "quantity",
			Description: "The amount of the item to give/take.",
			Required: true,
		},
		{
			Type: discordgo.ApplicationCommandOptionString,
			Name: "profile",
			Description: "common_core, athena, common_public, profile0, collections, creative",
			Required: true,
		},
	}, personOptions...)

	addCommand(&DiscordCommand{
		Command: &discordgo.ApplicationCommand{
			Name: "give",
			Description: "Grant a player an item in the game.",
			Options: grantOptions,
		},
		Handler: giveItemHandler,
		AdminOnly: true,
	})

	addCommand(&DiscordCommand{
		Command: &discordgo.ApplicationCommand{
			Name: "take",
			Description: "Take an item from a player in the game.",
			Options: grantOptions,
		},
		Handler: takeItemHandler,
		AdminOnly: true,
	})

	addCommand(&DiscordCommand{
		Command: &discordgo.ApplicationCommand{
			Name: "everything",
			Description: "Give a player full locker",
			Options: personOptions,
		},
		Handler: giveEverythingHandler,
		AdminOnly: true,
	})

	permissionOptionChoices := []*discordgo.ApplicationCommandOptionChoice{
		{
			Name: "All",
			Value: person.PermissionAll,
		},
		{
			Name: "Lookup",
			Value: person.PermissionLookup,
		},
		{
			Name: "Information",
			Value: person.PermissionInformation,
		},
		{
			Name: "Donator",
			Value: person.PermissionDonator,
		},
		{
			Name: "ItemControl",
			Value: person.PermissionItemControl,
		},
		{
			Name: "LockerControl",
			Value: person.PermissionLockerControl,
		},
		{
			Name: "Owner",
			Value: person.PermissionOwner,
		},
		{
			Name: "PermissionControl",
			Value: person.PermissionPermissionControl,
		},
	}

	permissionOptions := append([]*discordgo.ApplicationCommandOption{
		{
			Type: discordgo.ApplicationCommandOptionInteger,
			Name: "permission",
			Description: "The permission to add/take.",
			Required: true,
			Choices: permissionOptionChoices,
		},
	}, personOptions...)

	permissionSubCommands := []*discordgo.ApplicationCommandOption{
		{
			Type: discordgo.ApplicationCommandOptionSubCommand,
			Name: "add",
			Description: "Add a permission to a player.",
			Options: permissionOptions,
		},
		{
			Type: discordgo.ApplicationCommandOptionSubCommand,
			Name: "remove",
			Description: "Rake a permission from a player.",
			Options: permissionOptions,
		},
	}

	addCommand(&DiscordCommand{
		Command: &discordgo.ApplicationCommand{
			Name: "permission",
			Description: "Give or take permissions from a player.",
			Options: permissionSubCommands,
		},
		Handler: permissionHandler,
		AdminOnly: true,
	})
}

func getPersonFromOptions(opts []*discordgo.ApplicationCommandInteractionDataOption, s *discordgo.Session) *person.Person {
	if len(opts) <= 0 {
		return nil
	}

	for _, option := range opts {
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