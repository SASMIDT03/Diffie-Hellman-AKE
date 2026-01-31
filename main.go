package main

import (
	"diffie-hellman-ake/peer"
	"fmt"
)

func main() {
	fmt.Println("Hello World")

	alice := peer.NewPeer("Alice")
	bob := peer.NewPeer("Bob")

	listenerAddrAlice := alice.GetListenerAddr()
	fmt.Println("listener addr of Alice: ", listenerAddrAlice)
	fmt.Println("listener addr of Bob: ", bob.GetListenerAddr())

	bob.Connect(listenerAddrAlice)

	// cursed, I know
	for {
	}
}
