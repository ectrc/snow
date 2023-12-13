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

	for _, command := range StaticClient.Commands {
		StaticClient.RegisterCommand(command)
	}

	err := StaticClient.Client.Open()
	if err != nil {
		panic(err)
	}
}

func (c *DiscordClient) RegisterCommand(command *DiscordCommand) {
	if command.AdminOnly {
		adminDefaultPermission := int64(discordgo.PermissionAdministrator)
		command.Command.DefaultMemberPermissions = &adminDefaultPermission
	}

	_, err := c.Client.ApplicationCommandCreate(aid.Config.Discord.ID, aid.Config.Discord.Guild, command.Command)
	if err != nil {
		aid.Print("Failed to register command: " + command.Command.Name)
		return
	}
}

func (c *DiscordClient) readyHandler(s *discordgo.Session, event *discordgo.Ready) {
	aid.Print("Discord bot is ready")
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