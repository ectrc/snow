package socket

import (
	"sync"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/person"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type MatchmakerData struct {
	PlaylistID string
	Region string
}

type Socket[T JabberData | MatchmakerData] struct {
	ID string
	Connection *websocket.Conn
	Data *T
	Person *person.Person
	mutex sync.Mutex
}

func (s *Socket[T]) Write(payload []byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.Connection.WriteMessage(websocket.TextMessage, payload)
}

func newSocket[T JabberData | MatchmakerData](conn *websocket.Conn, data ...T) *Socket[T] {
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