package handlers

import (
	"net/url"

	"github.com/ectrc/snow/aid"
	"github.com/gofiber/fiber/v2"
)

func GetDiscordOAuthURL(c *fiber.Ctx) error {
	return c.Status(200).SendString("https://discord.com/api/oauth2/authorize?client_id="+ aid.Config.Discord.ID +"&redirect_uri="+ url.QueryEscape(aid.Config.API.Host + aid.Config.API.Port +"/snow/discord/callback") + "&response_type=code&scope=identify")
}