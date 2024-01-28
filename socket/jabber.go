package socket

import "github.com/ectrc/snow/aid"

func HandleNewJabberSocket(identifier string) {
	_, ok := JabberSockets.Get(identifier)
	if !ok {
		return
	}

	aid.Print("New jabber socket: " + identifier)
}