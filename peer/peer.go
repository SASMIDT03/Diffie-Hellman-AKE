package peer

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type peer struct {
	name             string
	ownListeningAddr string
	connection       net.Conn
}

func NewPeer(nameOfNewPeer string) *peer {

	listenerAddr := startListener()

	peer := &peer{
		name:             nameOfNewPeer,
		ownListeningAddr: listenerAddr,
	}

	return peer
}

func startListener() string {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		fmt.Println("Error while starting listener")
		log.Fatal(err)
	}

	go waitForIncommingConnections(listener)

	return listener.Addr().String()
}

func waitForIncommingConnections(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error while accepting new connection")
			log.Fatal(err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(connection net.Conn) {
	dec := json.NewDecoder(connection)

	var recivedMsg string

	for {
		err := dec.Decode(&recivedMsg)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(recivedMsg)
	}
}

func (p *peer) Connect(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	p.connection = conn
}

func (p *peer) SendMsg(msg string) {
	enc := json.NewEncoder(p.connection)
	if err := enc.Encode(msg); err != nil {
		fmt.Println("Error while sending: ", msg)
		log.Fatal(err)
	}
}

func (p *peer) GetListenerAddr() string {
	return p.ownListeningAddr
}
