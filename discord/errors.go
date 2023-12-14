package discord

import "github.com/bwmarrin/discordgo"

var ErrorNoAccount = discordgo.InteractionResponse{
	Type: discordgo.InteractionResponseChannelMessageWithSource,
	Data: &discordgo.InteractionResponseData{
		Content: "You do not have an account with the bot.",
		Flags:   discordgo.MessageFlagsEphemeral,
	},
}

var ErrorNoPermission = discordgo.InteractionResponse{
	Type: discordgo.InteractionResponseChannelMessageWithSource,
	Data: &discordgo.InteractionResponseData{
		Content: "You do not have permission to use this command.",
		Flags:   discordgo.MessageFlagsEphemeral,
	},
}

var ErrorInvalidArguments = discordgo.InteractionResponse{
	Type: discordgo.InteractionResponseChannelMessageWithSource,
	Data: &discordgo.InteractionResponseData{
		Content: "Please provide valid arguments.",
		Flags:   discordgo.MessageFlagsEphemeral,
	},
}

var ErrorInvalidDisplayOrDiscord = discordgo.InteractionResponse{
	Type: discordgo.InteractionResponseChannelMessageWithSource,
	Data: &discordgo.InteractionResponseData{
		Content: "Please provide a valid display name or discord user.",
		Flags:   discordgo.MessageFlagsEphemeral,
	},
}

var ErrorNoDisplay = discordgo.InteractionResponse{
	Type: discordgo.InteractionResponseChannelMessageWithSource,
	Data: &discordgo.InteractionResponseData{
		Content: "No display name found with that name, please try again with a valid name.",
		Flags:   discordgo.MessageFlagsEphemeral,
	},
}

var ErrorNoDiscordUser = discordgo.InteractionResponse{
	Type: discordgo.InteractionResponseChannelMessageWithSource,
	Data: &discordgo.InteractionResponseData{
		Content: "No discord user found with that ID, please try again with a valid ID.",
		Flags:   discordgo.MessageFlagsEphemeral,
	},
}