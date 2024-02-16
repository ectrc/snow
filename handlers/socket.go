package handlers

import (
	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/socket"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func MiddlewareWebsocket(c *fiber.Ctx) error {
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	var protocol string

	switch c.Get("Sec-WebSocket-Protocol") {
	case "xmpp":
		protocol = "jabber"
	default:
		protocol = "matchmaking"
	}

	c.Locals("identifier", "ws-"+aid.RandomString(8))
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

func GetSnowConnectedSockets(c *fiber.Ctx) error {
	jabber := aid.JSON{}
	socket.JabberSockets.Range(func(key string, value *socket.Socket[socket.JabberData]) bool {
		jabber[key] = value
		return true
	})

	matchmaking := aid.JSON{}
	socket.MatchmakerSockets.Range(func(key string, value *socket.Socket[socket.MatchmakerData]) bool {
		matchmaking[key] = value
		return true
	})

	return c.Status(200).JSON(aid.JSON{
		"jabber": jabber,
		"matchmaking": matchmaking,
	})
}