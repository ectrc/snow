package discord

import "github.com/bwmarrin/discordgo"

type EmbedBuilder struct {
	Embed discordgo.MessageEmbed
}

func NewEmbedBuilder() *EmbedBuilder {
	return &EmbedBuilder{
		Embed: discordgo.MessageEmbed{},
	}
}

func (e *EmbedBuilder) SetTitle(title string) *EmbedBuilder {
	e.Embed.Title = title
	return e
}

func (e *EmbedBuilder) SetDescription(description string) *EmbedBuilder {
	e.Embed.Description = description
	return e
}

func (e *EmbedBuilder) SetColor(color int) *EmbedBuilder {
	e.Embed.Color = color
	return e
}

func (e *EmbedBuilder) SetImage(url string) *EmbedBuilder {
	e.Embed.Image = &discordgo.MessageEmbedImage{
		URL: url,
	}
	return e
}

func (e *EmbedBuilder) SetThumbnail(url string) *EmbedBuilder {
	e.Embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: url,
	}
	return e
}

func (e *EmbedBuilder) SetFooter(text string) *EmbedBuilder {
	e.Embed.Footer = &discordgo.MessageEmbedFooter{
		Text: text,
	}
	return e
}

func (e *EmbedBuilder) AddField(name, value string, inline bool) *EmbedBuilder {
	e.Embed.Fields = append(e.Embed.Fields, &discordgo.MessageEmbedField{
		Name: name,
		Value: value,
		Inline: inline,
	})
	return e
}

func (e *EmbedBuilder) Build() *discordgo.MessageEmbed {
	return &e.Embed
}