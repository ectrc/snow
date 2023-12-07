package handlers

import (
	"strings"

	"github.com/ectrc/snow/fortnite"
	"github.com/gofiber/fiber/v2"
)

func GetPreloadedCosmetics(c *fiber.Ctx) error {
	return c.JSON(fortnite.Cosmetics)
}

func GetPlaylistImage(c *fiber.Ctx) error {
	playlist := c.Params("playlist")
	if playlist == "" {
		return c.SendStatus(404)
	}
	playlist = strings.Split(playlist, ".")[0]

	image, ok := fortnite.PlaylistImages[playlist]
	if !ok {
		return c.SendStatus(404)
	}
	
	c.Set("Content-Type", "image/png")
	return c.Send(image)
}
