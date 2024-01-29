package socket

import (
	"fmt"
	"reflect"

	"github.com/beevik/etree"
	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/person"
	"github.com/gofiber/contrib/websocket"
)

type JabberData struct {
	JabberID string
}

var jabberHandlers = map[string]func(*Socket[JabberData], *etree.Document) error {
	"open": jabberOpenHandler,
	"iq": jabberIqRootHandler,
}

func HandleNewJabberSocket(identifier string) {
	socket, ok := JabberSockets.Get(identifier)
	if !ok {
		return
	}
	defer JabberSockets.Delete(identifier)

	for {
		_, message, failed := socket.Connection.ReadMessage()
		if failed != nil {
			break
		}

		aid.Print(string(message))
	
		parsed := etree.NewDocument()
		if err := parsed.ReadFromBytes(message); err != nil {
			return
		}

		if handler, ok := jabberHandlers[parsed.Root().Tag]; ok {
			if err := handler(socket, parsed); err != nil {
				socket.Connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
				return
			}
		}
	}
}

func jabberOpenHandler(socket *Socket[JabberData], parsed *etree.Document) error {
	socket.Connection.WriteMessage(websocket.TextMessage, []byte(`<open xmlns="urn:ietf:params:xml:ns:xmpp-framing" from="prod.ol.epicgames.com" version="1.0" id="`+ socket.ID +`" />`))
	socket.Connection.WriteMessage(websocket.TextMessage, []byte(`<stream:features xmlns:stream="http://etherx.jabber.org/streams" />`))

	return nil
}

func jabberIqRootHandler(socket *Socket[JabberData], parsed *etree.Document) error {
	redirect := map[string]func(*Socket[JabberData], *etree.Document) error {
		"set": jabberIqSetHandler,
		"get": jabberIqGetHandler,
	}

	if handler, ok := redirect[parsed.Root().SelectAttr("type").Value]; ok {
		if err := handler(socket, parsed); err != nil {
			return err
		}
	}

	return nil
}

func jabberIqSetHandler(socket *Socket[JabberData], parsed *etree.Document) error {
	snowId, err := aid.GetSnowFromToken(parsed.FindElement("/iq/query/password").Text())
	if err != nil {
		return err
	}

	person := person.Find(snowId)
	if person == nil {
		return fmt.Errorf("person not found")
	}

	socket.Data.JabberID = snowId + "@prod.ol.epicgames.com/" + parsed.FindElement("/iq/query/resource").Text()
	socket.Person = person

	socket.Connection.WriteMessage(websocket.TextMessage, []byte(`<iq xmlns="jabber:client" type="result" id="_xmpp_auth1" from="prod.ol.epicgames.com" to="`+ socket.Data.JabberID +`" />`))
	return nil
}


func jabberIqGetHandler(socket *Socket[JabberData], parsed *etree.Document) error {
	socket.Connection.WriteMessage(websocket.TextMessage, []byte(`<iq xmlns="jabber:client" type="result" id="`+ parsed.Root().SelectAttr("id").Value +`" from="prod.ol.epicgames.com" to="`+ socket.Data.JabberID +`" />`))
	return nil
}

func GetJabberSocketByPersonID(id string) (*Socket[JabberData], bool) {
	var found *Socket[JabberData]

	JabberSockets.Range(func(key string, socket *Socket[JabberData]) bool {
		if socket.Person.ID == id {
			found = socket
			return false
		}

		return true
	})

	return found, found != nil
}

func (s *Socket[T]) JabberSendMessageToPerson(data aid.JSON) {
	if reflect.TypeOf(s.Data) != reflect.TypeOf(&JabberData{}) {
		return
	}

	jabberSocket, ok := JabberSockets.Get(s.ID)
	if !ok {
		aid.Print("jabber socket not found even though it should be")
		return
	}

	aid.Print(`<message xmlns="jabber:client" from="xmpp-admin@prod.ol.epicgames.com" to="`+ jabberSocket.Data.JabberID +`">
		<body>`+ string(data.ToBytes()) +`</body>
	</message>`)

	s.Connection.WriteMessage(1, []byte(`<message xmlns="jabber:client" from="xmpp-admin@prod.ol.epicgames.com" to="`+ jabberSocket.Data.JabberID +`">
		<body>`+ string(data.ToBytes()) +`</body>
	</message>`))
}