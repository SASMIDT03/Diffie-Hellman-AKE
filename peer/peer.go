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
	connectionAddr   string
}

func NewPeer(nameOfNewPeer string, port string) *peer {

	peer := &peer{
		name: nameOfNewPeer,
	}

	peer.ownListeningAddr = peer.startListener(port)

	return peer
}

func (p *peer) startListener(port string) string {
	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println("Error while starting listener")
		log.Fatal(err)
	}

	go p.waitForIncommingConnections(listener)

	return listener.Addr().String()
}

func (p *peer) waitForIncommingConnections(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error while accepting new connection")
			log.Fatal(err)
		}
		go p.handleConnection(conn)
	}
}

func (p *peer) handleConnection(connection net.Conn) {
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

			if listenerAddrOfSender == p.connectionAddr {
				continue
			}

			fmt.Println("Someone connected to me, who is listening at: ", listenerAddrOfSender)

			p.Connect(listenerAddrOfSender)
			p.connectionAddr = listenerAddrOfSender

		case MsgEnvelope:
			var recivedMsg string

			if err := json.Unmarshal(recivedEnvelope.Payload, &recivedMsg); err != nil {
				fmt.Println("Error while receiving msg envelope")
				continue
			}

			fmt.Println("I " + p.name + " received msg: " + recivedMsg)
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

func (p *peer) SendMsg(msg string) {
	msgToBeSend, err := json.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}

	envelopeToSend := Envelope{Type: MsgEnvelope, Payload: msgToBeSend}
	p.SendEnvelope(envelopeToSend)
}

func (p *peer) GetListenerAddr() string {
	return p.ownListeningAddr
}
