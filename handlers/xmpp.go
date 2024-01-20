package handlers

import "github.com/ectrc/snow/aid"

func handlePresenceSocket(id string) {
	aid.Print("(xmpp) connection opened", id)
	socket, ok := sockets.Get(id)
	if !ok {
		return
	}

	for {
		_, msg, err := socket.Connection.ReadMessage()
		if err != nil {
			aid.Print("(xmpp) error reading message", err)
			break
		}
		aid.Print("(xmpp) message received", string(msg))
	}
}