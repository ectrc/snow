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
			Name: "account",
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

	bansOptions := append([]*discordgo.ApplicationCommandOption{
		{
			Type: discordgo.ApplicationCommandOptionString,
			Name: "reason",
			Description: "The reason for the ban.",
			Required: true,
		},
		{
			Type: discordgo.ApplicationCommandOptionString,
			Name: "expires",
			Description: "The time the ban expires. (e.g. 1y, 1w, 1d, 1h, 1m, 1s) (default: 1w)",
			Required: false,
		},
	}, personOptions...)

	bansSubCommands := []*discordgo.ApplicationCommandOption{
		{
			Type: discordgo.ApplicationCommandOptionSubCommand,
			Name: "add",
			Description: "Ban a player from using this service.",
			Options: bansOptions,
		},
		{
			Type: discordgo.ApplicationCommandOptionSubCommand,
			Name: "clear",
			Description: "Clear all bans from a player.",
			Options: personOptions,
		},
		{
			Type: discordgo.ApplicationCommandOptionSubCommand,
			Name: "list",
			Description: "List all previous bans from a player.",
			Options: personOptions,
		},
	}

	addCommand(&DiscordCommand{
		Command: &discordgo.ApplicationCommand{
			Name: "bans",
			Description: "Perform an action on a player's bans.",
			Options: bansSubCommands,
		},
		Handler: bansHandler,
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
			Description: "The profile to give/take the item from.",
			Required: true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name: "Athena",
					Value: "athena",
				},
				{
					Name: "Common Core",
					Value: "common_core",
				},
				{
					Name: "Profile0",
					Value: "profile0",
				},
				{
					Name: "Creative",
					Value: "creative",
				},
				{
					Name: "Collections",
					Value: "collections",
				},
				{
					Name: "Common Public",
					Value: "common_public",
				},
			},
		},
	}, personOptions...)

	grantSubCommands := []*discordgo.ApplicationCommandOption{
		{
			Type: discordgo.ApplicationCommandOptionSubCommand,
			Name: "add",
			Description: "Grant a player an item in the game.",
			Options: grantOptions,
		},
		{
			Type: discordgo.ApplicationCommandOptionSubCommand,
			Name: "remove",
			Description: "Take an item from a player in the game.",
			Options: grantOptions,
		},
		{
			Type: discordgo.ApplicationCommandOptionSubCommand,
			Name: "fill",
			Description: "Grant a player all items in the game.",
			Options: personOptions,
		},
		{
			Type: discordgo.ApplicationCommandOptionSubCommand,
			Name: "clear",
			Description: "Reset their locker to default.",
			Options: personOptions,
		},
	}

	addCommand(&DiscordCommand{
		Command: &discordgo.ApplicationCommand{
			Name: "items",
			Description: "Perform an action on a player's items.",
			Options: grantSubCommands,
		},
		Handler: itemsHandler,
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
			Description: "Remove a permission from a player.",
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

	rewardOptions := append([]*discordgo.ApplicationCommandOption{
		{
			Type: discordgo.ApplicationCommandOptionString,
			Name: "reward",
			Description: "The reward bundle to give to the player.",
			Required: true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name: "Twitch Prime Drop 1",
					Value: "twitch_prime_1",
				},
				{
					Name: "Twitch Prime Drop 2",
					Value: "twitch_prime_2",
				},
				{
					Name: "Samsung Galaxy",
					Value: "samsung_galaxy",
				},
				{
					Name: "Samsung IKONIK",
					Value: "samsung_ikonik",
				},
				{
					Name: "Honor Guard",
					Value: "honor_guard",
				},
				{
					Name: "Multi Factor Authentication",
					Value: "mfa",
				},
			},
		},
	}, personOptions...)

	addCommand(&DiscordCommand{
		Command: &discordgo.ApplicationCommand{
			Name: "reward",
			Description: "Gift a player a item reward bundle!",
			Options: rewardOptions,
		},
		Handler: rewardHandler,
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