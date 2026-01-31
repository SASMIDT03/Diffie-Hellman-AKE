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

	var recivedEnvelope Envelope

	for {
		err := dec.Decode(&recivedEnvelope)
		if err != nil {
			log.Fatal(err)
		}

		switch recivedEnvelope.Type {

		case ConnectEnvelope:
			var listenerAddrOfSender string

			if err := json.Unmarshal(recivedEnvelope.Payload, &listenerAddrOfSender); err != nil {
				fmt.Println("Error while reciving connection")
				continue
			}

			fmt.Println("Someone connected to me, who is listening at: ", listenerAddrOfSender)
		}
	}
}

func (p *peer) Connect(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	p.connection = conn

	ownListenerAddrMarshaled, err := json.Marshal(p.ownListeningAddr)
	if err != nil {
		log.Fatal(err)
	}

	envelopeToSend := Envelope{Type: ConnectEnvelope, Payload: ownListenerAddrMarshaled}

	p.SendEnvelope(envelopeToSend)
}

func (p *peer) SendEnvelope(envelope Envelope) {
	enc := json.NewEncoder(p.connection)
	if err := enc.Encode(envelope); err != nil {
		fmt.Println("Error while envelope")
		log.Fatal(err)
	}
}

func (p *peer) GetListenerAddr() string {
	return p.ownListeningAddr
}
