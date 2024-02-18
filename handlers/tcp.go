package handlers

import (
	"fmt"
	"net"
	"strconv"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/socket"
)

type tcpClient struct {
	c *net.Conn
	buffer []byte
	jabber *socket.Socket[socket.JabberData]
}

func (t *tcpClient) WriteMessage(messageType int, data []byte) error {
	_, err := (*t.c).Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (t *tcpClient) ReadMessage() (messageType int, p []byte, err error) {
	n, err := (*t.c).Read(t.buffer)
	if err != nil {
		return 0, nil, err
	}

	return 1, t.buffer[:n], nil
}

func (t *tcpClient) loop() {
	defer t.close()

	for {
		_, p, err := t.ReadMessage()
		if err != nil {
			return
		}

		aid.Print("(tcp) received: " + string(p))
		socket.JabberSocketOnMessage(t.jabber, p)
	}
}

func (t *tcpClient) close() error {
	socket.JabberSockets.Delete(t.jabber.ID)
	(*t.c).Close()
	return nil
}

type tcpServer struct {
	ln net.Listener
	port string
	nope chan string
}

func NewServer() (*tcpServer) {
	portNumber, err := strconv.Atoi(aid.Config.API.Port[1:])
	if err != nil {
		return nil
	}
	portNumber++

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", portNumber))
	if err != nil {
		return nil
	}

	return &tcpServer{
		ln: ln, 
		port: fmt.Sprintf(":%d", portNumber),
	}
}

func (t *tcpServer) Listen() error {
	defer t.ln.Close()
	aid.Print("(tcp) listening on " + aid.Config.API.Host + t.port)

	go t.accept()
	<-t.nope

	return nil
}

func (t *tcpServer) accept() {
	for {
		conn, err := t.ln.Accept()
		if err != nil {
			return
		}

		aid.Print("(tcp) new connection from " + conn.RemoteAddr().String())

		go t.handle(conn)
	}
}

func (t *tcpServer) handle(conn net.Conn) {
	tcpClient := &tcpClient{
		c: &conn,
		buffer: make([]byte, 1024),
	}
	tcpClient.jabber = socket.NewJabberSocket(tcpClient, "tcp-"+aid.RandomString(18), socket.JabberData{})
	socket.JabberSockets.Set(tcpClient.jabber.ID, tcpClient.jabber)

	tcpClient.loop()
}