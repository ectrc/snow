package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/ectrc/snow/aid"
	p "github.com/ectrc/snow/person"
)

var (
	socketXmppMessageHandlers = map[string]func(*Socket, *etree.Document) error {
		"open": presenceSocketOpenEvent,
		"iq": presenceSocketIqRootEvent,
		"close": presenceSocketCloseEvent,
	}
)

func handlePresenceSocket(id string) {
	aid.Print("(xmpp) connection opened", id)
	socket, ok := sockets.Get(id)
	if !ok {
		return
	}

	for {
		_, message, err := socket.Connection.ReadMessage()
		if err != nil {
			break
		}

		parsed := etree.NewDocument()
		if err := parsed.ReadFromBytes(message); err != nil {
			return
		}

		aid.Print("(xmpp) message received", string(message))

		if handler, ok := socketXmppMessageHandlers[parsed.Root().Tag]; ok {
			if err := handler(socket, parsed); err != nil {
				aid.Print("(xmpp) connection closed", id, err.Error())
				return
			}
		}
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
	socket.JID = person.ID + "@prod.ol.epicgames.com/" + tree.Root().SelectElement("query").SelectElement("resource").Text()

	document := etree.NewDocument()
	iq := document.CreateElement("iq")
	iq.Attr = append(iq.Attr, etree.Attr{Key: "type", Value: "result"})
	iq.Attr = append(iq.Attr, etree.Attr{Key: "id", Value: "_xmpp_auth1"})
	iq.Attr = append(iq.Attr, etree.Attr{Key: "from", Value: "prod.ol.epicgames.com"})

	socket.WriteTree(document)
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