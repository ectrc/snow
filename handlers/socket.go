package handlers

import (
	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/person"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SocketType string
var (
	SocketTypeXmpp SocketType = "xmpp"
	SocketTypeUnknown SocketType = "unknown"
)

type Socket struct {
	ID string
	Type SocketType
	Connection *websocket.Conn
	Person *person.Person
}

var (
	handles = map[SocketType]func(string) {
		SocketTypeXmpp: handlePresenceSocket,
	}

	sockets = aid.GenericSyncMap[Socket]{}
)

func MiddlewareWebsocket(c *fiber.Ctx) error {
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	var protocol SocketType
	switch c.Get("Sec-WebSocket-Protocol") {
	case "xmpp":
		protocol = SocketTypeXmpp
	default:
		protocol = SocketTypeUnknown
	}

	c.Locals("uuid", uuid.New().String())
	c.Locals("protocol", protocol)
	return c.Next()
}

func WebsocketConnection(c *websocket.Conn) {
	protocol := c.Locals("protocol").(SocketType)
	uuid := c.Locals("uuid").(string)

	sockets.Add(uuid, &Socket{
		ID: uuid,
		Type: protocol,
		Connection: c,
	})
	defer close(uuid)

	if handle, ok := handles[protocol]; ok {
		handle(uuid)
	} 
}

func close(id string) {
	socket, ok := sockets.Get(id)
	if !ok {
		return
	}
	socket.Connection.Close()
	sockets.Delete(id)
	aid.Print("(xmpp) connection closed", id)
}