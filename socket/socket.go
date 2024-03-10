package socket

import (
	"reflect"
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

type WebSocket interface {
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (messageType int, p []byte, err error)
}

type Socket[T JabberData | MatchmakerData] struct {
	ID string
	Connection WebSocket
	Data *T
	Person *person.Person
	M sync.Mutex
}

func (s *Socket[T]) Write(payload []byte) {
	s.M.Lock()
	defer s.M.Unlock()

	err := s.Connection.WriteMessage(websocket.TextMessage, payload)
	if err != nil {
		s.Remove()
	}
}

func (s *Socket[T]) Remove() {
	reflectType := reflect.TypeOf(s.Data).String()
	switch reflectType {
	case "*socket.JabberData":
		JabberSockets.Delete(s.ID)
	case "*socket.MatchmakerData":
		MatchmakerSockets.Delete(s.ID)
	}
}

func newSocket[T JabberData | MatchmakerData](conn WebSocket, data ...T) *Socket[T] {
	additional := data[0]

	return &Socket[T]{
		ID: uuid.New().String(),
		Connection: conn,
		Data: &additional,
	}
}

func NewJabberSocket(conn WebSocket, id string, data JabberData) *Socket[JabberData] {
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