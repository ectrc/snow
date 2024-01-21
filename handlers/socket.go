package handlers

import (
	"github.com/beevik/etree"
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
	JID string
	Type SocketType
	Connection *websocket.Conn
	Person *person.Person
}

type MessageToWrite struct {
	Socket *Socket
	Message []byte
}

func (s *Socket) Write(message []byte) {
	socketWriteQueue <- MessageToWrite{
		Socket: s,
		Message: message,
	}
}

func (s *Socket) WriteTree(message *etree.Document) {
	bytes, err := message.WriteToBytes()
	if err != nil {
		return
	}

	socketWriteQueue <- MessageToWrite{
		Socket: s,
		Message: bytes,
	}
}

var (
	socketWriteQueue = make(chan MessageToWrite, 1000)
	socketHandlers = map[SocketType]func(string) {
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
	defer func() {
		socket, ok := sockets.Get(uuid)
		if !ok {
			return
		}
		socket.Connection.Close()
		sockets.Delete(uuid)
		aid.Print("(xmpp) connection closed", uuid)
	}()

	if handle, ok := socketHandlers[protocol]; ok {
		handle(uuid)
	} 
}

func init() {
	go func() {
		for {
			if aid.Config != nil {
				break
			}
		}

		aid.Print("(socket) write queue started")

		for {
			message := <-socketWriteQueue
			aid.Print("(socket) message sent", message.Socket.ID, string(message.Message))
			message.Socket.Connection.WriteMessage(websocket.TextMessage, message.Message)
		}
	}()
}