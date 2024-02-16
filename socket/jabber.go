package socket

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/beevik/etree"
	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/person"
	"github.com/gofiber/contrib/websocket"
)

type JabberData struct {
	JabberID string
	PartyID string
	LastPresence string
}

var jabberHandlers = map[string]func(*Socket[JabberData], *etree.Document) error {
	"stream": jabberStreamHandler,
	"open": jabberOpenHandler,
	"iq": jabberIqRootHandler,
	"presence": jabberPresenceRootHandler,
	"message": jabberMessageHandler,
}

func HandleNewJabberSocket(identifier string) {
	aid.Print("new jabber handle: " + identifier)

	socket, ok := JabberSockets.Get(identifier)
	if !ok {
		aid.Print("socket not found", identifier)
		return
	}
	defer JabberSockets.Delete(socket.ID)

	for {
		_, message, failed := socket.Connection.ReadMessage()
		if failed != nil {
			aid.Print("jabber message failed", failed)
			break
		}
		
		JabberSocketOnMessage(socket, message)
	}
}

func JabberSocketOnMessage(socket *Socket[JabberData], message []byte) {
	if strings.Contains(string(message), `">`) {
		message = []byte(strings.ReplaceAll(string(message), `">`, `"/>`))
	}

	aid.Print("jabber message", string(message))
	parsed := etree.NewDocument()
	if err := parsed.ReadFromString(string(message)); err != nil {
		aid.Print("jabber message failed to parse", err)
		return
	}

	if handler, ok := jabberHandlers[parsed.Root().Tag]; ok {
		if err := handler(socket, parsed); err != nil {
			socket.Connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error()))
			return
		}
		return
	}

	aid.Print("jabber message not handled", parsed.Root().Tag)
}

func jabberStreamHandler(socket *Socket[JabberData], parsed *etree.Document) error {
	socket.Write([]byte(`<stream:stream xmlns="jabber:client" xmlns:stream="http://etherx.jabber.org/streams" to="prod.ol.epicgames.com" id="`+ socket.ID +`" version="1.0" xml:lang="en">`))
	socket.Write([]byte(`<stream:features xmlns:stream="http://etherx.jabber.org/streams" />`))
	socket.Write([]byte(`<open xmlns="urn:ietf:params:xml:ns:xmpp-framing" from="prod.ol.epicgames.com" version="1.0" id="`+ socket.ID +`" />`))
	return nil
}

func jabberOpenHandler(socket *Socket[JabberData], parsed *etree.Document) error {
	socket.Write([]byte(`<open xmlns="urn:ietf:params:xml:ns:xmpp-framing" from="prod.ol.epicgames.com" version="1.0" id="`+ socket.ID +`" />`))
	socket.Write([]byte(`<stream:features xmlns:stream="http://etherx.jabber.org/streams" />`))
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

	JabberSockets.ChangeKey(socket.ID, person.ID)
	socket.ID = person.ID
	socket.Person = person
	socket.Data.JabberID = snowId + "@prod.ol.epicgames.com/" + parsed.FindElement("/iq/query/resource").Text()
	socket.Data.PartyID = person.DisplayName + ":" + person.ID + parsed.FindElement("/iq/query/resource").Text()

	socket.Write([]byte(`<iq xmlns="jabber:client" type="result" id="_xmpp_auth1" from="prod.ol.epicgames.com" to="`+ socket.Data.JabberID +`" />`))
	return nil
}

func jabberIqGetHandler(socket *Socket[JabberData], parsed *etree.Document) error {
	socket.Write([]byte(`<iq xmlns="jabber:client" type="result" id="`+ parsed.Root().SelectAttr("id").Value +`" from="prod.ol.epicgames.com" to="`+ socket.Data.JabberID +`" />`))
	socket.JabberNotifyFriends()
	return nil
}

func jabberPresenceRootHandler(socket *Socket[JabberData], parsed *etree.Document) error {
	status := parsed.FindElement("/presence/status")
	if status == nil {
		return jabberPresenceJoinGroupchat(socket, parsed)
	}

	socket.Data.LastPresence = status.Text()
	socket.JabberNotifyFriends()

	return nil
}

func jabberPresenceJoinGroupchat(socket *Socket[JabberData], parsed *etree.Document) error {
	towards := parsed.FindElement("/presence").SelectAttr("to").Value

	partyId := aid.Regex(towards, `Party-(.*?)@`)
	if partyId == nil {
		return nil
	}

	party, ok := person.Parties.Get(*partyId)
	if !ok {
		return nil
	}

	for _, member := range party.Members {
		if member.Person.ID == socket.ID {
			return nil
		}
	
		memberSocket, ok := JabberSockets.Get(member.Person.ID)
		if !ok {
			continue
		}
		memberPartyId := "Party-" + party.ID + "@muc.prod.ol.epicgames.com/" + memberSocket.Data.PartyID
		memberRole := aid.Ternary[string](party.Captain.Person.ID == member.Person.ID, "moderator", "participant")
		memberAffiliation := aid.Ternary[string](party.Captain.Person.ID == member.Person.ID, "owner", "none")

		socket.Write([]byte(`<presence xmlns="jabber:client" from="`+ memberPartyId +`" to="`+ socket.Data.JabberID +`">
			<x xmlns="http://jabber.org/protocol/muc#user">
				<item 
					affiliation="`+ memberAffiliation +`" 
					role="`+ memberRole +`" 
					jid="`+ memberSocket.Data.JabberID +`"
					nick="`+ memberPartyId +`"
				/>
			</x>
		</presence>`))
	}

	socketPartyId := "Party-" + party.ID + "@muc.prod.ol.epicgames.com/" + socket.Data.PartyID
	socketRole := aid.Ternary[string](party.Captain.Person.ID == socket.ID, "moderator", "participant")
	socketAffiliation := aid.Ternary[string](party.Captain.Person.ID == socket.ID, "owner", "none")

	socket.Write([]byte(`<presence xmlns="jabber:client" from="`+ socketPartyId +`" to="`+ socket.Data.JabberID +`">
		<x xmlns="http://jabber.org/protocol/muc#user">
			<item 
				affiliation="`+ socketAffiliation +`" 
				role="`+ socketRole +`" 
				jid="`+ socket.Data.JabberID +`"
				nick="`+ socketPartyId +`"
			/>
			<status code="110"/>
			<status code="100"/>
			<status code="170"/>
		</x>
	</presence>`))

	return nil
}

func jabberMessageHandler(socket *Socket[JabberData], parsed *etree.Document) error {
	message := parsed.FindElement("/message")
	if message.FindElement("/body") == nil {
		return nil
	}

	towards := message.SelectAttr("to")
	if towards == nil {
		return nil
	}

	redirect := map[string]func(*Socket[JabberData], *etree.Document) error {
		"chat": jabberMessageChatHandler,
		"groupchat": jabberMessageGroupchatHandler,
	}

	if handler, ok := redirect[towards.Value]; ok {
		if err := handler(socket, parsed); err != nil {
			return err
		}
	}

	return nil
}

// TODO: IMPLEMENT WHISPERING

func jabberMessageChatHandler(socket *Socket[JabberData], parsed *etree.Document) error {
	return nil
}

func jabberMessageGroupchatHandler(socket *Socket[JabberData], parsed *etree.Document) error {
	return nil
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

	s.Write([]byte(`<message xmlns="jabber:client" from="xmpp-admin@prod.ol.epicgames.com" to="`+ jabberSocket.Data.JabberID +`">
		<body>`+ string(data.ToBytes()) +`</body>
	</message>`))
}

func (s *Socket[T]) JabberNotifyFriends() {
	if reflect.TypeOf(s.Data) != reflect.TypeOf(&JabberData{}) {
		return
	}

	jabberSocket, ok := JabberSockets.Get(s.ID)
	if !ok {
		aid.Print("jabber socket not found even though it should be")
		return
	}

	s.Person.Relationships.Range(func(key string, value *person.Relationship) bool {
		friendSocket, found := JabberSockets.Get(value.From.ID)
		if value.Direction == person.RelationshipOutboundDirection {
			friendSocket, found = JabberSockets.Get(value.Towards.ID)
		}
		if !found {
			return true
		}

		friendSocket.Write([]byte(`<presence xmlns="jabber:client" type="available" from="`+ jabberSocket.Data.JabberID +`" to="`+ friendSocket.Data.JabberID +`">
			<status>`+ jabberSocket.Data.LastPresence +`</status>
		</presence>`))

		jabberSocket.Write([]byte(`<presence xmlns="jabber:client" type="available" from="`+ friendSocket.Data.JabberID +`" to="`+ jabberSocket.Data.JabberID +`">
			<status>`+ friendSocket.Data.LastPresence +`</status>
		</presence>`))

		return true
	})

	jabberSocket.Write([]byte(`<presence xmlns="jabber:client" type="available" from="`+ jabberSocket.Data.JabberID +`" to="`+ jabberSocket.Data.JabberID +`">
		<status>`+ jabberSocket.Data.LastPresence +`</status>
	</presence>`))
}