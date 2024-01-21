package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/ectrc/snow/aid"
	p "github.com/ectrc/snow/person"
	"github.com/ectrc/snow/storage"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

type PresenceStatus struct {
	Status          string `json:"Status"`
	IsPlaying       bool   `json:"bIsPlaying"`
	IsJoinable      bool   `json:"bIsJoinable"`
	HasVoiceSupport bool   `json:"bHasVoiceSupport"`
	SessionId       string `json:"SessionId"`
	Properties      map[string]struct {
		SourceId           string `json:"sourceId"`
		SourceDisplayName  string `json:"sourceDisplayName"`
		SourcePlatform     string `json:"sourcePlatform"`
		PartyId            string `json:"partyId"`
		PartyTypeId        int    `json:"partyTypeId"`
		Key                string `json:"key"`
		AppId              string `json:"appId"`
		BuildId            string `json:"buildId"`
		PartyFlags         int    `json:"partyFlags"`
		NotAcceptingReason int    `json:"notAcceptingReason"`
		Pc                 int    `json:"pc"`
	} `json:"Properties"`
}

type PresenceState struct {
	JID string
	Open bool
	RawStatus string
	ParsedStatus PresenceStatus
}

var (
	socketXmppMessageHandlers = map[string]func(*Socket, *etree.Document) error {
		"open": presenceSocketOpenEvent,
		"iq": presenceSocketIqRootEvent,
		"message": presenceSocketMessageEvent,
		"presence": presenceSocketPresenceEvent,
		"close": presenceSocketCloseEvent,
	}
)

func presenceSocketHandle(id string) {
	socket, ok := sockets.Get(id)
	if !ok {
		return
	}

	socket.Type = SocketTypeXmpp
	socket.PresenceState = &PresenceState{}

	for {
		_, message, err := socket.Connection.ReadMessage()
		if err != nil {
			break
		}

		parsed := etree.NewDocument()
		if err := parsed.ReadFromBytes(message); err != nil {
			return
		}

		if handler, ok := socketXmppMessageHandlers[parsed.Root().Tag]; ok {
			if err := handler(socket, parsed); err != nil {
				return
			}
		}
	}

	for _, partial := range storage.Repo.GetFriendsForPerson(socket.Person.ID) {
		friend := socket.Person.GetFriend(partial.ID)
		if friend == nil {
			continue
		}

		friendSocket := GetSocketByPerson(friend.Person)
		if friendSocket == nil {
			continue
		}

		friendDocument := etree.NewDocument()
		friendPresence := friendDocument.CreateElement("presence")
		friendPresence.Attr = append(friendPresence.Attr, etree.Attr{Key: "from", Value: socket.PresenceState.JID})
		friendPresence.Attr = append(friendPresence.Attr, etree.Attr{Key: "to", Value: friendSocket.PresenceState.JID})
		friendPresence.Attr = append(friendPresence.Attr, etree.Attr{Key: "type", Value: "available"})
		friendPresence.Attr = append(friendPresence.Attr, etree.Attr{Key: "xmlns", Value: "jabber:client"})
		friendPresence.CreateElement("status").SetText(aid.JSONStringify(aid.JSON{}))
		friendPresence.CreateElement("show").SetText("away")

		friendSocket.WriteTree(friendDocument)
	}

	for _, party := range socket.PresenceState.ParsedStatus.Properties {
		if party.PartyId == "" {
			continue
		}

		sockets.Range(func(_ string, recieverSocket *Socket) bool {
			if recieverSocket.Type != SocketTypeXmpp || recieverSocket.Person == nil {
				return true
			}

			document := etree.NewDocument()
			message := document.CreateElement("message")
			message.Attr = append(message.Attr, etree.Attr{Key: "id", Value: uuid.New().String()})
			message.Attr = append(message.Attr, etree.Attr{Key: "from", Value: socket.PresenceState.JID})
			message.Attr = append(message.Attr, etree.Attr{Key: "to", Value: recieverSocket.PresenceState.JID})
			message.Attr = append(message.Attr, etree.Attr{Key: "xmlns", Value: "jabber:client"})
			message.CreateElement("body").SetText(aid.JSONStringify(aid.JSON{
				"type": "com.epicgames.party.memberexited",
				"timestamp": time.Now().Format(time.RFC3339),
				"payload": aid.JSON{
					"memberId": socket.Person.ID,
					"partyId": party.PartyId,
					"wasKicked": false,
				},
			}))

			recieverSocket.WriteTree(document)
			return true
		})
	}
}

func presenceSocketOpenEvent(socket *Socket, tree *etree.Document) error {
	document := etree.NewDocument()
	open := document.CreateElement("open")
	open.Attr = append(open.Attr, etree.Attr{Key: "xmlns", Value: "urn:ietf:params:xml:ns:xmpp-framing"})
	open.Attr = append(open.Attr, etree.Attr{Key: "from", Value: "prod.ol.epicgames.com"})
	open.Attr = append(open.Attr, etree.Attr{Key: "version", Value: "1.0"})
	open.Attr = append(open.Attr, etree.Attr{Key: "id", Value: socket.ID})

	socket.WriteTree(document)
	return nil
}

func presenceSocketCloseEvent(socket *Socket, tree *etree.Document) error {
	return fmt.Errorf("safe exit")
}

func presenceSocketIqRootEvent(socket *Socket, tree *etree.Document) error {
	redirect := map[string]func(*Socket, *etree.Document) error{
		"set": presenceSocketIqSetEvent,
		"get": presenceSocketIqGetEvent,
	}

	if handler, ok := redirect[tree.Root().SelectAttr("type").Value]; ok {
		if err := handler(socket, tree); err != nil {
			return err
		}
	}

	return nil
}

func presenceSocketIqSetEvent(socket *Socket, tree *etree.Document) error {
	token := tree.Root().SelectElement("query").SelectElement("password")
	if token == nil || token.Text() == "" {
		return fmt.Errorf("invalid token")
	}
	real := strings.ReplaceAll(token.Text(), "eg1~", "")

	claims, err := aid.JWTVerify(real)
	if err != nil {
		return fmt.Errorf("invalid token")
	}

	if claims["snow_id"] == nil {
		return fmt.Errorf("invalid token")
	}

	snowId, ok := claims["snow_id"].(string)
	if !ok {
		return fmt.Errorf("invalid token")
	}

	person := p.Find(snowId)
	if person == nil {
		return fmt.Errorf("invalid token")
	}

	socket.Person = person
	socket.PresenceState.JID = person.ID + "@prod.ol.epicgames.com/" + tree.Root().SelectElement("query").SelectElement("resource").Text()

	document := etree.NewDocument()
	iq := document.CreateElement("iq")
	iq.Attr = append(iq.Attr, etree.Attr{Key: "type", Value: "result"})
	iq.Attr = append(iq.Attr, etree.Attr{Key: "id", Value: "_xmpp_auth1"})
	iq.Attr = append(iq.Attr, etree.Attr{Key: "from", Value: "prod.ol.epicgames.com"})

	socket.WriteTree(document)
	SendPresenceSocketStatusToFriends(socket)
	return nil
}

func presenceSocketIqGetEvent(socket *Socket, tree *etree.Document) error {
	document := etree.NewDocument()
	iq := document.CreateElement("iq")
	iq.Attr = append(iq.Attr, etree.Attr{Key: "type", Value: "result"})
	iq.Attr = append(iq.Attr, etree.Attr{Key: "id", Value: tree.Root().SelectAttr("id").Value})
	iq.Attr = append(iq.Attr, etree.Attr{Key: "to", Value: tree.Root().SelectAttr("from").Value})
	iq.Attr = append(iq.Attr, etree.Attr{Key: "from", Value: "prod.ol.epicgames.com"})
	ping := iq.CreateElement("ping")
	ping.Attr = append(ping.Attr, etree.Attr{Key: "xmlns", Value: "urn:xmpp:ping"})

	socket.WriteTree(document)
	return nil
}

func presenceSocketMessageEvent(socket *Socket, tree *etree.Document) error {
	reciever := p.Find(strings.Split(tree.Root().SelectAttr("to").Value, "@")[0])
	if reciever == nil {
		return nil
	}

	recieverSocket := GetSocketByPerson(reciever)
	if recieverSocket == nil {
		return nil
	}

	document := etree.NewDocument()
	message := document.CreateElement("message")
	message.Attr = append(message.Attr, etree.Attr{Key: "id", Value: tree.Root().SelectAttr("id").Value})
	message.Attr = append(message.Attr, etree.Attr{Key: "from", Value: socket.PresenceState.JID})
	message.Attr = append(message.Attr, etree.Attr{Key: "to", Value: recieverSocket.PresenceState.JID})
	message.Attr = append(message.Attr, etree.Attr{Key: "xmlns", Value: "jabber:client"})
	message.CreateElement("body").SetText(tree.Root().SelectElement("body").Text())

	recieverSocket.WriteTree(document)
	return nil
}

func presenceSocketPresenceEvent(socket *Socket, tree *etree.Document) error {
	status := tree.Root().SelectElement("status")
	if status == nil {
		return nil
	}

	socket.PresenceState.RawStatus = status.Text()
	json.NewDecoder(strings.NewReader(status.Text())).Decode(&socket.PresenceState.ParsedStatus)

	SendPresenceSocketStatusToFriends(socket)
	return nil
}

func SendPresenceSocketStatusToFriends(socket *Socket) {
	for _, partial := range storage.Repo.GetFriendsForPerson(socket.Person.ID) {
		friend := socket.Person.GetFriend(partial.ID)
		if friend == nil {
			continue
		}

		friendSocket := GetSocketByPerson(friend.Person)
		if friendSocket == nil {
			continue
		}

		friendDocument := etree.NewDocument()
		friendPresence := friendDocument.CreateElement("presence")
		friendPresence.Attr = append(friendPresence.Attr, etree.Attr{Key: "from", Value: socket.PresenceState.JID})
		friendPresence.Attr = append(friendPresence.Attr, etree.Attr{Key: "to", Value: friendSocket.PresenceState.JID})
		friendPresence.Attr = append(friendPresence.Attr, etree.Attr{Key: "type", Value: "available"})
		friendPresence.Attr = append(friendPresence.Attr, etree.Attr{Key: "xmlns", Value: "jabber:client"})
		friendPresence.CreateElement("status").SetText(socket.PresenceState.RawStatus)
		friendSocket.WriteTree(friendDocument)

		document := etree.NewDocument()
		presence := document.CreateElement("presence")
		presence.Attr = append(presence.Attr, etree.Attr{Key: "from", Value: friendSocket.PresenceState.JID})
		presence.Attr = append(presence.Attr, etree.Attr{Key: "to", Value: socket.PresenceState.JID})
		presence.Attr = append(presence.Attr, etree.Attr{Key: "type", Value: "available"})
		presence.Attr = append(presence.Attr, etree.Attr{Key: "xmlns", Value: "jabber:client"})
		presence.CreateElement("status").SetText(friendSocket.PresenceState.RawStatus)
		socket.WriteTree(document)
	}
}

func init() {
	go func() {
		for {
			if aid.Config != nil {
				break
			}
		}

		timer := time.NewTicker(30 * time.Second)
		for {
			<-timer.C
			sockets.Range(func(key string, socket *Socket) bool {
				if socket.Type != SocketTypeXmpp || socket.Person == nil {
					return true
				}

				document := etree.NewDocument()
				iq := document.CreateElement("iq")
				iq.Attr = append(iq.Attr, etree.Attr{Key: "id", Value: "_xmpp_auth1"})
				iq.Attr = append(iq.Attr, etree.Attr{Key: "type", Value: "get"})
				iq.Attr = append(iq.Attr, etree.Attr{Key: "from", Value: "prod.ol.epicgames.com"})
				ping := iq.CreateElement("ping")
				ping.Attr = append(iq.Attr, etree.Attr{Key: "xmlns", Value: "urn:xmpp:ping"})

				socket.WriteTree(document)
				return true
			})
		}
	}()
}