package socket

import (
	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/person"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type SocketType int

var (
	SocketTypeJabber      SocketType = 1
	SocketTypeMatchmaking SocketType = 2
)

type JabberData struct {
	JabberID string
}

type MatchmakerData struct {
	PlaylistID string
	Region string
}

type data interface {
	JabberData | MatchmakerData
}

type Socket[T data] struct {
	ID string
	Connection *websocket.Conn
	Data *T
	Person *person.Person
}

func newSocket[T data](conn *websocket.Conn, data ...T) *Socket[T] {
	additional := data[0]

	return &Socket[T]{
		ID: uuid.New().String(),
		Connection: conn,
		Data: &additional,
	}
}

func NewJabberSocket(conn *websocket.Conn, id string, data JabberData) *Socket[JabberData] {
	socket := newSocket[JabberData](conn, data)
	socket.ID = id
	return socket
}

func NewMatchmakerSocket(conn *websocket.Conn, data MatchmakerData) *Socket[MatchmakerData] {
	return newSocket[MatchmakerData](conn, data)
}

var (
	JabberSockets = aid.GenericSyncMap[Socket[JabberData]]{}
	MatchmakerSockets = aid.GenericSyncMap[Socket[MatchmakerData]]{}
)