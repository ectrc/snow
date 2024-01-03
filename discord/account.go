package discord

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/fortnite"
	"github.com/ectrc/snow/person"
	"github.com/ectrc/snow/storage"
)

func createHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	modal := &discordgo.InteractionResponseData{
		CustomID: "create://" + i.Member.User.ID,
		Title:    "Create an account",
		Components: []discordgo.MessageComponent{
			&discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "display",
						Label:       "DISPLAY NAME",
						Style:       discordgo.TextInputShort,
						Placeholder: "Enter your crazy display name here!",
						Required:    true,
						MaxLength:   20,
						MinLength:   2,
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
		return
	}

	components, ok := data.Components[0].(*discordgo.ActionsRow)
	if !ok {
		return
	}

	display, ok := components.Components[0].(*discordgo.TextInput)
	if !ok {
		return
	}

	found := person.FindByDiscord(i.Member.User.ID)
	if found != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You already have an account with the display name: `" + found.DisplayName + "`",
			},
		})
		return
	}

	found = person.FindByDisplay(display.Value)
	if found != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Someone already has an account with the display name: `" + found.DisplayName + "`, please choose another one.",
			},
		})
		return
	}

	account := fortnite.NewFortnitePerson(display.Value, false) // or aid.Config.Fortnite.Everything
	discord := &storage.DB_DiscordPerson{
		ID:       i.Member.User.ID,
		PersonID: account.ID,
		Username: i.Member.User.Username,
		Avatar:   i.Member.User.Avatar,
	}
	storage.Repo.SaveDiscordPerson(discord)
	account.Discord = discord
	account.Save()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Your account has been created with the display name: `" + account.DisplayName + "`",
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

func codeHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	player := person.FindByDiscord(i.Member.User.ID)
	if player == nil {
		s.InteractionRespond(i.Interaction, &ErrorNoAccount)
		return
	}

	code := player.ID + "=" + time.Now().Format("2006-01-02T15:04:05.999Z")
	encrypted, sig := aid.KeyPair.EncryptAndSignB64([]byte(code))
	decrypt, err := aid.KeyPair.DecryptAndVerifyB64(encrypted, sig)

	if err || string(decrypt) != code {
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				NewEmbedBuilder().
					SetColor(0x2b2d31).
					SetDescription("`" + encrypted + "`").
					Build(),
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
}