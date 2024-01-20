package handlers

import (
	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/person"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Socket struct {
	ID string
	Connection *websocket.Conn
	Person *person.Person
}

var (
	sockets = aid.GenericSyncMap[Socket]{}
)

func MiddlewareWebsocket(c *fiber.Ctx) error {
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	c.Locals("protocol", c.Get("Sec-WebSocket-Protocol"))
	c.Locals("uuid", uuid.New().String())
	return c.Next()
}

func WebsocketConnection(c *websocket.Conn) {
	protocol := c.Locals("protocol").(string)
	uuid := c.Locals("uuid").(string)

	sockets.Add(uuid, &Socket{
		ID: uuid,
		Connection: c,
	})
	defer close(uuid)

	switch protocol {
	case "xmpp":
		aid.Print("(xmpp) new connection: ", uuid)
	default:
		aid.Print("(unknown) new connection: ", uuid)
	}

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			break
		}
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