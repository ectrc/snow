package handlers

import (
	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/socket"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func MiddlewareWebsocket(c *fiber.Ctx) error {
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	var identifier string
	var protocol string

	switch c.Get("Sec-WebSocket-Protocol") {
	case "xmpp":
		identifier = c.Get("Sec-WebSocket-Key")
		protocol = "jabber"
	default:
		protocol = "matchmaking"
		identifier = uuid.New().String()
	}

	c.Locals("identifier", identifier)
	c.Locals("protocol", protocol)

	return c.Next()
}

func WebsocketConnection(c *websocket.Conn) {
	protocol := c.Locals("protocol").(string)
	identifier := c.Locals("identifier").(string)

	switch protocol {
	case "jabber":
		socket.JabberSockets.Set(identifier, socket.NewJabberSocket(c, identifier, socket.JabberData{}))
		socket.HandleNewJabberSocket(identifier)
	case "matchmaking":
		// socket.MatchmakerSockets.Set(identifier, socket.NewMatchmakerSocket(c, socket.MatchmakerData{}))
	default:
		aid.Print("Invalid protocol: " + protocol)
	}
}