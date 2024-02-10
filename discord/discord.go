package discord

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ectrc/snow/aid"
)

type DiscordCommand struct {
	Command *discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
	AdminOnly bool
}

type DiscordModal struct {
	ID string
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

type DiscordClient struct {
	Client *discordgo.Session
	Commands map[string]*DiscordCommand
	Modals map[string]*DiscordModal
}

var StaticClient *DiscordClient

func NewDiscordClient(token string) *DiscordClient {
	client, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	client.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	return &DiscordClient{
		Client: client,
		Commands: make(map[string]*DiscordCommand),
		Modals: make(map[string]*DiscordModal),
	}
}

func IntialiseClient() {
	StaticClient = NewDiscordClient(aid.Config.Discord.Token)
	StaticClient.Client.AddHandler(StaticClient.readyHandler)
	StaticClient.Client.AddHandler(StaticClient.interactionHandler)

	addCommands()

	if len(StaticClient.Commands) != len(StaticClient.GetRegisteredCommands()) {
		// StaticClient.UnregisterCommands()
		StaticClient.RegisterCommands()
	}

	err := StaticClient.Client.Open()
	if err != nil {
		panic(err)
	}
}

func (c *DiscordClient) UnregisterCommands() {
	commands := c.GetRegisteredCommands()
	if commands == nil {
		return
	}

	for _, command := range commands {
		err := c.Client.ApplicationCommandDelete(aid.Config.Discord.ID, aid.Config.Discord.Guild, command.ID)
		if err != nil {
			aid.Print("(discord) failed to delete command: " + command.Name)
		}
	}

	commands = c.GetGlobalRegisteredCommands()

	for _, command := range commands {
		err := c.Client.ApplicationCommandDelete(aid.Config.Discord.ID, "", command.ID)
		if err != nil {
			aid.Print("(discord) failed to delete command: " + command.Name)
		}
	}
}

func (c *DiscordClient) RegisterCommands() {
	adminPermission := int64(discordgo.PermissionAdministrator)

	update := []*discordgo.ApplicationCommand{}
	for _, command := range c.Commands {
		if command.AdminOnly {
			command.Command.DefaultMemberPermissions = &adminPermission
			command.Command.Description += " (admin only)"
		}

		update = append(update, command.Command)
	}

	_, err := c.Client.ApplicationCommandBulkOverwrite(aid.Config.Discord.ID, aid.Config.Discord.Guild, update)
	if err != nil {
		aid.Print("(discord) failed to register commands", err)
		return
	}
}

func (c *DiscordClient) GetRegisteredCommands() []*discordgo.ApplicationCommand {
	commands, err := c.Client.ApplicationCommands(aid.Config.Discord.ID, aid.Config.Discord.Guild)
	if err != nil {
		aid.Print("(discord) failed to get commands")
		return nil
	}

	return commands
}

func (c *DiscordClient) GetGlobalRegisteredCommands() []*discordgo.ApplicationCommand {
	commands, err := c.Client.ApplicationCommands(aid.Config.Discord.ID, "")
	if err != nil {
		aid.Print("(discord) failed to get commands")
		return nil
	}

	return commands
}

func (c *DiscordClient) readyHandler(s *discordgo.Session, event *discordgo.Ready) {
	aid.Print("(discord) bot is ready")
}

func (c *DiscordClient) interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if command, ok := c.Commands[i.ApplicationCommandData().Name]; ok {
			command.Handler(s, i)
		}

	case discordgo.InteractionModalSubmit:
		if modal, ok := c.Modals[strings.Split(i.ModalSubmitData().CustomID, "://")[0]]; ok {
			modal.Handler(s, i)
		}
	}
}