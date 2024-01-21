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
const SocketTypeXmpp SocketType = "xmpp"
const	SocketTypeUnknown SocketType = "unknown"

type Socket struct {
	ID string
	Connection *websocket.Conn
	Person *person.Person
	
	Type SocketType
	PresenceState *PresenceState
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
		SocketTypeXmpp: presenceSocketHandle,
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

func GetSocketByPerson(person *person.Person) *Socket {
	var recieverSocket *Socket
	sockets.Range(func(key string, value *Socket) bool {
		if value.Person == nil {
			return true
		}

		if value.Person.ID == person.ID {
			recieverSocket = value
			return false
		}

		return true
	})
	
	return recieverSocket
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
			aid.Print("(socket) writing message to", message.Socket.ID, string(message.Message))
			message.Socket.Connection.WriteMessage(websocket.TextMessage, message.Message)
		}
	}()
}